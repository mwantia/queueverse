package tools

import "github.com/mwantia/queueverse/pkg/plugin/shared"

var TimeGetCurrent = shared.ToolDefinition{
	Name:     "time_get_current",
	Type:     shared.TypeString,
	Required: []string{"timezone"},
	Description: `Get the current time in the specified timezone.
	The timezone must be a IANA compatible timezone.
	The output is in the following format 'Mon Jan 2 15:04:05'.
	Only use the toll, if the conversation specifically requires the current time.`,
	Properties: map[string]shared.ToolProperty{
		"timezone": {
			Type:        shared.TypeString,
			Description: "The timezone to use. Must be a IANA compatible time zone.",
		},
	},
}
