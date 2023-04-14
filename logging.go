package logging //льет логи и в аутапут и в файл

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

type writerHook struct {
	Writer    []io.Writer // сюда отправляются лог файлы
	LogLevels []logrus.Level
}

func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	for _, w := range hook.Writer {
		w.Write([]byte(line))
	}
	return err
}

var e *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func GetLogger() Logger {
	return Logger{e}
}


func (l *Logger) GetLoggerWithField(k string, v interface{}) Logger {
	return Logger{l.WithField(k, v)}
}

func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}

func Init() {
	l := logrus.New()
	l.SetReportCaller(true)
	l.Formatter = &logrus.TextFormatter{ //передается рантайм врейм, в котором есть информация о логировании и инфа о лайне и функция внутри которой мы находимс и возвращает функцию и имя файла
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
		},
		DisableColors: false,
		FullTimestamp: true,
	}

	err := os.Mkdir("logs", 0644)
	if err != nil {
		panic(err)

	}

	allFile, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		panic(err)
	}
	l.SetOutput(io.Discard) //не записываь

	l.AddHook(&writerHook{
		Writer:    []io.Writer{allFile, os.Stdout},
		LogLevels: logrus.AllLevels,
	})

	l.SetLevel(logrus.TraceLevel)
}

//info.log -- info, warning, error, critical
//debug.log -- info, warning, error, critical, debug
//trace.log -- info, warning, error, critical, debug, trace
//error.log -- error, critical

//kafka -- info, debug
//file -- error, trace
//stdout-- warning, critical
