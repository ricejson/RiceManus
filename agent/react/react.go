package react

import "github.com/ricejson/rice-manus/agent"

type ReAct interface {
	Think(*agent.AgentRuntime) (bool, error)
	Act(*agent.AgentRuntime) (string, error)
}

// ReActAgent Reason + Acting 智能体架构
type ReActAgent struct {
	reAct ReAct
}

func NewReActAgent(reAct ReAct) *ReActAgent {
	return &ReActAgent{reAct: reAct}
}

func (agent *ReActAgent) Step(runtime *agent.AgentRuntime) (string, error) {
	// 思考
	shouldAct, err := agent.reAct.Think(runtime)
	if err != nil {
		return "思考过程遇到错误", err
	}
	if !shouldAct {
		return "思考结束-无需行动", nil
	}
	return agent.reAct.Act(runtime)
}

func (agent *ReActAgent) Cleanup() error {
	return nil
}
