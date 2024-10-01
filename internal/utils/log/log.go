package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	global *logrus.Logger
)

func init() {
	global = logrus.New()
	global.Out = os.Stdout
}

// Data дополнительная информация об ошибке
type Data map[string]any

// Debug логирует сообщение с уровнем debug
func Debug(err error, data Data) {
	global.WithFields(logrus.Fields(data)).Debug(err)
}

// Info логирует сообщение с уровнем info
func Info(msg any, data Data) {
	global.WithFields(logrus.Fields(data)).Info(msg)
}

// Error логирует сообщение с уровнем error
func Error(err error, data Data) {
	global.WithError(err).WithFields(logrus.Fields(data)).Error(err)
}

// Println выводит сообщение в поток вывода
func Println(msg any, data Data) {
	global.WithFields(logrus.Fields(data)).Println(msg)
}
