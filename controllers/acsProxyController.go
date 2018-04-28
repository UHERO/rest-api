package controllers

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

// GetAcsData retrieves 5-Year 2016 Data Profile for all counties and census tracts in the state of Hawaii
func GetAcsData() func(http.ResponseWriter, *http.Request) {
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
	}
}
