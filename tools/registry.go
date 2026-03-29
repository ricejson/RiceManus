package tools

import "github.com/cloudwego/eino/components/tool"

type ToolRegistry struct {
	tools []tool.BaseTool
}

func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make([]tool.BaseTool, 0),
	}
}

func (registry *ToolRegistry) RegistryTool(tool tool.BaseTool) {
	registry.tools = append(registry.tools, tool)
}

func (registry *ToolRegistry) GetTools() []tool.BaseTool {
	return registry.tools
}
