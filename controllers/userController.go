package controllers

import (
	"net/http"
	"github.com/markbates/goth/gothic"
	"fmt"
	"github.com/UHERO/rest-api/common"
)

func ProviderCallback(provider string) func(http.ResponseWriter, *http.Request) {
	return func (w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		q.Add("provider", provider)
		r.URL.RawQuery = q.Encode()
		userProfile, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		token, err := common.GenerateJWT(userProfile.Email, "dataPortalUser")
		common.StoreJWT(w, r, token)
		if err != nil {
			panic(err)
		}

		http.Redirect(w, r, "http://localhost:4200/?login", http.StatusMovedPermanently)
	}

}

func GetEmail(w http.ResponseWriter, r *http.Request) {


}
