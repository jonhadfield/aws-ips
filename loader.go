package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/matryer/try"
)

func loadFromFile(filePath string) (ranges IPRangeDoc, err error) {
	var fileContent []byte
	fileContent, err = ioutil.ReadFile(filePath)
	if err != nil {
		return
	}
	err = json.Unmarshal(fileContent, &ranges)
	return
}

func loadFromURL() (ranges IPRangeDoc, err error) {
	var request *http.Request
	request, err = http.NewRequest(http.MethodGet, ipURL, nil)
	if err != nil {
		return
	}

	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 3,
			DisableKeepAlives:   false,
		},
		Timeout: time.Duration(6) * time.Second,
	}
	var resp *http.Response

	err = try.Do(func(attempt int) (bool, error) {
		var rErr error
		resp, rErr = client.Do(request)
		return attempt < 3, rErr
	})

	if err != nil {
		return
	}
	defer resp.Body.Close()

	var syncRespBodyBytes []byte
	syncRespBodyBytes, err = getResponseBody(resp)
	if err != nil {
		return
	}
	err = json.Unmarshal(syncRespBodyBytes, &ranges)
	return
}

func loadRanges(filePath string) (ranges IPRangeDoc, err error) {
	if filePath != "" {
		return loadFromFile(filePath)
	}
	return loadFromURL()

}

func getResponseBody(resp *http.Response) (body []byte, err error) {
	var output io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		output, err = gzip.NewReader(resp.Body)
		if err != nil {
			return
		}
		if err != nil {
			return
		}
	default:
		output = resp.Body
		if err != nil {
			return
		}
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(output)
	if err != nil {
		return
	}
	body = buf.Bytes()
	return
}
