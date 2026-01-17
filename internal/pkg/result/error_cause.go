package result

type ErrorCause uint8

const (
	ENTITY_DUPLICATE ErrorCause = iota + 1
	ENTITY_NOT_FOUND
	INTERNAL_SERVICE_ERROR
	CREDENTIALS_ERROR
	TOKEN_GENERATION_ERROR
	TOKEN_EXPIRED_ERROR
	INTERNAL_LOGIC_ERROR
	BAD_REQUEST_ERROR
	UNKNOWN_ERROR
)

var ErrorMessages = map[ErrorCause]string{
	ENTITY_DUPLICATE:       "Entity already exists",
	ENTITY_NOT_FOUND:       "Entity not found",
	INTERNAL_SERVICE_ERROR: "Internal service error",
	CREDENTIALS_ERROR:      "Invalid credentials",
	TOKEN_GENERATION_ERROR: "Token generation error",
	TOKEN_EXPIRED_ERROR:    "Token expired",
	INTERNAL_LOGIC_ERROR:   "Internal logic error",
	BAD_REQUEST_ERROR:      "Bad request",
	UNKNOWN_ERROR:          "Unknown error",
}

func (e ErrorCause) String() string {
	if v, ok := ErrorMessages[e]; ok {
		return v
	}

	return ErrorMessages[UNKNOWN_ERROR]
}
