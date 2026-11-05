package main

import (
	_ "embed"

	"github.com/invopop/client.go/invopop"
	"github.com/invopop/popapp/internal/interfaces/gateway"
)

/*
 * ABOUT: Actions are the steps the app offers to workflows. They're declared
 * here and synced with Invopop's Access service each time the app starts, so
 * Invopop always has the app's current definitions. Each action's ID must
 * start with the app's ID, and the gateway handles the tasks for it.
 *
 * Steps are configured with a PopUI JSON UI panel the console renders, whose
 * "config" data is the step's configuration, and checked against a JSON
 * Schema that's kept apart from the panel.
 */

//go:embed actions/example.ui.json
var exampleConfigUI []byte

//go:embed actions/example.schema.json
var exampleConfigSchema []byte

// actions returns the definitions of every action the app provides.
func actions() *invopop.SyncActions {
	return &invopop.SyncActions{
		Actions: []*invopop.ActionDefinition{
			{
				ID:           gateway.ActionExample,
				Name:         "Example",
				Description:  "Logs a message, to show how actions work.",
				ConfigUI:     exampleConfigUI,
				ConfigSchema: exampleConfigSchema,
				Schemas:      []string{"bill/invoice"},
				Results: []*invopop.ActionResult{
					{Status: invopop.ActionResultOK, Code: "logged", Description: "Message logged", Default: true},
				},
			},
		},
	}
}
