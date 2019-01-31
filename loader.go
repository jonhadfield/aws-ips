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

func loadRanges(filePath string) (ranges IPRangeDoc, err error) {
	if filePath != "" {
		var fileContent []byte
		fileContent, err = ioutil.ReadFile(filePath)
		if err != nil {
			return
		}
		err = json.Unmarshal(fileContent, &ranges)
		if err != nil {
			return
		}
	} else {
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

		// process response body
		var syncRespBodyBytes []byte
		syncRespBodyBytes, err = getResponseBody(resp)
		if err != nil {
			return
		}
		err = json.Unmarshal(syncRespBodyBytes, &ranges)
		if err != nil {
			return
		}
	}
	return ranges, err

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
