package agent

// ReActAgent Reason + Acting 智能体架构
type ReActAgent struct {
	BaseAgent
}

func (agent *ReActAgent) Step() (string, error) {
	// 思考
	shouldAct := agent.Think()
	if !shouldAct {
		return "思考结束-无需行动", nil
	}
	return agent.Act()
}

// Think 由具体子类实现
func (agent *ReActAgent) Think() bool {
	return false
}

// Act 由具体子类实现
func (agent *ReActAgent) Act() (string, error) {
	return "", nil
}
