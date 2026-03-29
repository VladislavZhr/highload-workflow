package kafka

import "fmt"

type ConfigError string

func (e ConfigError) Error() string {
	return string(e)
}

func ErrInvalidConfig(msg string) error {
	return ConfigError("invalid kafka config: " + msg)
}

var ErrMalformedMessage = fmt.Errorf("malformed message")
