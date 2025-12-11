package driver

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/billy4479/mc-runner/internal/config"
	"github.com/goccy/go-yaml"
	"github.com/rs/zerolog/log"
)

type DriverConfig struct {
	JavaVersion    int      `json:"java_version"`
	JarFile        string   `json:"jar_file"`
	JavaFlags      []string `json:"java_flags"`
	EulaAccepted   bool     `json:"eula_accepted"`
	ServerName     string   `json:"server_name"`
	HasIntegration bool     `json:"has_integration"`
}

type Driver struct {
	config  DriverConfig
	process *exec.Cmd

	stdioMutex        sync.Mutex
	stdin             *bufio.Writer
	stdoutMultiplexer *SubscribeWriter
	stdoutBuffer      *bytes.Buffer

	isRunning bool
}

func NewDriver(config *config.Config) (*Driver, error) {
	b, err := os.ReadFile(filepath.Join(config.WorldDir, "mc-runner.yml"))
	if err != nil {
		return nil, err
	}

	driver := &Driver{
		stdioMutex:   sync.Mutex{},
		stdoutBuffer: bytes.NewBuffer([]byte{}),
	}

	err = yaml.Unmarshal(b, &driver.config)
	if err != nil {
		return nil, err
	}

	java := func() string {
		switch driver.config.JavaVersion {
		case 8:
			return config.Java8
		case 17:
			return config.Java17
		case 21:
			return config.Java21
		default:
			log.Warn().Int("specified_java_version", driver.config.JavaVersion).Msg("unsupported java version, falling back to java 21")
			return config.Java21
		}
	}()

	flags := make([]string, len(driver.config.JavaFlags))
	_ = copy(flags, driver.config.JavaFlags)
	flags = append(flags, "-jar", driver.config.JarFile, "-nogui")
	log.Debug().Str("java", java).Strs("flags", flags).Msg("server command initialized")

	driver.process = exec.Command(java, flags...)
	driver.process.Dir = config.WorldDir

	stdout, err := driver.process.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := driver.process.StderrPipe()
	if err != nil {
		return nil, err
	}
	driver.stdoutMultiplexer = NewSubscribeWriter()
	driver.stdoutMultiplexer.Subscribe(driver.stdoutBuffer, math.MaxUint64)
	driver.stdoutMultiplexer.Subscribe(os.Stdout, math.MaxUint64-1)

	go func() {
		// We do not check errors here as (I believe) they should get caught by cmd.Wait()
		_, _ = io.Copy(driver.stdoutMultiplexer, stdout)
	}()
	go func() {
		// We do not check errors here as (I believe) they should get caught by cmd.Wait()
		_, _ = io.Copy(driver.stdoutMultiplexer, stderr)
	}()

	stdin, err := driver.process.StdinPipe()
	if err != nil {
		return nil, err
	}
	driver.stdin = bufio.NewWriter(stdin)

	return driver, err
}

func (d *Driver) Stdout() *SubscribeWriter {
	return d.stdoutMultiplexer
}

func (d *Driver) Stdin() *bufio.Writer {
	return d.stdin
}

func (d *Driver) IsRunning() bool {
	return d.isRunning
}

func (d *Driver) ServerName() string {
	return d.config.ServerName
}

func (d *Driver) OnlinePlayers() []string {
	return []string{
		"these names",
		"are",
		"placeholder",
		"some could be very long, like very looooooooooong",
	}
}

func (d *Driver) TimeBeforeStop() time.Duration {
	return 5 * time.Minute
}

func (d *Driver) Start() error {
	if d.isRunning {
		return nil
	}

	log.Debug().Str("server_name", d.config.ServerName).Msg("starting server")

	err := d.process.Start()
	if err != nil {
		return err
	}
	d.isRunning = true

	go func() {
		err := d.process.Wait()

		if err != nil {
			log.Error().Err(err).
				Int("exit_code", d.process.ProcessState.ExitCode()).
				Msg("server terminated with an error")

			// At this point we don't care if we fail
			_, _ = d.stdoutMultiplexer.Write(fmt.Appendf(nil, "\n\n[ERR]: Server terminated with an error: %v\n\n", err))
		}

		log.Info().Msg("server terminated with no errors")
		d.isRunning = false
	}()

	return nil
}
