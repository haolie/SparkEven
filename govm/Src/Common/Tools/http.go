package Tools

import (
	"io"
	"net/http"
)

func HttpRequest(url string, method string, headers map[string]string) (data []byte, err error) {
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return
	}

	if headers != nil {
		for key, value := range headers {
			request.Header.Add(key, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return
	}

	data, err = io.ReadAll(resp.Body)
	return
}
