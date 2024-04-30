package helper

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/logrusorgru/aurora"
)

func AskQuestion(message string, defaultValue string) string {
	reader := bufio.NewReader(os.Stdin)

	if defaultValue == "" {
		fmt.Printf(aurora.Sprintf(aurora.Cyan(message + ": ")))
	} else {
		fmt.Printf(aurora.Sprintf(aurora.Cyan(message+": "), aurora.White("("+defaultValue+")")))
	}

	value, _ := reader.ReadString('\n')
	value = strings.TrimSpace(value)

	if value == "" {
		return defaultValue
	}

	return value
}

func ShowInfo(step, emoji, message string) string {
	return fmt.Sprintf("%s %s %s...", emoji, aurora.Cyan(step), message)
}
