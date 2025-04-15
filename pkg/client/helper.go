package client

import (
	"context"
	"net/url"
	"strconv"

	"golang.org/x/oauth2"
)

const ItemsPerPage = 5

type PageOptions struct {
	Next     string
	PageSize int
}

func getTokenSource(_ context.Context, bearerToken string) oauth2.TokenSource {
	return oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: bearerToken,
		TokenType:   "Bearer",
	})
}

func WithPageLimit(pageSize int) ReqOpt {
	if pageSize != 0 {
		return WithQueryParam("per_page", strconv.Itoa(pageSize))
	}
	pageSize = ItemsPerPage
	return WithQueryParam("per_page", strconv.Itoa(pageSize))
}

func WithPageCursor(nextPageToken string) ReqOpt {
	return WithQueryParam("cursor", nextPageToken)
}

func WithQueryParam(key string, value string) ReqOpt {
	return func(reqURL *url.URL) {
		if value != "" {
			q := reqURL.Query()
			q.Set(key, value)
			reqURL.RawQuery = q.Encode()
		}
	}
}

type ReqOpt func(reqURL *url.URL)
