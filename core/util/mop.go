package util

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var (
	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
)

func NodeUsabe(node string) (bool, error) {
	response, err := client.Get("http://" + node + ":1633")
	if err != nil {
		return false, err
	}
	defer response.Body.Close()
	return true, nil
}

func VoucherUsabe(node string, voucher string) (bool, error) {
	response, err := client.Get("http://" + node + ":1635" + "/stamps/" + voucher)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	bts, _ := ioutil.ReadAll(response.Body)
	var ret map[string]interface{}
	if err := json.Unmarshal(bts, &ret); err != nil {
		return false, err
	}
	return ret["usable"].(bool), nil
}

func ReferenceUsabe(gateway string, reference string) (bool, error) {
	response, err := client.Get(gateway + "/access/" + reference + "/")
	if err != nil {
		return false, err
	}
	defer response.Body.Close()
	return true, nil
}
