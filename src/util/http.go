package util

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func MakePostRequest(url, contentType string, body any) (any, error) {
	marshaled, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	postBody := bytes.NewBuffer(marshaled)
	resp, err := http.Post(url, contentType, postBody)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var serviceResponse any
	err = json.Unmarshal(bytesResp, &serviceResponse)
	if err != nil {
		return nil, err
	}

	return &serviceResponse, nil
}
