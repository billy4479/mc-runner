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

	stdoutBufferMutex sync.Mutex
	stdoutBuffer      *bytes.Buffer

	stdin      io.Writer
	stdoutPipe io.Reader
	stderrPipe io.Reader

	stateBroadcaster *Broadcaster[*ServerState]
	chatBroadcaster  *Broadcaster[string]

	stopTime time.Time

	onlinePlayers *OnlinePlayers

	isRunning bool
}

type ServerState struct {
	Version       string   `json:"version"`
	ConnectUrl    string   `json:"connect_url"`
	ServerName    string   `json:"server_name"`
	IsRunning     bool     `json:"is_running"`
	OnlinePlayers []string `json:"online_players"` // TODO
	AutoStopTime  int64    `json:"auto_stop_time"` // TODO
	BotTag        string   `json:"bot_tag"`        // TODO
}

func (drv *Driver) getJavaExe() string {
	switch drv.config.JavaVersion {
	case 8:
		return drv.globalConfig.Java8
	case 17:
		return drv.globalConfig.Java17
	case 21:
		return drv.globalConfig.Java21
	case 25:
		return drv.globalConfig.Java25
	default:
		log.Warn().Int("specified_java_version", drv.config.JavaVersion).Msg("unsupported java version, falling back to java 21")
		return drv.globalConfig.Java21
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
	drv.stdin = stdin

	drv.stdoutBufferMutex.Lock()
	drv.stdoutBuffer = new(bytes.Buffer)
	drv.stdoutBufferMutex.Unlock()

	return nil
}

func NewDriver(config *config.Config) (*Driver, error) {
	b, err := os.ReadFile(filepath.Join(config.WorldDir, "mc-runner.yml"))
	if err != nil {
		return nil, err
	}

	drv := &Driver{
		globalConfig: config,

		stateBroadcaster: NewBroadcaster[*ServerState](),
		chatBroadcaster:  NewBroadcaster[string](),

		stdoutBufferMutex: sync.Mutex{},
		stopTime:          time.Unix(math.MaxInt64, 0),
	}

	drv.onlinePlayers = NewOnlinePlayers(drv)

	err = yaml.Unmarshal(b, &drv.config)
	if err != nil {
		return nil, err
	}

	err = drv.makeNewExec()

	return drv, err
}

func (drv *Driver) StateBroadcaster() *Broadcaster[*ServerState] {
	return drv.stateBroadcaster
}

func (drv *Driver) GetState() *ServerState {
	return &ServerState{
		Version:       config.Version,
		ConnectUrl:    drv.globalConfig.ConnectUrl,
		ServerName:    drv.config.ServerName,
		IsRunning:     drv.isRunning,
		OnlinePlayers: drv.onlinePlayers.Get(),
		AutoStopTime:  drv.stopTime.Unix(),
		BotTag:        "@my_todo_bot", // TODO
	}
}

func (drv *Driver) Start() error {
	if drv.isRunning {
		return nil
	}

	logIfErr := func(err error) bool {
		if err == nil {
			return false
		}
		log.Error().Err(fmt.Errorf("driver: %w", err))
		return true
	}

	log.Debug().Str("server_name", drv.config.ServerName).Msg("starting server")

	err := drv.process.Start()
	if err != nil {
		return err
	}
	drv.isRunning = true

	drv.ScheduleStop()

	go func() {
		drv.stateBroadcaster.sendUpdate(drv.GetState())

		{
			// I believe errors should have been already caught by cmd.Wait()

			readFrom := func(reader io.Reader) {
				scanner := bufio.NewScanner(reader)

				for scanner.Scan() {
					line := append(scanner.Bytes(), '\n')
					lineStr := string(line)

					go func() {
						_, err = os.Stdout.Write(line)
						logIfErr(err)
					}()

					go drv.onlinePlayers.parseLine(lineStr)

					go func() {
						drv.stdoutBufferMutex.Lock()
						_, err = drv.stdoutBuffer.Write(line)
						drv.stdoutBufferMutex.Unlock()
						logIfErr(err)
					}()

					go drv.chatBroadcaster.sendUpdate(lineStr)
				}
			}

			go readFrom(drv.stdoutPipe)
			go readFrom(drv.stderrPipe)

			go func() {
				scanner := bufio.NewScanner(os.Stdin)
				for scanner.Scan() {
					_, _ = drv.stdin.Write(append(scanner.Bytes(), '\n'))
				}

				// This should be unreachable, since os.Stdin will never EOF.
				// Currently, typing commands in stdin while the server is closed, results in the command
				// being sent to the server as soon as it opens, since Write waits until it is actually
				// possible to write.
			}()
		}

		err := drv.process.Wait()

		// If the error is waitid... it means that the server has stopped naturally, therefore there is
		// no need to wait for it.
		if err != nil && err.Error() != "waitid: no child processes" {
			log.Error().Err(err).
				Int("exit_code", drv.process.ProcessState.ExitCode()).
				Msg("server terminated with an error")
		} else {
			log.Info().Msg("server terminated with no errors")
		}

		drv.isRunning = false
		err = drv.makeNewExec()
		if err != nil {
			log.Fatal().Err(err).Msg("could not create a new process")
		}

		drv.stateBroadcaster.sendUpdate(drv.GetState())
	}()

	return nil
}

func (drv *Driver) GetChatHistory() string {
	drv.stdoutBufferMutex.Lock()
	defer drv.stdoutBufferMutex.Unlock()

	return drv.stdoutBuffer.String()
}

func (drv *Driver) Stop() {
	if !drv.isRunning {
		return
	}

	drv.isRunning = false

	_, _ = drv.stdin.Write([]byte("/stop\n"))

	go func() {
		time.Sleep(10 * time.Second)

		if drv.process.ProcessState != nil && !drv.process.ProcessState.Exited() {
			err := drv.process.Process.Signal(os.Interrupt)
			if err != nil {
				log.Error().Err(err).Msg("could not send SIGINT, killing")
				err = drv.process.Process.Kill()
				if err != nil {
					log.Error().Err(err).Msg("could not kill process")
				}
			}
		}
	}()

	// Errors are already handled somewhere else
	_ = drv.process.Wait()
}

func (drv *Driver) ChatBroadcaster() *Broadcaster[string] {
	return drv.chatBroadcaster
}

func (drv *Driver) ScheduleStop() {
	go func() {
		thisStopTime := time.Now().Add(time.Duration(drv.globalConfig.MinutesBeforeStop) * time.Minute)

		log.Debug().Msg("schedule stop")
		drv.stopTime = thisStopTime
		drv.stateBroadcaster.sendUpdate(drv.GetState())

		time.Sleep(time.Duration(drv.globalConfig.MinutesBeforeStop) * time.Minute)
		if drv.onlinePlayers.Count() == 0 && thisStopTime.Equal(drv.stopTime) {
			log.Info().Msgf("no player online for %d minutes, stopping", drv.globalConfig.MinutesBeforeStop)
			drv.Stop()
		}
	}()
}
