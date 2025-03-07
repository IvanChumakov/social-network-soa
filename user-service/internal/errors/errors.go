package errors

type NotFoundUserError struct{}

func (nf *NotFoundUserError) Error() string {
	return "User not found. Check credentials"
}

type LoginAlreadyTakenError struct{}

func (lfa *LoginAlreadyTakenError) Error() string {
	return "User with same login already exists"
}

type UpdateCredentialsError struct{}

func (usd *UpdateCredentialsError) Error() string {
	return "You can't update credentials"
}
