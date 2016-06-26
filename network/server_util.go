package network

type ServerError struct {
	E error
}

func NewServerError(e error) *ServerError {
	return &ServerError{
		E: e,
	}
}
