package controllers

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/UHERO/rest-api/data"
	"github.com/gorilla/mux"
)

// GetCensusData retrieves data from api.census.gov
func GetCensusData(cacheRepository *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		target := "https://api.census.gov/" + mux.Vars(r)["census_endpoint"]
		remote, err := url.Parse(target)
		if err != nil {
			panic(err)
		}
		proxy := httputil.NewSingleHostReverseProxy(remote)
		r.URL.Path = ""
		proxy.Transport = &censusTransport{CacheRepository: cacheRepository}
		log.Print("request ", r)
		proxy.ServeHTTP(w, r)
		// response, err := http.DefaultTransport.RoundTrip(r)
		// body, err := ioutil.ReadAll(response.Body)
		// WriteCachePair(r, cacheRepository, body)
	}
}

type censusTransport struct {
	// CapturedTransport http.RoundTripper
	CacheRepository *data.CacheRepository
}

func (t *censusTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	response, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		log.Print("Error retrieving data from census.gov: ", err)
		log.Print("fail")
		log.Print("cache repo is nil", &data.CacheRepository{} == nil)
		log.Print("census req is nil", GetCensusReqURI(r))
		cachedVal, _ := t.CacheRepository.GetCache(GetCensusReqURI(r))
		if cachedVal != nil {
			t := &http.Response{
				Status:        "200 OK",
				StatusCode:    200,
				Proto:         "HTTP/1.1",
				ProtoMajor:    1,
				ProtoMinor:    1,
				Body:          ioutil.NopCloser(bytes.NewBuffer(cachedVal)),
				ContentLength: int64(len(cachedVal)),
				Request:       r,
				Header:        make(http.Header, 0),
			}
			log.Print("err", err)
			err = nil
			return t, err
		}
		log.Print("no response from acs")
		t := &http.Response{
			Status:        "200 OK",
			StatusCode:    200,
			Proto:         "HTTP/1.1",
			ProtoMajor:    1,
			ProtoMinor:    1,
			Body:          ioutil.NopCloser(bytes.NewBufferString("Error retrieving data")),
			ContentLength: int64(len("Error retrieving data")),
			Request:       r,
			Header:        make(http.Header, 0),
		}
		log.Print("custom response")
		err = nil
		return t, err
	}
	log.Print("test")
	body, err := httputil.DumpResponse(response, true)
	log.Print("body")
	if err != nil {
		log.Print("Error retrieving data from census.gov: ", err)
		t := &http.Response{
			Status:        "200 OK",
			StatusCode:    200,
			Proto:         "HTTP/1.1",
			ProtoMajor:    1,
			ProtoMinor:    1,
			Body:          ioutil.NopCloser(bytes.NewBufferString("Error retrieving data")),
			ContentLength: int64(len("Error retrieving data")),
			Request:       r,
			Header:        make(http.Header, 0),
		}
		err = nil
		return t, err
	}
	WriteCachePair(r, t.CacheRepository, body)
	return response, err
}
