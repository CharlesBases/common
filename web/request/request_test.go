package request

import (
	"testing"

	"charlesbases/common/log"
)

var (
	url = "http://www.baidu.com"
)

func TestRequest(t *testing.T) {
	defer log.Flush()

	req1, _ := Request("GET", url)
	log.Debug(req1.Body)

	req2, _ := Request("POST", url, map[string]interface{}{"user": "赵铁柱"})
	log.Debug(req2.Body)

	req3, _ := Request("POST", url, map[string]interface{}{"user": "赵铁柱"}, map[string]interface{}{"Content-type": "application/json"})
	log.Debug(req3.Body)

}
