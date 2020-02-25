package errors

import (
	"errors"
	"fmt"
)

// Error types to be used when creating a new error to flag what type
// it is. Optional.
const (
	NotImplemented = 501
	ThingNotFound  = 404
	NoID           = 92184239084
	InternalError  = 500
	Forbidden      = 403
)

// Error is container for a normal go error and extra information useful
// for this specific program, including a trace of which methods have called it
// and if the message should be passed externally
type Error struct {
	Ops    []string
	Err    error
	Type   interface{}
	Public bool
}

func (e *Error) Error() string {
	if e != nil && e.Err != nil {
		return e.Err.Error()
	}
	return ""
}

func (e *Error) opStack() string {
	var err string
	for i := len(e.Ops) - 1; i >= 0; i-- {
		if err != "" {
			err += " - "
		}
		err += e.Ops[i]
	}
	return err
}

// Wrap an existing error into the custom error type for this program
// The op string is the best name (usually function name) for the operation
// that's wrapping the error
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
		// assuming any error we did not explicitly create is not for public consumption
		return New(e, nil, op, false)
	default:
		panic("Unknown error type sent to function")
	}
}

// New returns an instatiated errors.Error
func New(err interface{}, typ interface{}, op string, public bool) *Error {
	e := new(Error)
	e.Public = public
	e.Ops = append(e.Ops, op)
	e.Type = typ
	switch err := err.(type) {
	case error:
		e.Err = err
	case string:
		e.Err = errors.New(err)
	case int:
		e.Err = fmt.Errorf("%d", err)
	default:
		panic("Unable to create error with argument type")
	}
	return e
}

// ErrorType returns the type of the error if it is a errors.Error type
// and if it is set (nil otherwise)
func ErrorType(e interface{}) interface{} {
	err, ok := e.(*Error)
	if !ok {
		return nil
	}
	return err.Type
}

// Print prints out the details of the provided error to the console included
// meta information if the error is a type to contain it (*Error)
func Print(err error) {
	if err == nil {
		return
	}
	errMsg := ` Error:  `
	errMsg += err.Error()
	errMsg += `
-------------------------------------------------------------------------------`
	if e, ok := err.(*Error); ok {
		errMsg += `
 Stack:  `
		errMsg += e.opStack()
		errMsg += `
 Type:   `
		errMsg += fmt.Sprintf("%v", e.Type)
		errMsg += `
 Public: `
		errMsg += fmt.Sprintf("%t", e.Public)
	}
	errMsg += `
`
	fmt.Println(errMsg)
}
