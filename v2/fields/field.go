package fields

const (
	quoteChar byte = 34
	jsonNull       = "null"
)

const (
	ErrUnsupportedType = "unsupported data type for %s field: %v"
)

func ValIsNull(extVal []byte) bool {
	return string(extVal) == `null`
}

func RemoveQuotes(extVal []byte) string {
	var vStr string
	if extVal[0] == quoteChar && byte(extVal[len(extVal)-1]) == quoteChar {
		vStr = string(extVal[1 : len(extVal)-1])
	} else {
		vStr = string(extVal)
	}
	return vStr
}
