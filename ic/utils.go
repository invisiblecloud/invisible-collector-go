package ic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const jsonMime = "application/json"

func isWhitespaceString(value string) bool {
	return len(strings.TrimSpace(value)) == 0
}

func readCloseableBuffer(buffer io.ReadCloser) ([]byte, error) {
	byteBuffer := bytes.Buffer{}
	if _, err := byteBuffer.ReadFrom(buffer); err != nil {
		return nil, err
	}

	if err := buffer.Close(); err != nil {
		return nil, err
	}

	return byteBuffer.Bytes(), nil
}

func bufferToMap(buffer io.ReadCloser) (map[string]interface{}, error) {
	if b, err := readCloseableBuffer(buffer); err != nil {
		return nil, err
	} else {
		m := make(map[string]interface{})
		if jsonErr := json.Unmarshal(b, &m); jsonErr != nil {
			fmt.Println(jsonErr.Error())
			return nil, jsonErr
		} else {
			return m, nil
		}
	}
}

func jsonContentType(header *http.Header) bool {
	return strings.Contains(header.Get("Content-Type"), jsonMime)
}
