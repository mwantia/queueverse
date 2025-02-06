package anthropic

import (
	"context"
	"encoding/json"
	"log"
	"testing"

	"github.com/liushuangls/go-anthropic/v2"
	"github.com/mwantia/queueverse/internal/config"
	"github.com/mwantia/queueverse/internal/tools"
	"github.com/mwantia/queueverse/pkg/plugin/base"
	"github.com/mwantia/queueverse/pkg/plugin/shared"
)

const (
	ConfigPath = "../../tests/config.hcl"
)

func TestAnthropicProvider(t *testing.T) {
	cfg, err := config.ParseConfig(ConfigPath)
	if err != nil {
		t.Fatalf("Failed to parse config: %v", err)
	}

	pcm := cfg.GetPluginConfigMap(PluginName)
	plugin := AnthropicProvider{
		Context: context.TODO(),
	}

	if err := plugin.SetConfig(&base.PluginConfig{ConfigMap: pcm}); err != nil {
		t.Fatalf("Failed to set plugin config: %v", err)
	}

	handler := tools.NewTest()
	request := shared.ChatRequest{
		Model: string(anthropic.ModelClaude3Dot5HaikuLatest),
		Message: shared.Message{
			Content: "Tell me the current time in germany.",
		},
	}

	response, err := plugin.Chat(request, handler)
	if err != nil {
		t.Fatalf("Failed to perform chat request: %v", err)
	}

	debug, _ := json.Marshal(response)
	log.Println(string(debug))
}
