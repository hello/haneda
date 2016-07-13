package core

import (
	"github.com/go-kit/kit/log"
	"io/ioutil"
	"os"
)

type Loggers struct {
	Debug log.Logger
	Info  log.Logger
	Warn  log.Logger
	Error log.Logger
}

func LoggersFromConfig(config *LoggersConfig, app string) *Loggers {

	important := log.NewLogfmtLogger(os.Stderr)
	enabled := log.NewLogfmtLogger(os.Stdout)
	disabled := log.NewLogfmtLogger(ioutil.Discard)

	enabled = log.NewContext(enabled).With("ts", log.DefaultTimestampUTC, "app", app)

	loggers := &Loggers{
		Debug: disabled,
		Info:  disabled,
		Warn:  log.NewContext(important).With("level", "warn"),
		Error: log.NewContext(important).With("level", "error"),
	}

	// overrides defaults based on config
	if config != nil {
		if config.Debug {
			loggers.Debug = log.NewContext(enabled).With("level", "debug")
		}

		if config.Info {
			loggers.Info = log.NewContext(enabled).With("level", "info")
		}
	}
	return loggers
}
