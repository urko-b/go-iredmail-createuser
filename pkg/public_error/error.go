package public_error

type Public interface {
	error
	IsPublic() bool
	StatusCode() int
	OriginalError() error
}

type PublicError struct {
	publicError   error
	originalError error
	statusCode    int
}

func (e PublicError) IsPublic() bool {
	return true
}

func (e PublicError) OriginalError() error {
	return e.originalError
}

func (e PublicError) Error() string {
	return e.publicError.Error()
}

func (e PublicError) StatusCode() int {
	return e.statusCode
}

func New(publicError, originalError error, statusCode int) *PublicError {
	return &PublicError{
		publicError,
		originalError,
		statusCode,
	}
}

func IsPublicErr(err error) bool {
	if publicErr, ok := err.(Public); ok && publicErr.IsPublic() { // nolint: errorlint
		return true
	}
	return false
}
