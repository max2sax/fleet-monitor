package dto

type ErrorNotFound struct {
	What string
}

func (e *ErrorNotFound) Error() string {
	return e.What + ": not found"
}
