package controllers

import (
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
		proxy.Transport = &data.CensusTransport{CacheRepository: cacheRepository}
		proxy.ServeHTTP(w, r)
	}
}
