package customerror

type NotFoundError struct{}

func (nfe NotFoundError) Error() string {
	return "Ресурс не найден"
}
