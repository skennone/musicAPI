package data

import (
	"fmt"
	"strconv"
)

type Length int32

func (l Length) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d mins", l)
	// Needs to be surrounded by double quotes in order to be a valid JSON string
	quotedJSONValue := strconv.Quote(jsonValue)

	return []byte(quotedJSONValue), nil
}
