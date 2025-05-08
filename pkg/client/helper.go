package client

import (
	"context"
	"fmt"
	"net/url"

	"golang.org/x/oauth2"
)

func getTokenSource(_ context.Context, bearerToken string) oauth2.TokenSource {
	return oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: bearerToken,
		TokenType:   "Bearer",
	})
}

func (e FluidTopicsAPIError) Message() string {
	return fmt.Sprintf("%s (HTTP %d): %s - %s", e.ErrorText, e.Status, e.MessageStr, e.Path)
}

func (e FluidTopicsAPIError) Error() string {
	return fmt.Sprintf("%s (HTTP %d): %s - %s", e.ErrorText, e.Status, e.MessageStr, e.Path)
}

type ReqOpt func(reqURL *url.URL)
