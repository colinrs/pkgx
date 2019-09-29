package httpclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// PostASJSON ...
func PostASJSON(url string, timeout time.Duration, v interface{}, headers map[string]string) (response []byte, code int, err error) {
	var bs []byte
	bs, err = json.Marshal(v)
	if err != nil {
		return
	}

	bf := bytes.NewBuffer(bs)

	client := http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("POST", url, bf)
	req.Header.Set("Content-Type", "application/json")

	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return
	}

	code = resp.StatusCode

	if resp.Body != nil {
		defer resp.Body.Close()
		response, err = ioutil.ReadAll(resp.Body)
	}

	return
}

// GetURLAsJSON ...
func GetURLAsJSON(url string, res interface{}) {

}

// GetURLAsText ...
func GetURLAsText(url string) (res string) {

	return
}

// GetURLAsBytes ...
func GetURLAsBytes(url string) (res []byte) {

	return
}
