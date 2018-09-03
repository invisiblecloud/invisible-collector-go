package ic

import "fmt"

type ServerError string

func (e ServerError) Error() string {
	return string(e)
}

type ConflictError struct {
	Msg        string
	ConflictId string
}

func (e ConflictError) Error() string {
	return e.Msg + " (conflicting id: " + e.ConflictId + ")"
}

type HttpStatusCodeError int

func (e HttpStatusCodeError) Error() string {
	return fmt.Sprintf("%v", int(e))
}
