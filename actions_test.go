package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActions(t *testing.T) {
	req := actions()
	require.NotEmpty(t, req.Actions)
	for _, a := range req.Actions {
		assert.True(t, json.Valid(a.ConfigUI), "%s config UI", a.ID)
		assert.True(t, json.Valid(a.ConfigSchema), "%s config schema", a.ID)
	}
}
