package main

import (
	"bytes"
	"io"
	"io/ioutil"
)

var (
	// backend endpoint => map[request_query]response
	// req/resp in json
	cache = make(map[string]map[string]string)
)

func AsJson(reader *io.ReadCloser) (string, error) {
	var buffer bytes.Buffer
	tmpReader := io.TeeReader(*reader, &buffer)
	jsonReq, err := ioutil.ReadAll(tmpReader)
	if err != nil {
		return "", err
	}
	rcloser := ioutil.NopCloser(bytes.NewReader(buffer.Bytes()))
	reader = &rcloser
	jsonString := string(jsonReq)
	return jsonString, nil
}

func UpdateCache(endpoint string, jsonReq string, jsonVal string) bool {
	if _, ok := cache[endpoint]; !ok {
		cache[endpoint] = make(map[string]string)
	}
	cache[endpoint][jsonReq] = jsonVal
	return true
}

func GetCachedValue(endpoint string, jsonString string) (string, bool) {
	epCache, ok := cache[endpoint]
	if !ok {
		return "", false
	}
	cachedReq, exists := epCache[jsonString]
	if !exists {
		return "", false
	}
	return cachedReq, true
}
