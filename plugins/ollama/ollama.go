package ollama

import (
	"log"

	"github.com/mwantia/queueverse/pkg/plugin/provider"
	"github.com/mwantia/queueverse/plugins/ollama/api"
)

func (p *OllamaProvider) GetModels() (*[]provider.Model, error) {
	tags, err := p.Client.Tags(p.Context)
	if err != nil {
		return nil, err
	}

	resp := make([]provider.Model, 0)
	for _, tag := range tags {
		resp = append(resp, provider.Model{
			Name: tag.Name,
			Metadata: map[string]any{
				"size":   tag.Size,
				"digest": tag.Digest,
			},
		})
	}
	return &resp, nil
}

func (p *OllamaProvider) Chat(input provider.ChatRequest) (*provider.ChatResponse, error) {
	tools := make([]api.ToolDefinition, 0)
	for _, tool := range input.Tools {

		properties := make(map[string]api.ToolProperty, 0)
		for name, property := range tool.Properties {
			properties[name] = api.ToolProperty{
				Type:        string(property.Type),
				Description: property.Description,
				Enum:        property.Enum,
			}
		}

		tools = append(tools, api.ToolDefinition{
			Type: "function",
			Function: api.ToolFunction{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters: api.ToolParameters{
					Type:       string(tool.Type),
					Required:   tool.Required,
					Properties: properties,
				},
			},
		})
	}

	request := api.ChatRequest{
		Model: input.Model,
		Messages: []api.ChatMessage{
			{
				Role:    "user",
				Content: input.Message.Content,
			},
		},
		Tools:       tools,
		Stream:      false,
		KeepAlive:   -1,
		ContextSize: 4096,
	}

	var output provider.ChatResponse

	if err := p.Client.Chat(p.Context, request, func(response api.ChatResponse) error {
		log.Printf("%v", response)
		return nil
	}); err != nil {
		return nil, err
	}

	return &output, nil
}

func (p *OllamaProvider) Embed(input provider.EmbedRequest) (*provider.EmbedResponse, error) {
	request := api.EmbedRequest{
		Model: input.Model,
		Input: input.Message.Content,
	}

	var output provider.EmbedResponse

	if err := p.Client.Embed(p.Context, request, func(response api.EmbedResponse) error {
		output = provider.EmbedResponse{
			Model:      response.Model,
			Embeddings: response.Embeddings,
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &output, nil
}
