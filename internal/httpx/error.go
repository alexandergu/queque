package httpx

type ConvertError struct {
	Message string
}

func (err *ConvertError) Error() string {
	return err.Message
}

type ValidationError struct {
	Message string
}

func (err *ValidationError) Error() string {
	return err.Message
}
