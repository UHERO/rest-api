package controllers

import (
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
		proxy.ServeHTTP(w, r)
		response, err := http.DefaultTransport.RoundTrip(r)
		if err != nil {
			log.Printf("Error retrieving data from census.gov: ", err)
			cached_val, _ := cacheRepository.GetCache(GetCensusReqURI(r))
			if cached_val != nil {
				//log.Printf("DEBUG: Cache HIT: " + url)
				WriteResponse(w, cached_val)
				// w.Header().Set("Content-Type", "application/json")
				// w.Write(cached_val)
				return
			}
			return
		}
		//body, err := httputil.DumpResponse(response, true)
		body, err := ioutil.ReadAll(response.Body)
		WriteCachePair(r, cacheRepository, body)
	}
}
