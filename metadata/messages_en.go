package metadata

const ER_MODEL_NOT_A_POINTER_OR_STRUCT = "%s() failed: model must be a struct or pointer to a struct"
const ER_MODEL_NOT_A_POINTER = "%s() failed: model must be a pointer to a struct"
const ER_VAL_CAST = "value of field %s can not be cast to %s"

// validation errors
const (
	ER_VALIDATION        = "validation error: %s"
	ER_VAL_REQUIRED      = "field '%s', value is required"
	ER_VAL_NOT_IN_LIST   = "field '%s', value is not in the list "
	ER_VAL_LEN_TOO_LONG  = "field '%s', value is too long"
	ER_VAL_TOO_SMALL     = "field '%s', value is too small"
	ER_VAL_TOO_BIG       = "field '%s', value is too big"
	ER_VAL_LEN_TOO_SHORT = "field '%s', value is too short"
	ER_VAL_REG_EXP       = "field '%s', value should comply with regular expression"
	ER_VAL_LEN_NOT_FIX   = "field '%s', value length should be fixed"
	ER_VAL_PRECISION     = "field '%s', float precision is exceeded maximum value"
	ER_VAL_VAL_LIST      = "field '%s', value should be in list of values"

	ER_VAL_AR_INVALID   = "field '%s', array format is invalid"
	ER_VAL_AR_MAX_LEN   = "field '%s', array count exceeded maximum value"
	ER_VAL_AR_MIN_LEN   = "field '%s', array count less than minimal value"
	ER_VAL_AR_FIXED_LEN = "field '%s', array count should be of fixed length"
)
