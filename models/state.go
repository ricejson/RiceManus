package models

type AgentState int

const (
	AgentStateIDLE = iota
	AgentStateRunning
	AgentStateFinished
	AgentStateError
)
