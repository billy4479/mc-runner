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
	MinecraftFlags []string `json:"minecraft_flags"`
	ServerName     string   `json:"server_name"`
	BotPrefix      string   `json:"bot_prefix"`
}

type Driver struct {
	config       DriverConfig
	globalConfig *config.Config
	process      *exec.Cmd

	stdioMutex        sync.Mutex
	stdin             *bufio.Writer
	stdoutMultiplexer *SubscribeWriter
	stdoutBuffer      *bytes.Buffer

	stdoutPipe io.Reader
	stderrPipe io.Reader

	stateBroadcaster *StateBroadcaster

	isRunning bool
}

func (driver *Driver) getJavaExe() string {
	switch driver.config.JavaVersion {
	case 8:
		return driver.globalConfig.Java8
	case 17:
		return driver.globalConfig.Java17
	case 21:
		return driver.globalConfig.Java21
	case 25:
		return driver.globalConfig.Java25
	default:
		log.Warn().Int("specified_java_version", driver.config.JavaVersion).Msg("unsupported java version, falling back to java 21")
		return driver.globalConfig.Java21
	}
}

func (drv *Driver) makeNewExec() error {
	java := drv.getJavaExe()

	flags := make([]string, len(drv.config.JavaFlags))
	_ = copy(flags, drv.config.JavaFlags)
	flags = append(flags, "-jar", drv.config.JarFile)
	flags = append(flags, drv.config.MinecraftFlags...)

	log.Debug().Str("java", java).Strs("flags", flags).Msg("server command initialized")

	drv.process = exec.Command(java, flags...)
	drv.process.Dir = drv.globalConfig.WorldDir

	stdout, err := drv.process.StdoutPipe()
	if err != nil {
		return err
	}
	drv.stdoutPipe = stdout
	stderr, err := drv.process.StderrPipe()
	if err != nil {
		return err
	}
	drv.stderrPipe = stderr
	stdin, err := drv.process.StdinPipe()
	if err != nil {
		return err
	}
	drv.stdin = bufio.NewWriter(stdin)

	return nil
}

func NewDriver(config *config.Config) (*Driver, error) {
	b, err := os.ReadFile(filepath.Join(config.WorldDir, "mc-runner.yml"))
	if err != nil {
		return nil, err
	}

	driver := &Driver{
		stdioMutex:   sync.Mutex{},
		stdoutBuffer: bytes.NewBuffer([]byte{}),
		globalConfig: config,
	}
	driver.stateBroadcaster = NewStateBroadcaster(driver)

	err = yaml.Unmarshal(b, &driver.config)
	if err != nil {
		return nil, err
	}

	err = driver.makeNewExec()
	driver.stdoutMultiplexer = NewSubscribeWriter()
	driver.stdoutMultiplexer.Subscribe(driver.stdoutBuffer, math.MaxUint64)
	driver.stdoutMultiplexer.Subscribe(os.Stdout, math.MaxUint64-1)

	return driver, err
}

func (d *Driver) StateBroadcaster() *StateBroadcaster {
	return d.stateBroadcaster
}

func (d *Driver) Stdout() *SubscribeWriter {
	return d.stdoutMultiplexer
}

func (driver *Driver) GetState() *ServerState {
	return &ServerState{
		Version:       config.Version,
		ConnectUrl:    driver.globalConfig.ConnectUrl,
		ServerName:    driver.config.ServerName,
		IsRunning:     driver.isRunning,
		OnlinePlayers: driver.OnlinePlayers(),                         // TODO
		AutoStopTime:  time.Now().Add(driver.TimeBeforeStop()).Unix(), // TODO
		BotTag:        "@my_todo_bot",                                 // TODO
	}
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

func (drv *Driver) Start() error {
	if drv.isRunning {
		return nil
	}

	log.Debug().Str("server_name", drv.config.ServerName).Msg("starting server")

	err := drv.process.Start()
	if err != nil {
		return err
	}
	drv.isRunning = true

	{
		// We do not check errors here as (I believe) they should get caught by cmd.Wait()
		go func() {
			_, _ = io.Copy(drv.stdoutMultiplexer, drv.stdoutPipe)
		}()
		go func() {
			_, _ = io.Copy(drv.stdoutMultiplexer, drv.stderrPipe)
		}()
		go func() {
			_, _ = io.Copy(drv.stdin, os.Stdin)
		}()
	}

	go func() {
		drv.stateBroadcaster.sendUpdate()
		err := drv.process.Wait()

		if err != nil {
			log.Error().Err(err).
				Int("exit_code", drv.process.ProcessState.ExitCode()).
				Msg("server terminated with an error")

			// At this point we don't care if we fail
			_, _ = drv.stdoutMultiplexer.Write(fmt.Appendf(nil, "\n\n[ERR]: Server terminated with an error: %v\n\n", err))
		} else {
			log.Info().Msg("server terminated with no errors")
		}

		drv.isRunning = false
		err = drv.makeNewExec()
		if err != nil {
			log.Fatal().Err(err).Msg("could not create a new process")
		}

		drv.stateBroadcaster.sendUpdate()
	}()

	return nil
}

func (d *Driver) Stop() {
	if !d.isRunning {
		return
	}

	err := d.process.Process.Signal(os.Interrupt)
	if err != nil {
		log.Error().Err(err).Msg("could not send SIGINT, killing")
		err = d.process.Process.Kill()
		if err != nil {
			log.Error().Err(err).Msg("could not kill process")
		}
	}

	// Errors will already be handled somewhere else
	_ = d.process.Wait()
}
