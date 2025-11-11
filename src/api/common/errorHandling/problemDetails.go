package errorHandling

type ProblemDetails struct {
	ErrorType string

	Title string

	Status int

	Detail string

	Instance *string

	Extenstions *map[string]string
}
