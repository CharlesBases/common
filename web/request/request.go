package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type Response struct {
	Body map[string]interface{}
}

/*
 arg[0] booy

 arg[1] header
	dedault
	{
		"Content-type": "application/json"
	}
*/
func Request(method string, url string, args ...map[string]interface{}) (*Response, error) {
	client := new(http.Client)
	// body
	req, err := http.NewRequest(strings.ToUpper(method), url, func() io.Reader {
		for k := range args {
			if bs, err := json.Marshal(args[k]); err == nil {
				return bytes.NewBuffer(bs)
			}
		}
		return nil
	}())
	if err != nil {
		return nil, err
	}
	// header
	for i, j := range args {
		if i == 1 {
			for k, v := range j {
				req.Header.Set(k, fmt.Sprintf("%v", v))
			}
			break
		}
	}

	if resp, err := client.Do(req); err == nil && resp.StatusCode == 200 {
		defer resp.Body.Close()
		response := new(Response)
		if bs, err := ioutil.ReadAll(resp.Body); err != nil {
			return nil, err
		} else {
			/*
				var buff bytes.Buffer
				json.Indent(&buff, bs, "", "   ")
				response.Body = buff.String()
			*/
			json.Unmarshal(bs, &response.Body)
			return response, nil
		}
	}
	return nil, errors.New("request error")
}
