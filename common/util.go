package common

import (
	"bytes"
	"encoding/gob"
)

func Encode(data interface{}) ([]byte, error) {
	gob.Register(data)
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(&data)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func Decode(data []byte, src interface{}) (interface{}, error) {
	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&src)
	if err != nil {
		return src, err
	}
	return src, nil
}
