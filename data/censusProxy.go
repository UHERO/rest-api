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
		// Check cache
		cachedVal, _ := t.CacheRepository.GetCache(GetCensusReqURI(r))
		if cachedVal != nil {
			rBody := bytes.NewBuffer(cachedVal)
			resp := SetCensusResponse(rBody)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		}
		// If cache is empty, return error message
		rBody := bytes.NewBufferString("Error retrieving data")
		resp := SetCensusResponse(rBody)
		resp.Header.Set("Content-Type", "plain/text")
		return resp, nil
	}
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Print("Error parsing response body")
		rBody := bytes.NewBufferString("Error retrieving data")
		resp := SetCensusResponse(rBody)
		resp.Header.Set("Content-Type", "plain/text")
		return nil, err
	}
	err = response.Body.Close()
	if err != nil {
		return nil, err
	}
	body := ioutil.NopCloser(bytes.NewReader(b))
	response.Body = body
	WriteCachePair(r, t.CacheRepository, b)
	return response, err
}

func SetCensusResponse(body *bytes.Buffer) *http.Response {
	resp := &http.Response{}
	resp.StatusCode = 200
	resp.Body = ioutil.NopCloser(body)
	resp.Header = make(http.Header, 0)
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
	//log.Printf("DEBUG: Stored in cache: %s", url)
}
