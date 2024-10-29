package providers

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"golang.org/x/oauth2"
)

const (
	oauth2TokenKey = "oauth2Token"
	idTokenKey     = "idToken"
	userInfoKey    = "userInfo"

	StateQueryParameter = "state"
	CodeQueryParameter  = "code"

	idTokenExtraKey = "id_token"
)

type IdentityProviderChildrenProps struct {
	IDToken  string
	UserInfo oidc.UserInfo

	Logout func(ctx app.Context)

	Error   error
	Recover func(ctx app.Context)
}

type IdentityProvider struct {
	app.Compo

	Issuer        string
	ClientID      string
	RedirectURL   string
	HomeURL       string
	Scopes        []string
	StoragePrefix string
	Children      func(IdentityProviderChildrenProps) app.UI

	oauth2Token oauth2.Token
	idToken     string
	userInfo    oidc.UserInfo

	err error
}

func (c *IdentityProvider) Render() app.UI {
	return c.Children(
		IdentityProviderChildrenProps{
			IDToken:  c.idToken,
			UserInfo: c.userInfo,

			Logout: func(ctx app.Context) {
				c.logout(true, ctx)
			},

			Error:   c.err,
			Recover: c.recover,
		},
	)
}

func (c *IdentityProvider) OnMount(ctx app.Context) {
	// Only continue if there is no error state; this prevents endless loops
	if c.err == nil {
		c.authorize(ctx)
	}
}

func (c *IdentityProvider) OnNav(ctx app.Context) {
	// Only continue if there is no error state; this prevents endless loops
	if c.err == nil {
		c.authorize(ctx)
	}
}

func (c *IdentityProvider) panic(err error, ctx app.Context) {
	go func() {
		// Set the error
		c.err = err

		// Prevent infinite retries
		time.Sleep(time.Second)

		// Unset the error & enable re-trying
		c.err = err
	}()
}

func (c *IdentityProvider) recover(ctx app.Context) {
	// Clear the error
	c.err = nil

	// Logout
	c.logout(false, ctx)
}

func (c *IdentityProvider) watch(ctx app.Context) {
	for {
		// Wait till token expires
		if c.oauth2Token.Expiry.After(time.Now()) {
			time.Sleep(c.oauth2Token.Expiry.Sub(time.Now()))
		}

		// Fetch new OAuth2 token
		oauth2Token, err := oauth2.StaticTokenSource(&c.oauth2Token).Token()
		if err != nil {
			c.panic(err, ctx)

			return
		}

		// Parse ID token
		idToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			c.panic(err, ctx)

			return
		}

		// Persist state in storage
		if err := c.persist(*oauth2Token, idToken, c.userInfo, ctx); err != nil {
			c.panic(err, ctx)

			return
		}

		// Set the login state
		c.oauth2Token = *oauth2Token
		c.idToken = idToken
	}
}

func (c *IdentityProvider) logout(withRedirect bool, ctx app.Context) {
	// Remove from storage
	c.clear(ctx)

	// Reload the app
	if withRedirect {
		ctx.Reload()
	}
}

func (c *IdentityProvider) rehydrate(ctx app.Context) (oauth2.Token, string, oidc.UserInfo, error) {
	// Read state from storage
	oauth2Token := oauth2.Token{}
	idToken := ""
	userInfo := oidc.UserInfo{}

	if err := ctx.LocalStorage().Get(c.getKey(oauth2TokenKey), &oauth2Token); err != nil {
		return oauth2.Token{}, "", oidc.UserInfo{}, err
	}
	if err := ctx.LocalStorage().Get(c.getKey(idTokenKey), &idToken); err != nil {
		return oauth2.Token{}, "", oidc.UserInfo{}, err
	}
	if err := ctx.LocalStorage().Get(c.getKey(userInfoKey), &userInfo); err != nil {
		return oauth2.Token{}, "", oidc.UserInfo{}, err
	}

	return oauth2Token, idToken, userInfo, nil
}

func (c *IdentityProvider) persist(oauth2Token oauth2.Token, idToken string, userInfo oidc.UserInfo, ctx app.Context) error {
	// Write state to storage
	if err := ctx.LocalStorage().Set(c.getKey(oauth2TokenKey), oauth2Token); err != nil {
		return err
	}
	if err := ctx.LocalStorage().Set(c.getKey(idTokenKey), idToken); err != nil {
		return err
	}
	return ctx.LocalStorage().Set(c.getKey(userInfoKey), userInfo)
}

func (c *IdentityProvider) clear(ctx app.Context) {
	// Remove from storage
	ctx.LocalStorage().Del(c.getKey(oauth2TokenKey))
	ctx.LocalStorage().Del(c.getKey(idTokenKey))
	ctx.LocalStorage().Del(c.getKey(userInfoKey))

	// Remove cookies
	app.Window().Get("document").Set("cookie", "")
}

func (c *IdentityProvider) getKey(key string) string {
	// Get a prefixed key
	return fmt.Sprintf("%v.%v", c.StoragePrefix, key)
}

func (c *IdentityProvider) authorize(ctx app.Context) {
	// Read state from storage
	oauth2Token, idToken, userInfo, err := c.rehydrate(ctx)
	if err != nil {
		c.panic(err, ctx)

		return
	}

	// Create the OIDC provider
	provider, err := oidc.NewProvider(context.Background(), c.Issuer)
	if err != nil {
		c.panic(err, ctx)

		return
	}

	// Create the OAuth2 config
	config := &oauth2.Config{
		ClientID:    c.ClientID,
		RedirectURL: c.RedirectURL,
		Endpoint:    provider.Endpoint(),
		Scopes:      append([]string{oidc.ScopeOpenID}, c.Scopes...),
	}

	// Log in
	if oauth2Token.AccessToken == "" || userInfo.Email == "" {
		// Logged out state, info neither in storage nor in URL: Redirect to login
		if app.Window().URL().Query().Get(StateQueryParameter) == "" {
			ctx.Navigate(config.AuthCodeURL(c.RedirectURL, oauth2.AccessTypeOffline))

			return
		}

		// Intermediate state, info is in URL: Parse OAuth2 token
		oauth2Token, err := config.Exchange(context.Background(), app.Window().URL().Query().Get(CodeQueryParameter))
		if err != nil {
			c.panic(err, ctx)

			return
		}

		// Parse ID token
		idToken, ok := oauth2Token.Extra(idTokenExtraKey).(string)
		if !ok {
			c.panic(err, ctx)

			return
		}

		// Parse user info
		userInfo, err := provider.UserInfo(context.Background(), oauth2.StaticTokenSource(oauth2Token))
		if err != nil {
			c.panic(err, ctx)

			return
		}

		// Persist state in storage
		if err := c.persist(*oauth2Token, idToken, *userInfo, ctx); err != nil {
			c.panic(err, ctx)

			return
		}

		// Test validity of storage
		if _, _, _, err = c.rehydrate(ctx); err != nil {
			c.panic(err, ctx)

			return
		}

		// Update and navigate to home URL
		ctx.Navigate(c.HomeURL)

		return
	}

	// Validation state

	// Create the OIDC config
	oidcConfig := &oidc.Config{
		ClientID: c.ClientID,
	}

	// Create the OIDC verifier and validate the token (i.e. check for it's expiry date)
	verifier := provider.Verifier(oidcConfig)
	if _, err := verifier.Verify(context.Background(), idToken); err != nil {
		// Invalid token; clear and re-authorize
		c.clear(ctx)
		c.authorize(ctx)

		return
	}

	// Logged in state

	// Set the login state
	c.oauth2Token = oauth2Token
	c.idToken = idToken
	c.userInfo = userInfo

	// Watch and renew token once expired
	go c.watch(ctx)
}
