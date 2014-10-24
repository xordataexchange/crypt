// Package standard implements base64 encoding in the following format:
//
//   base64(data)
//
package standard

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
)

// Deocde decodes data using the standard codec.
func Decode(data []byte) ([]byte, error) {
	decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewBuffer(data))
	bytes, err := ioutil.ReadAll(decoder)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// Encode encodes data to a base64 encoded using the standard codec.
func Encode(data []byte) ([]byte, error) {
	buffer := new(bytes.Buffer)
	encoder := base64.NewEncoder(base64.StdEncoding, buffer)
	encoder.Write(data)
	encoder.Close()
	return buffer.Bytes(), nil
}
