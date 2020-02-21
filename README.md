# 2020-02-21_auth_code

This is the result from a stream talking about the OAuth2 Authorization Code flow.

What we covered:

1. There are 4 different flows in "OAuth2", all designed to delegate access to a user's resources (via "scopes") to a 3rd party application.
2. We implemented the "Authorization Code" flow, which is for secure (server to server) communication.
    1. Step 1: Redirect from your app to the oauth2 provider
        1. Step 1.1: The user interacts with the oauth2 provider to not provide personal details (such as a password) to our app. On authentication with the provider and authorizing our appp, they come back to our "redirect uri"
    2. Step 2: Handle the redirect _back_ from the provider and exchange the short-lived "code" (from the query params) for an access token
    3. Step 3: If you want to use any allowed resources (allowed via "scopes"), use the "access\_token" to access them.

## Usage

1. Create an OAuth application on [Github](https://github.com/settings/applications/new?oauth_application[callback_url]=http://localhost:8000/auth/github/callback&oauth_application[name]=alpinehq%20demo&oauth_application[url]=http://localhost:8000)
    1. Use `http://localhost:8000/auth/github/callback` as the Authorization callback URL
2. Export the provided client id and client secret in the environment: `GITHUB_CLIENT_ID` and `GITHUB_CLIENT_SECRET`
3. `go run main.go`
4. `open http://localhost:8000/auth/github`
5. Watch the terminal (we will make this more useful later)

# Resources

* [Github "Authorizing OAuth Apps"](https://developer.github.com/apps/building-oauth-apps/authorizing-oauth-apps/)
