package v1

var (
	// common errors
	ErrSuccess             = newError(0, "ok")
	ErrBadRequest          = newError(400, "Bad Request")
	ErrForbidden           = newError(403, "Forbidden")
	ErrNotFound            = newError(404, "Not Found")
	ErrInternalServerError = newError(500, "Internal Server Error")

	// more biz errors
	ErrUnauthorized     = newError(2000, "Unauthorized")
	ErrEmailAlreadyUse  = newError(1001, "The email is already in use.")
	ErrEmailAlreadyUse1 = newError(1002, "T211212he email is already in use.")
)
