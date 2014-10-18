package secconf

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io"
	"io/ioutil"

	"code.google.com/p/go.crypto/openpgp"
)

func Decode(data []byte, secertKeyring io.Reader) ([]byte, error) {
	decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewBuffer(data))
	entityList, err := openpgp.ReadKeyRing(secertKeyring)
	if err != nil {
		return nil, err
	}
	md, err := openpgp.ReadMessage(decoder, entityList, nil, nil)
	if err != nil {
		return nil, err
	}
	gzReader, err := gzip.NewReader(md.UnverifiedBody)
	if err != nil {
		return nil, err
	}
	defer gzReader.Close()
	bytes, err := ioutil.ReadAll(gzReader)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func Encode(data []byte, keyring io.Reader) ([]byte, error) {
	entityList, err := openpgp.ReadKeyRing(keyring)
	if err != nil {
		return nil, err
	}
	buffer := new(bytes.Buffer)
	encoder := base64.NewEncoder(base64.StdEncoding, buffer)
	pgpWriter, err := openpgp.Encrypt(encoder, entityList, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	gzWriter := gzip.NewWriter(pgpWriter)
	if _, err := gzWriter.Write(data); err != nil {
		return nil, err
	}
	if err := gzWriter.Close(); err != nil {
		return nil, err
	}
	if err := pgpWriter.Close(); err != nil {
		return nil, err
	}
	if err := encoder.Close(); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
