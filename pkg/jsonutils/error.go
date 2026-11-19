package jsonutils

// JsonError represents an error during JSON serialization or deserialization.
type JsonError struct {
	Message string
}

func (e JsonError) Error() string {
	return e.Message
}

// ValidationError represents an error during data validation.
type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}
