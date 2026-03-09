package metadata

const (
	ErrModelNotAPointerOrStruct = "%s() failed: model must be a struct or pointer to a struct"
	ErrModelNotAPointer         = "%s() failed: model must be a pointer to a struct"
	ErrValCast                  = "value of field %s can not be cast to %s"
)

// validation errors
const (
	ErrValidation     = "validation error: %s"
	ErrValRequired    = "field '%s', value is required"
	ErrValNotInList   = "field '%s', value is not in the list "
	ErrValLenTooLong  = "field '%s', value is too long"
	ErrValTooSmall    = "field '%s', value is too small"
	ErrValTooBig      = "field '%s', value is too big"
	ErrValLenTooShort = "field '%s', value is too short"
	ErrValRegExp      = "field '%s', value should comply with regular expression"
	ErrValLenNotFix   = "field '%s', value length should be fixed"
	ErrValPrecision   = "field '%s', float precision is exceeded maximum value"
	ErrValValList     = "field '%s', value should be in list of values"

	ErrValArInvalid  = "field '%s', array format is invalid"
	ErrValArMaxLen   = "field '%s', array count exceeded maximum value"
	ErrValArinLen    = "field '%s', array count less than minimal value"
	ErrValArFixedLen = "field '%s', array count should be of fixed length"
)
