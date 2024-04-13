package errors

type InvalidDateFormat struct {
	message string
	err     error
}

func (e InvalidDateFormat) Error() string {
	return e.message
}

func (e InvalidDateFormat) Unwrap() error {
	return e.err
}

func NewInvalidDateFormat(message string, err error) error {
	return InvalidDateFormat{message, err}
}

type InvalidRepeatFormat struct {
	message string
	err     error
}

func (e InvalidRepeatFormat) Error() string {
	return e.message
}

func (e InvalidRepeatFormat) Unwrap() error {
	return e.err
}

func NewInvalidRepeatFormat(message string, err error) error {
	return InvalidRepeatFormat{message, err}
}

type InvalidTitleFormat struct {
	message string
	err     error
}

func (e InvalidTitleFormat) Error() string {
	return e.message
}

func (e InvalidTitleFormat) Unwrap() error {
	return e.err
}

func NewInvalidTitleFormat(message string, err error) error {
	return InvalidTitleFormat{message, err}
}

type TaskNotExists struct {
	message string
	err     error
}

func (e TaskNotExists) Error() string {
	return e.message
}

func (e TaskNotExists) Unwrap() error {
	return e.err
}

func NewTaskNotExists(message string, err error) error {
	return InvalidTitleFormat{message, err}
}

type AuthenticationError struct {
	message string
	err     error
}

func (e AuthenticationError) Error() string {
	return e.message
}

func (e AuthenticationError) Unwrap() error {
	return e.err
}

func NewAuthenticationError(message string, err error) error {
	return AuthenticationError{message, err}
}
