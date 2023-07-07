package cachelet

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"k8s.io/apimachinery/pkg/util/json"
	"net/http"
)

func sendConnectionRequest(endpoint string) bool {
	url := fmt.Sprintf("%s/peers", endpoint)
	req := &PeerConnectRequest{endpoint: endpoint}
	res := &PeerConnectResponse{}
	requestTo(http.MethodPost, url, req, res)
	if res.message != "ok" {
		return false
	} else {
		return true
	}
}

func requestTo(url, method string, reqPayload interface{}, resPayload interface{}) interface{} {
	client := &http.Client{}
	var req *http.Request
	var err error
	var reqBody io.Reader

	if reqPayload != nil {
		data, err := json.Marshal(reqPayload)
		if err != nil {
			log.Error(err.Error())
			return nil
		}
		reqBody = bytes.NewReader(data)
	}

	req, err = http.NewRequest(method, url, reqBody)

	if err != nil {
		log.Error(err.Error())
		return nil
	}

	res, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		log.Error(err.Error())
		return nil
	}

	err = json.Unmarshal(body, resPayload)

	if err != nil {
		log.Error(err.Error())
		return nil
	}
	return nil
}
