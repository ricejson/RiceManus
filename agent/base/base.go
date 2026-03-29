package base

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/ricejson/rice-manus/agent"
	"github.com/ricejson/rice-manus/models"
)

// BaseAgent 仅仅定义循环步骤以及 agent 状态
type BaseAgent struct {
	currentStep int
	maxStep     int

	runtime *agent.AgentRuntime

	stepper agent.Stepper
}

func NewBaseAgent(chatModel model.BaseChatModel, name string, maxStep int, stepper agent.Stepper) *BaseAgent {
	return &BaseAgent{
		currentStep: 0,
		maxStep:     maxStep,
		stepper:     stepper,
		runtime: &agent.AgentRuntime{
			Name:       name,
			Messages:   make([]*schema.Message, 0),
			Results:    make([]string, 0),
			AgentState: models.AgentStateIDLE,
			ChatModel:  chatModel,
		},
	}
}

// Run 模板方法：接收 Stepper 接口，调用子类的 Step()
func (agent *BaseAgent) Run(userPrompt string) (string, error) {
	defer agent.stepper.Cleanup()
	if agent.runtime.AgentState != models.AgentStateIDLE {
		return "", errors.New(fmt.Sprintf("Cannot run agent from state:%v", agent.runtime.AgentState))
	}
	if userPrompt == "" {
		return "", errors.New("user prompt cannot be empty")
	}
	// 修改 agent 状态
	agent.runtime.AgentState = models.AgentStateRunning
	// 记录用户消息
	agent.runtime.Messages = append(agent.runtime.Messages, schema.UserMessage(userPrompt))
	for i := 0; i < agent.maxStep && agent.runtime.AgentState != models.AgentStateFinished; i++ {
		agent.currentStep = i + 1
		fmt.Printf("Executing step %d / %d\n", agent.currentStep, agent.maxStep)
		res, err := agent.stepper.Step(agent.runtime)
		if err != nil {
			// TODO: 记录日志
			agent.runtime.AgentState = models.AgentStateError
			agent.currentStep = 0
			return "", err
		}
		fmt.Printf("Executing step %d, result:%s\n", agent.currentStep, res)
		// 记录结果
		agent.runtime.Results = append(agent.runtime.Results, res)
	}
	// 超出最大步骤
	if agent.currentStep >= agent.maxStep {
		agent.runtime.AgentState = models.AgentStateFinished
		agent.runtime.Results = append(agent.runtime.Results, fmt.Sprintf("Agent finished with %d steps", agent.currentStep))
	}
	return strings.Join(agent.runtime.Results, "\n"), nil
}
