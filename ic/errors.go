package ic

import "fmt"

// Represent an invisible corrector error.
//
// This error can be for example an invalid company API key.
type ServerError string

func (e ServerError) Error() string {
	return string(e)
}

// Represents a model conflict error, where the ConflictId field
// represents the conflicting ID in the invisible collector's database
type ConflictError struct {
	msg        string
	ConflictId string
}

func (e ConflictError) Error() string {
	return e.msg + " (conflicting id: " + e.ConflictId + ")"
}

// Represents a HTTP protocol/status code error that isn't a specific invisible collector error.
type HttpStatusCodeError int

func (e HttpStatusCodeError) Error() string {
	return fmt.Sprintf("%v", int(e))
}
