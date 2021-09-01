package utils

import (
    "strings"
    "errors"
    "fmt"
)

func ReplaceStringAt(input string, startString string, endString string, newString string) (string, error) {
	startIndex := strings.Index(input, startString)
	if startIndex == -1 {
		return "", errors.New(fmt.Sprintf("Could not find: %s", startString))
	}
	startIndex += len(startString)

	endIndex := strings.Index(input[startIndex:], endString)
	if endIndex == -1 {
		return "", errors.New(fmt.Sprintf("Could not find: %s", endString))
	}
	endIndex = startIndex + endIndex

	var sb strings.Builder

	sb.WriteString(input[:startIndex])
	sb.WriteString(newString)
	sb.WriteString(input[endIndex:])

	return sb.String(), nil
}
