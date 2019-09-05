package modules

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

func GzipEncode(in []byte) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	writer := gzip.NewWriter(buffer)
	defer writer.Close()

	_, err := writer.Write(in)
	if err != nil {
		return nil, err
	}

	err = writer.Flush()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func GzipDecode(in []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	decoded, _ := ioutil.ReadAll(reader)
	return decoded, nil
}
