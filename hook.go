package filehook

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
	"encoding/json"
)

func NewHook(file string) (f *FileHook) {
	path := strings.Split(file, "/")
	if len(path) > 1 {
		exec.Command("mkdir", path[0]).Run()
	}
	w := NewFileWriter()
	config := fmt.Sprintf(`{"filename":"%s","maxdays":7}`, file)
	w.Init(config)
	return &FileHook{w}
}

type FileHook struct {
	W LoggerInterface
}

func (hook *FileHook) Fire(entry *logrus.Entry) (err error) {
	message, err := getMessage(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}
	switch entry.Level {
	case logrus.PanicLevel:
		fallthrough
	case logrus.FatalLevel:
		fallthrough
	case logrus.ErrorLevel:
		return hook.W.WriteMsg(fmt.Sprintf("[ERROR] %s", message), LevelError)
	case logrus.WarnLevel:
		return hook.W.WriteMsg(fmt.Sprintf("[WARN] %s", message), LevelWarn)
	case logrus.InfoLevel:
		return hook.W.WriteMsg(fmt.Sprintf("[INFO] %s", message), LevelInfo)
	case logrus.DebugLevel:
		return hook.W.WriteMsg(fmt.Sprintf("[DEBUG] %s", message), LevelDebug)
	default:
		return nil
	}
	return
}

func (hook *FileHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}

func getMessage(entry *logrus.Entry) (message string, err error) {

	message = message + fmt.Sprintf("%s", entry.Message)
	if len(entry.Data) > 0 {
		data, err := json.Marshal(entry.Data)
		if err == nil {
			message += " " + string(data)
		}
	}
	message += "\r\n"

	return
}
