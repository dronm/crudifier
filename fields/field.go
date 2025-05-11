package fields

const (
	QUOTE_CHAR byte = 34
	JSON_NULL       = "null"
)

const (
	ER_UNSUPPORTED_TYPE = "unsupported data type for %s field: %v"
)

func ValIsNull(extVal []byte) bool {
	return string(extVal) == `null`
}

func RemoveQuotes(extVal []byte) string {
	var v_str string
	if extVal[0] == QUOTE_CHAR && byte(extVal[len(extVal)-1]) == QUOTE_CHAR {
		v_str = string(extVal[1 : len(extVal)-1])
	} else {
		v_str = string(extVal)
	}
	return v_str
}
