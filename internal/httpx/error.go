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

type BadRequestError struct {
	ExternalError error
}

func (err *BadRequestError) Error() string {
	return err.ExternalError.Error()
}

type NotFoundError struct {
	ExternalError error
}

func (err *NotFoundError) Error() string {
	return err.ExternalError.Error()
}

type ConflictError struct {
	ExternalError error
}

func (err *ConflictError) Error() string {
	return err.ExternalError.Error()
}

type UnexpectedRuntimeError struct {
	ExternalError error
}

func (err *UnexpectedRuntimeError) Error() string {
	return err.ExternalError.Error()
}
