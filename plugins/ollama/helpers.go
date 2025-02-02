package ollama

import (
	"github.com/mwantia/queueverse/pkg/plugin/provider"
	"github.com/mwantia/queueverse/plugins/ollama/api"
)

func CreateMessageRequest(req provider.ChatRequest) api.ChatRequest {
	messages := []api.ChatMessage{}
	for _, msg := range req.Messages {
		tools := []api.ToolCall{}
		for _, tool := range msg.ToolCalls {
			tools = append(tools, api.ToolCall{
				Function: api.ToolCallFunction{
					Index:     tool.Function.Index,
					Name:      tool.Function.Name,
					Arguments: tool.Function.Arguments,
				},
			})
		}

		messages = append(messages, api.ChatMessage{
			Role:      msg.Role,
			Content:   msg.Content,
			ToolCalls: tools,
		})
	}

	tools := []api.ToolDefinition{}
	for _, tool := range req.Tools {
		definition := api.ToolDefinition{
			Type: string(provider.ToolDefinitionFunction),
			Function: api.ToolFunction{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters: api.ToolParameters{
					Type:       string(tool.Parameters.Type),
					Required:   tool.Parameters.Required,
					Properties: make(map[string]api.ToolProperty),
				},
			},
		}

		for name, property := range tool.Parameters.Properties {
			definition.Function.Parameters.Properties[name] = api.ToolProperty{
				Type:        string(property.Type),
				Description: property.Description,
				Enum:        property.Enum,
			}
		}

		tools = append(tools, definition)
	}

	return api.ChatRequest{
		Model:       req.Model,
		Messages:    messages,
		Tools:       tools,
		Stream:      false,
		KeepAlive:   -1,
		ContextSize: 4096,
	}
}
