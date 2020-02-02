package errors

import (
	"errors"
	"fmt"
)

const (
	NotImplemented = 213412341234
	ThingNotFound  = 404
	NoID           = 92184239084
	InternalError  = 500
)

type Error struct {
	Ops  []string
	Err  error
	Type interface{}
}

func (e *Error) Error() string {
	if e != nil && e.Err != nil {
		return e.Err.Error()
	}
	return ""
}

func (e *Error) OpStack() string {
	var err string
	for i := len(e.Ops) - 1; i >= 0; i-- {
		if err != "" {
			err += " - "
		}
		err += e.Ops[i]
	}
	return err
}

func Wrap(e interface{}, op string) error {
	if e == nil {
		return nil
	}
	switch e := e.(type) {
	case *Error:
		e.Ops = append(e.Ops, op)
		return e
	case Error:
		e.Ops = append(e.Ops, op)
		return &e
	case error:
		return New(e, nil, op)
	default:
		panic("Unknown error type sent to function")
	}
}

func New(err interface{}, typ interface{}, op string) *Error {
	e := new(Error)
	e.Ops = append(e.Ops, op)
	e.Type = typ
	switch err := err.(type) {
	case error:
		e.Err = err
	case string:
		e.Err = errors.New(err)
	case int:
		e.Err = errors.New(fmt.Sprintf("%d", err))
	default:
		panic("Unable to create error with argument type")
	}
	return e
}

func ErrorType(e interface{}) interface{} {
	err, ok := e.(*Error)
	if !ok {
		return nil
	}
	return err.Type
}
