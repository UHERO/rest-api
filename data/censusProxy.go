package data

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

type CensusTransport struct {
	CacheRepository *CacheRepository
}

func (t *CensusTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	response, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		return RetrieveCached(t, r)
	}
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Print("Error parsing response body", GetCensusReqURI(r))
		return SetErrorMsg(r)
	}
	err = response.Body.Close()
	if err != nil {
		return nil, err
	}
	response = SetCensusResponse(bytes.NewBuffer(b), "application/json")
	WriteCachePair(r, t.CacheRepository, b)
	return response, err
}

func SetCensusResponse(body *bytes.Buffer, contentType string) *http.Response {
	resp := &http.Response{}
	resp.StatusCode = 200
	resp.Body = ioutil.NopCloser(body)
	resp.Header = make(http.Header, 0)
	resp.Header.Set("Content-Type", contentType)
	return resp
}

func GetCensusReqURI(r *http.Request) string {
	return r.RequestURI
}

func WriteCachePair(r *http.Request, c *CacheRepository, payload []byte) {
	url := GetCensusReqURI(r)
	err := c.SetCachePair(url, payload)
	if err != nil {
		log.Printf("Cache store FAILURE: %s", url)
		return
	}
}

func RetrieveCached(t *CensusTransport, r *http.Request) (*http.Response, error) {
	// Check cache
	cachedVal, _ := t.CacheRepository.GetCache(GetCensusReqURI(r))
	if cachedVal != nil {
		resp := SetCensusResponse(bytes.NewBuffer(cachedVal), "application/json")
		return resp, nil
	}
	// If cache is empty, return error message
	return SetErrorMsg(r)
}

func SetErrorMsg(r *http.Request) (*http.Response, error) {
	resp := SetCensusResponse(bytes.NewBufferString("Error retrieving data"), "plain/text")
	resp.StatusCode = 404
	return resp, nil
}
