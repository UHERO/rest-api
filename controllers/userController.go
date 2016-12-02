package controllers

import (
	"fmt"
	"github.com/UHERO/rest-api/common"
	"github.com/markbates/goth/gothic"
	"net/http"
)

func ProviderCallback(provider string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		q.Add("provider", provider)
		r.URL.RawQuery = q.Encode()
		userProfile, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		// add user to user table if not already there

		// get

		// attach the user Id to the JWT
		token, err := common.GenerateJWT(nil, userProfile.Email, "dataPortalUser")
		common.StoreJWT(w, r, token)
		if err != nil {
			panic(err)
		}

		http.Redirect(w, r, "http://localhost:4200/?login", http.StatusMovedPermanently)
	}

}

func GetEmail(w http.ResponseWriter, r *http.Request) {

}
