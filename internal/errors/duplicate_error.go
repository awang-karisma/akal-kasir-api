package errors

type DuplicateError struct {
	Field string
	Value string
}

func (e *DuplicateError) Error() string {
	return "duplicate entry: " + e.Field + " = " + e.Value
}

func NewDuplicateError(field, value string) error {
	return &DuplicateError{Field: field, Value: value}
}

func IsDuplicateError(err error) bool {
	_, ok := err.(*DuplicateError)
	return ok
}
