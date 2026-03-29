package tools

import (
	"context"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

type TerminateTool struct {
}

func NewTerminateTool(toolName string, toolDesc string) tool.InvokableTool {
	terminateTool, err := utils.InferTool[any, string](
		toolName,
		toolDesc,
		func(ctx context.Context, args any) (string, error) {
			return doTerminate(), nil
		},
	)
	if err != nil {
		panic(err)
	}
	return terminateTool
}

func doTerminate() string {
	return "任务结束"
}
