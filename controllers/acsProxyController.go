package controllers

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/UHERO/rest-api/data"
)

// GetAcsData retrieves 5-Year 2016 Data Profile for all counties and census tracts in the state of Hawaii
func GetAcsData(cacheRepository *data.CacheRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ids, ok := getAcsIdsList(w, r)
		if !ok {
			return
		}
		acsKey := os.Getenv("ACS_KEY")
		joinIds := strings.Join(ids, ",")
		target := "https://api.census.gov/data/2016/acs/acs5/profile?get=" + joinIds + "&for=tract:*&in=state:15%20county:*&key=" + acsKey
		remote, err := url.Parse(target)
		if err != nil {
			panic(err)
		}
		proxy := httputil.NewSingleHostReverseProxy(remote)
		r.URL.Path = ""
		proxy.ServeHTTP(w, r)
		response, err := http.DefaultTransport.RoundTrip(r)
		if err != nil {
			log.Printf("Error retrieving data from ACS: ", err)
			cached_val, _ := cacheRepository.GetCache(GetFullRelativeURL(r))
			if cached_val != nil {
				//log.Printf("DEBUG: Cache HIT: " + url)
				WriteResponse(w, cached_val)
				return
			}
			return
		}
		body, err := httputil.DumpResponse(response, true)
		WriteCache(r, cacheRepository, body)
	}
}
