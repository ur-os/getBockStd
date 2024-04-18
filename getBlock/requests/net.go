package requests

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func doNodeRequest(body interface{}, nodeEndpoint string) (*http.Response, error) {
	requestBodyBytes, err := json.Marshal(&body)
	if err != nil {
		return nil, err
	}
	requestBodyBytesReader := bytes.NewReader(requestBodyBytes)

	req, err := http.NewRequest(http.MethodPost, nodeEndpoint, requestBodyBytesReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", strconv.Itoa(len(requestBodyBytes)))
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "keep-alive")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(strconv.Itoa(resp.StatusCode))
	}

	return resp, nil
}
