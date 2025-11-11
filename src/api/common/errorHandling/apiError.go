package errorHandling

// No need to import errorHandling since we're in the same package

const BASEURI = "aceclub.com/problem"

func NotFound(detail string, extentions *map[string]string) ProblemDetails {

	return ProblemDetails{
		ErrorType:   BASEURI + "/not-found",
		Title:       "Not Found",
		Status:      404,
		Detail:      detail,
		Extenstions: extentions,
	}
}

func Unauthorized(detail string, extentions *map[string]string) ProblemDetails {

	return ProblemDetails{
		ErrorType:   BASEURI + "/unauthorized",
		Title:       "Unauthorized",
		Status:      401,
		Detail:      detail,
		Extenstions: extentions,
	}
}

func BadRequest(detail string, extentions *map[string]string) ProblemDetails {

	return ProblemDetails{
		ErrorType:   BASEURI + "/bad-request",
		Title:       "Bad Request",
		Status:      400,
		Detail:      detail,
		Extenstions: extentions,
	}
}

func InternalServerError(detail string, exstentions *map[string]string) ProblemDetails {
	return ProblemDetails{
		ErrorType:   BASEURI + "/internal-error",
		Title:       "Internal Server Error",
		Status:      500,
		Detail:      detail,
		Extenstions: exstentions,
	}
}

func Forbidden(detail string, extentions *map[string]string) ProblemDetails {
	return ProblemDetails{
		ErrorType:   BASEURI + "/forbidden",
		Title:       "Forbidden",
		Status:      403,
		Detail:      detail,
		Extenstions: extentions,
	}
}
