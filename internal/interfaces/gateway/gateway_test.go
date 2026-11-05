package gateway

import (
	"context"
	"testing"

	"github.com/invopop/client.go/gateway"
	"github.com/stretchr/testify/assert"
)

func TestExecuteExample(t *testing.T) {
	s := new(Service)
	ctx := context.Background()

	res := s.executeAction(ctx, &gateway.Task{Action: ActionExample, Config: []byte(`{"message":"hi","uppercase":true}`)})
	assert.Equal(t, gateway.TaskStatus_OK, res.Status)
	assert.Equal(t, "logged", res.Code)

	res = s.executeAction(ctx, &gateway.Task{Action: ActionExample, Config: []byte(`not json`)})
	assert.Equal(t, gateway.TaskStatus_KO, res.Status)

	res = s.executeAction(ctx, &gateway.Task{Action: "popapp.unknown"})
	assert.Equal(t, gateway.TaskStatus_KO, res.Status)
}
