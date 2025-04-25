package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/ratelimit"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"golang.org/x/oauth2"
)

const (
	baseURL          = "https://powin-staging.fluidtopics.net/api"
	getUsers         = "/users"
	getUserRolesById = "/users/%s/roles"
	getUserUsage     = "/users/%s/dump"
)

type FluidTopicsClient struct {
	wrapper     *uhttp.BaseHttpClient
	TokenSource oauth2.TokenSource
}

func New(ctx context.Context, bearerToken string) (*FluidTopicsClient, error) {
	httpClient, err := uhttp.NewClient(ctx, uhttp.WithLogger(true, ctxzap.Extract(ctx)))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP client: %w", err)
	}
	cli, err := uhttp.NewBaseHttpClientWithContext(context.Background(), httpClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create base HTTP client: %w", err)
	}
	client := FluidTopicsClient{
		wrapper:     cli,
		TokenSource: getTokenSource(ctx, bearerToken),
	}
	return &client, nil
}

func (c *FluidTopicsClient) ListUsers(ctx context.Context, _ *pagination.Token) ([]UserList, string, annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)
	var res []UserList
	var annotation annotations.Annotations

	queryUrl, err := url.JoinPath(baseURL, getUsers)
	if err != nil {
		l.Error(fmt.Sprintf("Error creating UserResponse URL: %s", err))
		return nil, "", nil, err
	}

	annotation, err = c.getResourcesFromAPI(ctx, queryUrl, &res)
	if err != nil {
		l.Error(fmt.Sprintf("Error getting resources: %s", err))
		return nil, "", nil, err
	}

	return res, "", annotation, nil
}

func (c *FluidTopicsClient) GetUserUsage(ctx context.Context, userID string) (UserUsage, annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)
	var res UserUsageResponse
	var annotation annotations.Annotations

	queryUrl, err := url.JoinPath(baseURL, fmt.Sprintf(getUserUsage, userID))
	if err != nil {
		l.Error(fmt.Sprintf("Error creating URL: %s", err))
		return res.User, nil, err
	}

	annotation, err = c.getResourcesFromAPI(ctx, queryUrl, &res)
	if err != nil {
		l.Error(fmt.Sprintf("Error getting resource: %s", err))
		return res.User, nil, err
	}

	return res.User, annotation, nil
}

func (c *FluidTopicsClient) ListRoles(ctx context.Context, options PageOptions, _ *pagination.Token) ([]UserRole, string, annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)
	var res UserRoleResponse
	var annotation annotations.Annotations

	queryUrl, err := url.JoinPath(baseURL, getUsers)
	if err != nil {
		l.Error(fmt.Sprintf("Error creating UserResponse URL: %s", err))
		return nil, "", nil, err
	}

	annotation, err = c.getResourcesFromAPI(ctx, queryUrl, &res, WithPageCursor(options.Next), WithPageLimit(options.PageSize))
	if err != nil {
		l.Error(fmt.Sprintf("Error getting resources: %s", err))
		return nil, "", nil, err
	}

	return res.User, "", annotation, nil
}

func (c *FluidTopicsClient) GetRolesByUserID(ctx context.Context, userID string) (UserRole, annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)
	var user UserRole
	var annotation annotations.Annotations

	queryUrl, err := url.JoinPath(baseURL, fmt.Sprintf(getUserRolesById, userID))
	if err != nil {
		l.Error(fmt.Sprintf("Error creating URL: %s", err))
		return user, nil, err
	}

	annotation, err = c.getResourcesFromAPI(ctx, queryUrl, &user)
	if err != nil {
		l.Error(fmt.Sprintf("Error getting resource: %s", err))
		return user, nil, err
	}

	return user, annotation, nil
}

func (c *FluidTopicsClient) getResourcesFromAPI(
	ctx context.Context,
	urlAddress string,
	res any,
	reqOptions ...ReqOpt,
) (annotations.Annotations, error) {
	_, annotation, err := c.doRequest(ctx, http.MethodGet, urlAddress, &res, reqOptions...)

	if err != nil {
		return nil, err
	}

	return annotation, nil
}

func (c *FluidTopicsClient) doRequest(
	ctx context.Context,
	method string,
	endpointUrl string,
	res interface{},
	reqOptions ...ReqOpt,
) (http.Header, annotations.Annotations, error) {
	var (
		resp *http.Response
		err  error
	)

	urlAddress, err := url.Parse(endpointUrl)

	if err != nil {
		return nil, nil, err
	}

	for _, o := range reqOptions {
		o(urlAddress)
	}

	authToken, err := c.TokenSource.Token()
	if err != nil {
		return nil, nil, err
	}

	req, err := c.wrapper.NewRequest(
		ctx,
		method,
		urlAddress,
		uhttp.WithContentTypeJSONHeader(),
		uhttp.WithAcceptJSONHeader(),
	)
	authToken.SetAuthHeader(req)

	if err != nil {
		return nil, nil, err
	}

	switch method {
	case http.MethodGet, http.MethodPut, http.MethodPost:
		var doOptions []uhttp.DoOption
		if res != nil {
			doOptions = append(doOptions, uhttp.WithResponse(&res))
		}
		resp, err = c.wrapper.Do(req, doOptions...)
		if resp != nil {
			defer resp.Body.Close()
		}
	case http.MethodDelete:
		resp, err = c.wrapper.Do(req)
		if resp != nil {
			defer resp.Body.Close()
		}
	}

	if err != nil {
		return nil, nil, err
	}

	annotation := annotations.Annotations{}
	if resp != nil {
		if desc, err := ratelimit.ExtractRateLimitData(resp.StatusCode, &resp.Header); err == nil {
			annotation.WithRateLimiting(desc)
		} else {
			return nil, annotation, err
		}

		return resp.Header, annotation, nil
	}

	return nil, nil, err
}
