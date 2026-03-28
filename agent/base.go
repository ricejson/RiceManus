package agent

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/ricejson/rice-manus/models"
)

// BaseAgent 仅仅定义循环步骤以及 agent 状态
type BaseAgent struct {
	chatModel model.BaseChatModel

	systemPrompt   string
	nextStepPrompt string

	currentStep int
	maxStep     int
	messages    []*schema.Message // 上下文列表
	results     []string
	agentState  models.AgentState
}

func NewBaseAgent(agent *adk.ChatModelAgent, maxStep int) *BaseAgent {
	return &BaseAgent{
		currentStep: 0,
		maxStep:     maxStep,
		messages:    make([]*schema.Message, 0),
		results:     make([]string, 0),
		agentState:  models.AgentStateIDLE,
	}
}

func (agent *BaseAgent) Run(userPrompt string) (string, error) {
	defer agent.Cleanup()
	if agent.agentState != models.AgentStateIDLE {
		return "", errors.New(fmt.Sprintf("Cannot run agent from state:%v", agent.agentState))
	}
	if userPrompt == "" {
		return "", errors.New("user prompt cannot be empty")
	}
	// 修改 agent 状态
	agent.agentState = models.AgentStateRunning
	// 记录用户消息
	agent.messages = append(agent.messages, schema.UserMessage(userPrompt))
	for i := 0; i < agent.maxStep && agent.agentState != models.AgentStateFinished; i++ {
		agent.currentStep = i + 1
		fmt.Printf("Executing step %d / %d\n", agent.currentStep, agent.maxStep)
		res, err := agent.Step()
		if err != nil {
			// TODO: 记录日志
			agent.agentState = models.AgentStateError
			agent.currentStep = 0
			return "", err
		}
		fmt.Printf("Executing step %d, result:%s\n", agent.currentStep, res)
		// 记录结果
		agent.results = append(agent.results, res)
	}
	// 超出最大步骤
	if agent.currentStep >= agent.maxStep {
		agent.agentState = models.AgentStateFinished
		agent.results = append(agent.results, fmt.Sprintf("Agent finished with %d steps", agent.currentStep))
	}
	return strings.Join(agent.results, "\n"), nil
}

func (agent *BaseAgent) Step() (string, error) {
	return "", nil
}

func (agent *BaseAgent) Cleanup() error {
	return nil
}
