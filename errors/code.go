package errors

// Define error codes as constants.
const (
	NotFound             = "notFound"
	Unauthorized         = "unauthorized"
	InternalServerError  = "internalServerError"
	BodyParse            = "bodyParse"
	Validator            = "validator"
	QueryError           = "queryError"
	CacheError           = "cacheError"
	Forbidden            = "forbidden"
	MissingRequiredParam = "missingRequiredParam"
	InvalidParam         = "invalidParam"
	OutOfSync            = "outOfSync"
	// Add more error codes as needed.
)
