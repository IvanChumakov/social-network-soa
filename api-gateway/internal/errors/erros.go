package errors

type JWTTokenEmpty struct{}

func (jte *JWTTokenEmpty) Error() string {
	return "Token is empty"
}

type JWTTokenInvalid struct{}

func (jti *JWTTokenInvalid) Error() string {
	return "Token is invalid"
}
