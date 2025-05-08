package connector

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/conductorone/baton-fluid-topics/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/stretchr/testify/assert"
)

var (
	ctx              = context.Background()
	bearerToken      = os.Getenv("FLUID_TOPICS_BEARER_TOKEN")
	domain           = os.Getenv("FLUID_TOPICS_DOMAIN")
	parentResourceID = &v2.ResourceId{}
	pToken           = &pagination.Token{}
)

func initClient(t *testing.T) *client.FluidTopicsClient {
	if bearerToken == "" {
		message :=
			fmt.Sprintf("Any of the required params not found. Bearer token: %s", bearerToken)
		t.Skip(message)
	}

	c, err := client.New(
		ctx,
		bearerToken,
		domain,
	)

	if err != nil {
		t.Errorf("ERROR: Failed to create client: %v", err)
	}
	return c
}

func TestUserBuilderList(t *testing.T) {
	c := initClient(t)

	u := newUserBuilder(c)
	res, _, _, err := u.List(ctx, parentResourceID, pToken)
	assert.Nil(t, err)
	assert.NotNil(t, res)

	message := fmt.Sprintf("Amount of users obtained: %d", len(res))
	t.Log(message)
}

func TestRoleBuilderList(t *testing.T) {
	c := initClient(t)

	r := newRoleBuilder(c)

	res, _, _, err := r.List(ctx, parentResourceID, pToken)
	assert.Nil(t, err)
	assert.NotNil(t, res)

	message := fmt.Sprintf("Amount of roles obtained: %d", len(res))
	t.Log(message)
}
