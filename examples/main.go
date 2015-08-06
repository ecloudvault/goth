package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/pat"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	//"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/amazon"
	//"github.com/markbates/goth/providers/gplus"
	//"github.com/markbates/goth/providers/lastfm"
	//"github.com/markbates/goth/providers/linkedin"
	//"github.com/markbates/goth/providers/spotify"
	//"github.com/markbates/goth/providers/twitter"
	)

func main() {
	goth.UseProviders(
		//twitter.New(os.Getenv("TWITTER_KEY"), os.Getenv("TWITTER_SECRET"), "http://localhost:3000/auth/twitter/callback"),
		//facebook.New(os.Getenv("FACEBOOK_KEY"), os.Getenv("FACEBOOK_SECRET"), "http://localhost:3000/auth/facebook/callback"),
		//gplus.New(os.Getenv("GPLUS_KEY"), os.Getenv("GPLUS_SECRET"), "http://localhost:3000/auth/gplus/callback"),
		//github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), "http://localhost:3000/auth/github/callback"),
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), "https://ecloud.nimbostrati.com:9898/auth/github/callback"),
    amazon.New(os.Getenv("AMAZON_KEY"), os.Getenv("AMAZON_SECRET"), "https://ecloud.nimbostrati.com:9898/auth/amazon/callback", "profile"),

		//spotify.New(os.Getenv("SPOTIFY_KEY"), os.Getenv("SPOTIFY_SECRET"), "http://localhost:3000/auth/spotify/callback"),
		//linkedin.New(os.Getenv("LINKEDIN_KEY"), os.Getenv("LINKEDIN_SECRET"), "http://localhost:3000/auth/linkedin/callback"),
		//lastfm.New(os.Getenv("LASTFM_KEY"), os.Getenv("LASTFM_SECRET"), "http://localhost:3000/auth/lastfm/callback"),
	)

	// Assign the GetState function variable so we can return the
	// state string we want to get back at the end of the oauth process.
	// Only works with facebook and gplus providers.
	gothic.GetState = func(req *http.Request) string {
		// Get the state string from the query parameters.
		return req.URL.RawQuery //.Query() //.Get("state")
	}

	p := pat.New()
	p.Get("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {

		// print our state string to the console
		fmt.Println("State: " + gothic.GetState(req))

		user, err := gothic.CompleteUserAuth(res, req)
		if err != nil {
			fmt.Fprintln(res, err)
			return
		}
		t, _ := template.New("foo").Parse(userTemplate)
		t.Execute(res, user)
	})

	p.Get("/auth/{provider}", gothic.BeginAuthHandler)
	p.Get("/", func(res http.ResponseWriter, req *http.Request) {
		t, _ := template.New("foo").Parse(indexTemplate)
		t.Execute(res, nil)
	})
	//http.ListenAndServe(":9898", p)
  http.ListenAndServeTLS(":9898", os.Getenv("GOTH_SSL_CERT"), os.Getenv("GOTH_SSL_KEY"), p) 
}

var indexTemplate = `
<p><a href="/auth/github">Log in with Github</a></p>
<p><a href="/auth/amazon">Log in with Amazon</a></p>
`

var userTemplate = `
<p>Name: {{.Name}}</p>
<p>Email: {{.Email}}</p>
<p>NickName: {{.NickName}}</p>
<p>Location: {{.Location}}</p>
<p>AvatarURL: {{.AvatarURL}} <img src="{{.AvatarURL}}"></p>
<p>Description: {{.Description}}</p>
<p>UserID: {{.UserID}}</p>
<p>AccessToken: {{.AccessToken}}</p>
`
