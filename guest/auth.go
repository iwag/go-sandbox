package main

import (
	"encoding/gob"
	"errors"
	"os"
	"net/http"
	"net/url"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"

	"google.golang.org/api/plus/v1"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

const (
	defaultSessionID = "default"
	googleProfileSessionKey = "google_profile"
	oauthTokenSessionKey    = "oauth_token"

	oauthFlowRedirectKey = "redirect"
)

var (
	OAuthConfig *oauth2.Config
	SessionStore sessions.Store
)

func init() {
	// Gob encoding for gorilla/sessions
	gob.Register(&oauth2.Token{})
	gob.Register(&Profile{})

	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")

	OAuthConfig = configureOAuthClient(clientId, clientSecret)

	cookieStore := sessions.NewCookieStore([]byte("something-very-secret"))
	cookieStore.Options = &sessions.Options{
		HttpOnly: true,
	}
	SessionStore = cookieStore

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/oauth2callback", oauthCallbackHandler)
}

// loginHandler initiates an OAuth flow to authenticate the user.
func loginHandler(w http.ResponseWriter, r *http.Request)  {
	sessionID := uuid.NewUUID().String()

	oauthFlowSession, err := SessionStore.New(r, sessionID)
	if err != nil {
		return // appErrorf(err, "could not create oauth session: %v", err)
	}
	oauthFlowSession.Options.MaxAge = 10 * 60 // 10 minutes

	redirectURL, err := validateRedirectURL(r.FormValue("redirect"))
	if err != nil {
		return // appErrorf(err, "invalid redirect URL: %v", err)
	}
	oauthFlowSession.Values[oauthFlowRedirectKey] = redirectURL

	if err := oauthFlowSession.Save(r, w); err != nil {
		return // appErrorf(err, "could not save session: %v", err)
	}

	// Use the session ID for the "state" parameter.
	// This protects against CSRF (cross-site request forgery).
	// See https://godoc.org/golang.org/x/oauth2#Config.AuthCodeURL for more detail.
	url := OAuthConfig.AuthCodeURL(sessionID, oauth2.ApprovalForce,	oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusFound)
	// return nil
}

// validateRedirectURL checks that the URL provided is valid.
// If the URL is missing, redirect the user to the application's root.
// The URL must not be absolute (i.e., the URL must refer to a path within this
// application).
func validateRedirectURL(path string) (string, error) {
	if path == "" {
		return "/", nil
	}

	// Ensure redirect URL is valid and not pointing to a different server.
	parsedURL, err := url.Parse(path)
	if err != nil {
		return "/", err
	}
	if parsedURL.IsAbs() {
		return "/", errors.New("URL must be absolute")
	}
	return path, nil
}

// oauthCallbackHandler completes the OAuth flow, retreives the user's profile
// information and stores it in a session.
func oauthCallbackHandler(w http.ResponseWriter, r *http.Request) /* *appError */ {
	oauthFlowSession, err := SessionStore.Get(r, r.FormValue("state"))
	ctx := appengine.NewContext(r)
	if err != nil {
		log.Infof(ctx, "invalid state parameter. try logging in again.")
		return // appErrorf(err, "invalid state parameter. try logging in again.")
	}

	redirectURL, ok := oauthFlowSession.Values[oauthFlowRedirectKey].(string)
	// Validate this callback request came from the app.
	if !ok {
		log.Infof(ctx, "invalid state parameter. try logging in again. 2")
		return // appErrorf(err, "invalid state parameter. try logging in again.")
	}

	code := r.FormValue("code")
	tok, err := OAuthConfig.Exchange(ctx, code)
	if err != nil {
		log.Infof(ctx, "could not get auth token: %v", err)
		return // appErrorf(err, "could not get auth token: %v", err)
	}

	session, err := SessionStore.New(r, defaultSessionID)
	if err != nil {
		log.Infof(ctx, "could not get default session: %v", err)
		return // appErrorf(err, "could not get default session: %v", err)
	}

	profile, err := fetchProfile(ctx, tok)
	if err != nil {
		log.Infof(ctx, "could not fetch Google profile: %v", err)
		return // appErrorf(err, "could not fetch Google profile: %v", err)
	}

	session.Values[oauthTokenSessionKey] = tok
	// Strip the profile to only the fields we need. Otherwise the struct is too big.
	stripped := stripProfile(profile)
	log.Infof(ctx, "success! ", stripped)

	session.Values[googleProfileSessionKey] = stripped
	if err := session.Save(r, w); err != nil {
		log.Infof(ctx, "could not save session: %v", err)
		return // appErrorf(err, "could not save session: %v", err)
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
	// return nil
}

// fetchProfile retrieves the Google+ profile of the user associated with the
// provided OAuth token.
func fetchProfile(ctx context.Context, tok *oauth2.Token) (*plus.Person, error) {
	client := oauth2.NewClient(ctx, OAuthConfig.TokenSource(ctx, tok))
	plusService, err := plus.New(client)
	if err != nil {
		return nil, err
	}
	return plusService.People.Get("me").Do()
}

// logoutHandler clears the default session.
func logoutHandler(w http.ResponseWriter, r *http.Request)  {
	session, err := SessionStore.New(r, defaultSessionID)
	if err != nil {
		return // appErrorf(err, "could not get default session: %v", err)
	}
	session.Options.MaxAge = -1 // Clear session.
	if err := session.Save(r, w); err != nil {
		return // appErrorf(err, "could not save session: %v", err)
	}
	redirectURL := r.FormValue("redirect")
	if redirectURL == "" {
		redirectURL = "/"
	}
	http.Redirect(w, r, redirectURL, http.StatusFound)
	// return nil
}

// profileFromSession retreives the Google+ profile from the default session.
// Returns nil if the profile cannot be retreived (e.g. user is logged out).
func profileFromSession(r *http.Request) *Profile {
	session, err := SessionStore.Get(r, defaultSessionID)
	if err != nil {
		return nil
	}
	tok, ok := session.Values[oauthTokenSessionKey].(*oauth2.Token)
	if !ok || !tok.Valid() {
		return nil
	}
	profile, ok := session.Values[googleProfileSessionKey].(*Profile)
	if !ok {
		return nil
	}
	return profile
}

type Profile struct {
	ID, DisplayName, ImageURL string
}

// stripProfile returns a subset of a plus.Person.
func stripProfile(p *plus.Person) *Profile {
	return &Profile{
		ID:          p.Id,
		DisplayName: p.DisplayName,
		ImageURL:    p.Image.Url,
	}
}

func configureOAuthClient(clientID, clientSecret string) *oauth2.Config {
	redirectURL := os.Getenv("OAUTH2_CALLBACK")
	if redirectURL == "" {
		redirectURL = "http://localhost:8080/oauth2callback"
	}
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}
}
