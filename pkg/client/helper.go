package client

import (
	"fmt"
	"net/url"

	"golang.org/x/oauth2"
)

func getTokenSource(bearerToken string) oauth2.TokenSource {
	return oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: bearerToken,
	})
}

func (e FluidTopicsAPIError) Message() string {
	return fmt.Sprintf("%s (HTTP %d): %s - %s", e.ErrorText, e.Status, e.MessageStr, e.Path)
}

func (e FluidTopicsAPIError) Error() string {
	return fmt.Sprintf("%s (HTTP %d): %s - %s", e.ErrorText, e.Status, e.MessageStr, e.Path)
}

type ReqOpt func(reqURL *url.URL)
