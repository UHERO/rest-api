package controllers

import (
	"github.com/uhero/rest-api/data"
	"net/http"
	"github.com/markbates/goth/gothic"
	"fmt"
	"html/template"
)

func Display(applicationRepository *data.ApplicationRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		applications, err := applicationRepository.GetAll(user.NickName)
		if err != nil {
			panic(err)
		}
		userResource := UserResource{User: user, Applications: applications}
		t, _ := template.New("userinfo").Parse(applicationTemplate)
		t.Execute(w, userResource)
	}
}

var applicationTemplate = `
<h1>Applications for {{.User.NickName}} ({{.User.UserID}})</h1>
{{range .Applications}}
<h2>{{.Name}}</h2>
<dl class="dl-horizontal">
	<dt>Hostname</dt>
	<dd>{{.Hostname}}</dd>
	<dt>API Key</dt>
	<dd>{{.APIKey}}</dd>
</dl>
{{else}}
No Applications
{{end}}
<button>Add Application</button>
`
