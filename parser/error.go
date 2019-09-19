package parser

import "fmt"

type ParserError struct {
	Message       string
	StartPosition Position
	EndPosition   Position
	ErrorStack    []error
	Fatal         bool
}

func (e *ParserError) Append(err error) *ParserError {
	if e.ErrorStack == nil {
		e.ErrorStack = make([]error, 0)
	}
	e.ErrorStack = append(e.ErrorStack, err)
	return e
}

func (e ParserError) Error() string {
	txt := fmt.Sprintf("Parser error at %s (up to %s): %s", e.StartPosition.String(), e.EndPosition.String(), e.Message)
	if e.ErrorStack != nil {
		txt += "\n" + "Following errors:\n"
		for _, err := range e.ErrorStack {
			txt += "    " + err.Error() + "\n"
		}
	}
	return txt
}

type ParserErrors []*ParserError

func (e ParserErrors) Error() string {
	str := ""
	for _, err := range e {
		str += err.Error() + "\n"
	}
	return str
}
