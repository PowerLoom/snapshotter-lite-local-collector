package helpers

import (
	"fmt"
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
	if len(os.Args) < 2 {
		fmt.Println("Pass loglevel as an argument if you don't want default(INFO) to be set.")
		fmt.Println("Values to be passed for logLevel: ERROR(2),INFO(4),DEBUG(5)")
		log.SetLevel(log.DebugLevel)
	} else {
		logLevel, err := strconv.ParseUint(os.Args[1], 10, 32)
		if err != nil || logLevel > 6 {
			log.SetLevel(log.DebugLevel) //TODO: Change default level to error
		} else {
			//TODO: Need to come up with approach to dynamically update logLevel.
			log.SetLevel(log.Level(logLevel))
		}
	}
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	// Set up log rotation for all logs
	logPath := "logs/all.log"
	allLogger := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    100, // megabytes
		MaxBackups: 7,
		MaxAge:     30, //days
		Compress:   true,
	}

	// Set up log rotation for error logs
	errorLogPath := "logs/error.log"
	errorLogger := &lumberjack.Logger{
		Filename:   errorLogPath,
		MaxSize:    100, // megabytes
		MaxBackups: 7,
		MaxAge:     30, //days
		Compress:   true,
	}

	// Hook to write logs to the allLogger
	log.AddHook(&writer.Hook{
		Writer: allLogger,
		LogLevels: []log.Level{
			log.PanicLevel,
			log.FatalLevel,
			log.ErrorLevel,
			log.WarnLevel,
			log.InfoLevel,
			log.DebugLevel,
			log.TraceLevel,
		},
	})

	// Hook to write error logs to errorLogger
	log.AddHook(&writer.Hook{
		Writer: errorLogger,
		LogLevels: []log.Level{
			log.PanicLevel,
			log.FatalLevel,
			log.ErrorLevel,
		},
	})

	log.Infof("Logger initialized with log level %s", log.GetLevel().String())
}
