package internal

import "fmt"

type HttpStatusCodeError int

func (e HttpStatusCodeError) Error() string {
	return fmt.Sprintf("")
}
