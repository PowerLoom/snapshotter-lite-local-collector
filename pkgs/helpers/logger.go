package helpers

import (
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"strconv"
)

func InitLogger() {
	log.SetOutput(io.Discard) // Send all logs to nowhere by default

	log.SetReportCaller(true)

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	// Function to create a new lumberjack logger
	newLogger := func(fileName string) io.Writer {
		return &lumberjack.Logger{
			Filename:   fileName,
			MaxSize:    100, // megabytes
			MaxBackups: 7,
			MaxAge:     30,   //days
			Compress:   true, // enabled by default
		}
	}

	log.AddHook(&writer.Hook{ // Send logs with level higher than warning to stderr
		Writer: os.Stderr,
		LogLevels: []log.Level{
			log.PanicLevel,
			log.FatalLevel,
			log.ErrorLevel,
			log.WarnLevel,
		},
	})

	log.AddHook(&writer.Hook{ // Send info and debug logs to stdout
		Writer: os.Stdout,
		LogLevels: []log.Level{
			log.TraceLevel,
			log.InfoLevel,
			log.DebugLevel,
		},
	})

	log.AddHook(&writer.Hook{
		Writer: newLogger("logs/snapshotter-lite-local-collector/critical.log"),
		LogLevels: []log.Level{
			log.PanicLevel,
			log.FatalLevel,
		},
	})

	log.AddHook(&writer.Hook{
		Writer: newLogger("logs/snapshotter-lite-local-collector/error.log"),
		LogLevels: []log.Level{
			log.ErrorLevel,
		},
	})

	log.AddHook(&writer.Hook{
		Writer: newLogger("logs/snapshotter-lite-local-collector/warning.log"),
		LogLevels: []log.Level{
			log.WarnLevel,
		},
	})

	log.AddHook(&writer.Hook{
		Writer: newLogger("logs/snapshotter-lite-local-collector/info.log"),
		LogLevels: []log.Level{
			log.InfoLevel,
		},
	})

	log.AddHook(&writer.Hook{
		Writer: newLogger("logs/snapshotter-lite-local-collector/debug.log"),
		LogLevels: []log.Level{
			log.DebugLevel,
		},
	})

	log.AddHook(&writer.Hook{
		Writer: newLogger("logs/snapshotter-lite-local-collector/trace.log"),
		LogLevels: []log.Level{
			log.TraceLevel,
		},
	})

	// Set the default log level
	if len(os.Args) < 2 {
		log.SetLevel(log.DebugLevel)
	} else {
		logLevel, err := strconv.ParseUint(os.Args[1], 10, 32)
		if err != nil || logLevel > 6 {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.Level(logLevel))
		}
	}
}
