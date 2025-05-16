package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

const (
	getUsers              = "/users"
	getUserRolesById      = "/users/%s/roles"
	getUserInfoById       = "/users/%s/dump"
	getAuthenticationInfo = "/authentication/current-session"
	createUser            = "/users/register"
)

type FluidTopicsClient struct {
	httpClient  *uhttp.BaseHttpClient
	tokenSource oauth2.TokenSource
	baseURL     string
}

func New(ctx context.Context, bearerToken string, domain string) (*FluidTopicsClient, error) {
	if !strings.HasPrefix(domain, "https://") {
		return nil, fmt.Errorf("domain must start with http://")
	}
	domain = strings.TrimRight(domain, "/")
	baseURL := fmt.Sprintf("%s/api", domain)

	httpClient, err := uhttp.NewClient(ctx, uhttp.WithLogger(true, ctxzap.Extract(ctx)))

	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP client: %w", err)
	}

	cli, err := uhttp.NewBaseHttpClientWithContext(context.Background(), httpClient)

	if err != nil {
		return nil, fmt.Errorf("failed to create base HTTP client: %w", err)
	}

	client := FluidTopicsClient{
		httpClient:  cli,
		tokenSource: getTokenSource(bearerToken),
		baseURL:     baseURL,
	}
	return &client, nil
}

func (c *FluidTopicsClient) ListUsers(ctx context.Context) ([]User, string, annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)
	var res []User
	var annotation annotations.Annotations

	queryUrl, err := url.JoinPath(c.baseURL, getUsers)
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

func (c *FluidTopicsClient) GetUserDetails(ctx context.Context, userID string) (User, annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)
	var res UserDataResponse
	var annotation annotations.Annotations

	queryUrl, err := url.JoinPath(c.baseURL, fmt.Sprintf(getUserInfoById, userID))
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

func (c *FluidTopicsClient) GetAuthenticationInfo(ctx context.Context) (AuthenticationInfo, annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)
	var res AuthenticationInfo
	var annotation annotations.Annotations

	queryUrl, err := url.JoinPath(c.baseURL, getAuthenticationInfo)
	if err != nil {
		l.Error(fmt.Sprintf("Error creating URL: %s", err))
		return res, nil, err
	}

	annotation, err = c.getResourcesFromAPI(ctx, queryUrl, &res)
	if err != nil {
		l.Error(fmt.Sprintf("Error getting resource: %s", err))
		return res, nil, err
	}

	return res, annotation, nil
}

func (c *FluidTopicsClient) GetRolesByUserID(ctx context.Context, userID string) (UserRoles, annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)
	var user UserRoles
	var annotation annotations.Annotations

	queryUrl, err := url.JoinPath(c.baseURL, fmt.Sprintf(getUserRolesById, userID))
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

func (c *FluidTopicsClient) UpdateUserManualRoles(ctx context.Context, userID string, manualRoles []string) (annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)

	queryUrl, err := url.JoinPath(c.baseURL, fmt.Sprintf(getUserRolesById, userID))
	if err != nil {
		l.Error("error creating URL", zap.Error(err))
		return nil, err
	}

	body := map[string]interface{}{
		"manualRoles": manualRoles,
	}

	_, annotation, err := c.doRequest(ctx, http.MethodPut, queryUrl, nil, body)
	if err != nil {
		return nil, err
	}

	return annotation, nil
}

func (c *FluidTopicsClient) CreateUser(ctx context.Context, newUser NewUserInfo) (annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)

	queryUrl, err := url.JoinPath(c.baseURL, createUser)
	if err != nil {
		l.Error("error creating URL", zap.Error(err))
		return nil, err
	}

	_, annotation, err := c.doRequest(ctx, http.MethodPost, queryUrl, nil, newUser)
	if err != nil {
		l.Error(fmt.Sprintf("Error getting resources: %s", err))
		return nil, err
	}

	return annotation, nil
}

func (c *FluidTopicsClient) getResourcesFromAPI(
	ctx context.Context,
	urlAddress string,
	out interface{},
) (annotations.Annotations, error) {
	_, annotation, err := c.doRequest(ctx, http.MethodGet, urlAddress, out, nil)
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
	body interface{},
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

	authToken, err := c.tokenSource.Token()
	if err != nil {
		return nil, nil, err
	}

	req, err := c.httpClient.NewRequest(ctx,
		method,
		urlAddress,
		uhttp.WithContentTypeJSONHeader(),
		uhttp.WithAcceptJSONHeader(),
		uhttp.WithJSONBody(body),
	)
	if err != nil {
		return nil, nil, err
	}

	authToken.SetAuthHeader(req)

	var errRes FluidTopicsAPIError
	var rateLimitDesc v2.RateLimitDescription

	doOptions := []uhttp.DoOption{
		uhttp.WithErrorResponse(&errRes),
		uhttp.WithRatelimitData(&rateLimitDesc),
	}
	switch method {
	case http.MethodGet, http.MethodPut, http.MethodPost:
		if res != nil {
			doOptions = append(doOptions, uhttp.WithResponse(res))
		}
		resp, err = c.httpClient.Do(req, doOptions...)
		if resp != nil {
			defer resp.Body.Close()
		}
		if err != nil {
			if errRes.MessageStr != "" || errRes.ErrorText != "" {
				return nil, nil, errRes
			}
			return nil, nil, err
		}
	case http.MethodDelete:
		resp, err = c.httpClient.Do(req)
		if resp != nil {
			defer resp.Body.Close()
		}
	}

	annotation := annotations.Annotations{}
	annotation.WithRateLimiting(&rateLimitDesc)

	if resp != nil {
		return resp.Header, annotation, nil
	}

	return nil, annotation, err
}
