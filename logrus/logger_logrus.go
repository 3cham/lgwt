package logrus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
)

var llog *log.Logger = nil
var cid, usecaseName string = "", ""

// Entry returns logrus entry with two different fields
// to log with Entry().Info("")
func Entry() *log.Entry {
	return llog.WithFields(log.Fields{
		"CID":     cid,
		"Usecase": usecaseName,
	})
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
func (f *GaiaFormatter) Format(entry *log.Entry) ([]byte, error) {
	b := &bytes.Buffer{}

	b.WriteString(entry.Time.Format(time.RFC3339Nano))

	b.WriteString(fmt.Sprintf(" [%s] ", entry.Level))

	b.WriteString(fmt.Sprintf(" [CID=%s Usecase=%s] ", cid, usecaseName))

	b.WriteString(marshal(encode(entry.Message)))
	return b.Bytes(), nil
}

func marshal(o interface{}) string {
	str, ok := o.(string)
	if ok {
		return str
	}
	data, err := json.Marshal(o)
	if err != nil {
		return fmt.Sprint(o)
	}
	return string(data)
}

func encode(message string) interface{} {
	if data := encodeForJsonString(message); data != nil {
		return data
	} else {
		return message
	}
}

func encodeForJsonString(message string) map[string]interface{} {
	// jsonstring
	inInterface := make(map[string]interface{})
	if err := json.Unmarshal([]byte(message), &inInterface); err != nil {
		//fmt.Print("err !!!! " , err.Error())
		return nil
	}
	return inInterface
}
