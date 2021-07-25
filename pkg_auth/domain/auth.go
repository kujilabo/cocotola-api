package domain

type TokenSet struct {
	AccessToken  string
	RefreshToken string
}

type UnauthorizedError struct {
	message string
}

func NewUnauthorizedError(message string) *UnauthorizedError {
	return &UnauthorizedError{
		message: message,
	}
}

func (e *UnauthorizedError) Error() string {
	return e.message
}
