package cmdshim

import (
	"regexp"
	"strings"
)

func ConvertToSetCommand(key, value string) string {
	var line string
	key = strings.TrimSpace(key)
	value = strings.TrimSpace(value)
	if key != "" && value != "" && len(value) > 0 {
		line = "@SET " + key + "=" + ReplaceDollarWithPercentPair(value) + "\r\n"
	}
	return line
}

func ExtractVariableValuePairs(declarations []string) map[string]string {
	pairs := make(map[string]string)
	for _, declaration := range declarations {
		split := strings.Split(declaration, "=")
		if len(split) == 2 {
			pairs[strings.TrimSpace(split[0])] = strings.TrimSpace(split[1])
		}
	}
	return pairs
}

func ConvertToSetCommands(variableString string) string {
	variableValuePairs := ExtractVariableValuePairs(strings.Split(variableString, " "))
	var variableDeclarationsAsBatch string
	for key, value := range variableValuePairs {
		variableDeclarationsAsBatch += ConvertToSetCommand(key, value)
	}
	return variableDeclarationsAsBatch
}

func ReplaceDollarWithPercentPair(value string) string {
	dollarExpressions := regexp.MustCompile(`\$\{?([^$@#?\- \t{}:]+)\}?`)
	result := ""
	startIndex := 0
	for _, match := range dollarExpressions.FindAllStringSubmatchIndex(value, -1) {
		betweenMatches := value[startIndex:match[0]]
		result += betweenMatches + "%" + value[match[2]:match[3]] + "%"
		startIndex = match[1]
	}
	result += value[startIndex:]
	return result
}
