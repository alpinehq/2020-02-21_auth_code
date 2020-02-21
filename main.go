package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	// 1. Redirect from your app to the authorization provider
	//   1.1. The user now interacts with the provider directly
	//     YES: go to 2
	//     NO: go to 2, but with a query param "error"
	// 2. Back in our app, we get either a "error" or "code"
	//    query param. We then have to exchange that with the
	//    provider for more info.
	// 3. Exchanging the code for an access token.
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//     fmt.Printf("%+v\n", r)
	// })

	// Step 1: redirect a user to here, to get to the AP
	http.HandleFunc("/auth/github", func(w http.ResponseWriter, r *http.Request) {
		vals := url.Values{
			// client_id=first&client_id=second
			"client_id":    []string{os.Getenv("GITHUB_CLIENT_ID")},
			"redirect_uri": []string{"http://localhost:8000/auth/github/callback"},
			"scope": []string{
				strings.Join([]string{
					"read:user",
					"user:email",
				}, " "),
			},
		}
		http.Redirect(w, r, "https://github.com/login/oauth/authorize?"+vals.Encode(), http.StatusTemporaryRedirect)
	})

	// Step 2: handle ?error= OR ?code=
	http.HandleFunc("/auth/github/callback", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if err := q.Get("error"); err != "" {
			fmt.Printf("ERROR LOGGING IN: %s\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		vals := url.Values{
			"client_id":     []string{os.Getenv("GITHUB_CLIENT_ID")},
			"client_secret": []string{os.Getenv("GITHUB_CLIENT_SECRET")},
			"code":          []string{q.Get("code")},
		}
		req, err := http.NewRequest(
			http.MethodPost,
			"https://github.com/login/oauth/access_token",
			strings.NewReader(vals.Encode()),
		)
		if err != nil {
			fmt.Printf("unable to generate request: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		req.Header.Set("Accept", "application/json")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("unable to generate request: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()

		var accessToken struct {
			AccessToken string `json:"access_token"`
			Type        string `json:"type"`
			Scope       string `json:"scope"`
		}

		raw, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(raw, &accessToken)
		fmt.Printf("access token: %+v\n", accessToken)

		req, _ = http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
		req.Header.Set("Authorization", "token "+accessToken.AccessToken)
		res, err = http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("unable to get user info: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()

		raw, _ = ioutil.ReadAll(res.Body)
		fmt.Printf("\n\n%q\n", raw)

		var githubUser struct {
			Email string `json:"email"`
			Login string `json:"login"`
			ID    int64  `json:"id"`
		}

		err = json.Unmarshal(raw, &githubUser)
		if err != nil {
			fmt.Printf("err: %s\n", err)
		}

		fmt.Printf("github user: %+v\n", githubUser)
	})

	fmt.Printf("Serving...\n")
	http.ListenAndServe(":8000", nil)
}
