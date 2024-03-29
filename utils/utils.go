package utils

import (
	"strings"
)

func GetChannelsFromParams(params []string) (channels []string) {
	size := len(params)
	for i, param := range params {
		if i == size - 1 { // skip trailing
			break
		}
		if strings.HasPrefix(param, "#") {
			channels = append(channels, param)
		}
	}
	return channels
}

func SplitTrailingBySpace(input string) (lexemes []string) {
	size := len(input)
	for i := 0; i < size; i++ {
		if input[i] != ' ' {
			var j int
			for j = i + 1; j < size; j++ {
				if input[j] == ' ' {
					break
				}
			}
			if i < j {
				lexemes = append(lexemes, input[i:j])
				i = j
			}
		} else if input[i] == '(' || input[i] == ')' {
			lexemes = append(lexemes, input[i:i+1])
		}
	}
	return lexemes
}
