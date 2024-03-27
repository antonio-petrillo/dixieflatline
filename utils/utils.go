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
