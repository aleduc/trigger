package logger

import (
	"fmt"
)

//go:generate go run github.com/golang/mock/mockgen --source=logger.go --destination=logger_mock.go --package=logger

// Logger can be defined just once, since it is low level interface, which will be imported in many places.
type Logger interface {
	Error(args ...interface{})
	Info(input string)
}

type Simple struct {
}

func NewSimple() *Simple {
	return &Simple{}
}

func (s *Simple) Info(input string) {
	fmt.Println(input)
}

func (s *Simple) Error(args ...interface{}) {
	fmt.Println(args...)
}
