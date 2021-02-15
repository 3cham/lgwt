package logrus

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
)

var llog *log.Logger = nil
var cid, usecaseName string = "", ""

// LLogger returns logrus entry with two different fields
// to log with LLogger().Info("")
func LLogger() *log.Logger {
	return llog
}

// InitLogger return a new logger with default TextFormatter that writes to a io.Writer
// use file, _ := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
// to pass as log file writer
func InitLogger(writer io.Writer) {
	llog = log.New()
	llog.SetLevel(log.DebugLevel)
	llog.SetOutput(writer)
	llog.SetFormatter(&GaiaFormatter{})
}

// WithCid extends the default logger with command id
func WithCid(commandID string) {
	cid = commandID
}

// WithUsecase extends the default logger with usecase name
func WithUsecase(usecase string) {
	usecaseName = usecase
}

// AttachStdout attaches stdout as an extra channel to our standard file logger
func AttachStdout() {
	llog.AddHook(&writer.Hook{
		Writer: os.Stdout,
		LogLevels: []log.Level{
			log.InfoLevel,
			log.WarnLevel,
			log.ErrorLevel,
		},
	})
}

type GaiaFormatter struct {
}

// Format implement logrus Formatter interface
// so we have a log format like:
// 2021-02-15T17:08:36.507917+01:00 [info]  [CID=123 Usecase=test-usecase] Test Info
// 2021-02-15T17:08:36.508076+01:00 [debug]  [CID=123 Usecase=test-usecase] Only in file
func (f *GaiaFormatter) Format(entry *log.Entry) ([]byte, error) {
	b := &bytes.Buffer{}

	b.WriteString(entry.Time.Format(time.RFC3339Nano))

	b.WriteString(fmt.Sprintf(" [%s] ", entry.Level))

	b.WriteString(fmt.Sprintf(" [CID=%s Usecase=%s] ", cid, usecaseName))

	b.WriteString(entry.Message)
	b.WriteString("\n")
	return b.Bytes(), nil
}
