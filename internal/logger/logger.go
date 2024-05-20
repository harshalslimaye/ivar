package logger

import (
	"log"
	"runtime"

	"github.com/logrusorgru/aurora"
)

func Error(message string, err error) {
	if err == nil {
		return
	}

	_, file, line, found := runtime.Caller(1)

	if !found {
		file = "unknown"
		line = 0
	}

	log.Println(aurora.Sprintf("%s:%s - %s", aurora.Red(file), aurora.Red(line), aurora.Red(message)))
	log.Println(aurora.Sprintf("%s", aurora.Red(err.Error())))
}
