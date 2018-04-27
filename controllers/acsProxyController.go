package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/UHERO/rest-api/common"
	"github.com/gorilla/mux"
)

func GetAcsData() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ids, ok := getAcsIdsList(w, r)
		if !ok {
			return
		}
		joinIds := strings.Join(ids, ",")
		// target := "https://api.census.gov/data/2016/acs/acs5/profile?get=DP02_0061PE,DP03_0062E,DP02_0064PE,DP02_0065PE,DP03_0009PE,DP03_0021PE,DP04_0005PE,DP04_0004PE,DP03_0025E,DP04_0134E,NAME&for=tract:*&in=state:15%20county:*&key=ad57a888cd72bea7153fa37026fca3dc19eb0134"
		target := "https://api.census.gov/data/2016/acs/acs5/profile?get=" + joinIds + "&for=tract:*&in=state:15%20county:*&key=ad57a888cd72bea7153fa37026fca3dc19eb0134"
		remote, err := url.Parse(target)
		if err != nil {
			panic(err)
		}
		proxy := httputil.NewSingleHostReverseProxy(remote)
		// r.URL.Path = ""
		fmt.Print(r)

		proxy.ServeHTTP(w, r)
	}
}

func getAcsIdsList(w http.ResponseWriter, r *http.Request) (ids []string, ok bool) {
	ok = true
	idsList, gotIds := mux.Vars(r)["ids_list"]
	fmt.Print(mux.Vars(r))
	if !gotIds {
		common.DisplayAppError(w, errors.New("Couldn't get id from request"), "Bad request.", 400)
		ok = false
		return
	}
	idStrArr := strings.Split(idsList, ",")
	for _, idStr := range idStrArr {
		id := idStr
		ids = append(ids, id)
	}
	return
}
