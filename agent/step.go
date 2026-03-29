package agent

type Stepper interface {
	Step(runtime *AgentRuntime) (string, error)
	Cleanup() error
}
