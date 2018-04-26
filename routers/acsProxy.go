package routers

import (
	"net/http/httputil"
	"net/url"
	"github.com/UHERO/rest-api/controllers"
	"github.com/gorilla/mux"
)

func SetAcsProxyRoute(
	router *mux.Router,
	// feedbackRepository *data.FeedbackRepository,
) *mux.Router {
	target := "https://api.census.gov/data/2016/acs/acs5/profile?get=DP02_0061PE,DP03_0062E,DP02_0064PE,DP02_0065PE,DP03_0009PE,DP03_0021PE,DP04_0005PE,DP04_0004PE,DP03_0025E,DP04_0134E,NAME&for=tract:*&in=state:15%20county:*&key=ad57a888cd72bea7153fa37026fca3dc19eb0134"
	remote, err := url.Parse(target)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	router.HandleFunc("/v1/acs", controllers.GetAcsData(proxy))
	return router
}

/* func handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request ){
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = mux.Vars(r)["rest"]
		p.ServeHTTP(w, r)
	}
} */
