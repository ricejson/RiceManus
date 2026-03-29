package agent

import (
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/ricejson/rice-manus/models"
)

type AgentRuntime struct {
	Name       string
	Messages   []*schema.Message // 上下文列表
	Results    []string
	AgentState models.AgentState

	NextStepPrompt string
	SystemPrompt   string

	ChatModel model.BaseChatModel
}
