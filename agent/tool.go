package agent

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

type ToolAgent struct {
	ReActAgent
	availableTools []tool.BaseTool
	toolCallResp   *schema.Message
}

// Think 思考是否执行工具调用
func (agent *ToolAgent) Think() bool {
	ctx := context.Background()

	// 1. 如果有 nextStepPrompt，添加到消息列表
	if agent.nextStepPrompt != "" {
		agent.messages = append(agent.messages, schema.UserMessage(agent.nextStepPrompt))
		agent.nextStepPrompt = "" // 清空，避免重复添加
	}

	// 2. 准备工具信息
	toolInfos := make([]*schema.ToolInfo, 0, len(agent.availableTools))
	for _, t := range agent.availableTools {
		toolInfo, err := t.Info(ctx)
		if err != nil {
			log.Printf("获取工具信息失败: %v", err)
			continue
		}
		toolInfos = append(toolInfos, toolInfo)
	}

	// 3. 构建调用选项
	opts := []model.Option{
		model.WithTools(toolInfos),
	}

	// 4. 调用模型获取响应
	resp, err := agent.chatModel.Generate(ctx, agent.messages, opts...)
	if err != nil {
		log.Printf("%s 的思考过程遇到了问题: %v", "ToolAgent", err)
		agent.messages = append(agent.messages, schema.AssistantMessage(
			fmt.Sprintf("处理时遇到错误: %v", err),
			nil, // toolCalls
		))
		return false
	}

	// 5. 获取模型输出
	if resp == nil {
		return false
	}

	agent.toolCallResp = resp

	// 获取文本回复
	textContent := resp.Content
	log.Printf("%s 的思考: %s", "ToolAgent", textContent)

	// 获取工具调用列表
	toolCalls := resp.ToolCalls
	log.Printf("%s 选择了 %d 个工具来使用", "ToolAgent", len(toolCalls))

	// 打印工具调用详情
	for _, tc := range toolCalls {
		log.Printf("工具名称：%s，参数：%s", tc.Function.Name, tc.Function.Arguments)
	}

	// 6. 根据是否有工具调用决定是否执行 Act
	if len(toolCalls) == 0 {
		// 只有不调用工具时，才记录助手消息
		agent.messages = append(agent.messages, resp)
		return false
	} else {
		return true
	}
}

// Act 执行工具调用
func (agent *ToolAgent) Act() (string, error) {
	ctx := context.Background()

	if agent.toolCallResp == nil || len(agent.toolCallResp.ToolCalls) == 0 {
		return "", fmt.Errorf("没有工具需要调用")
	}

	// 执行工具调用
	var results []string
	for _, tc := range agent.toolCallResp.ToolCalls {
		toolName := tc.Function.Name
		arguments := tc.Function.Arguments

		// 查找对应的工具
		var targetTool tool.BaseTool
		for _, t := range agent.availableTools {
			info, err := t.Info(ctx)
			if err != nil {
				continue
			}
			if info.Name == toolName {
				targetTool = t
				break
			}
		}

		if targetTool == nil {
			errMsg := fmt.Sprintf("未找到工具: %s", toolName)
			log.Printf("%s", errMsg)
			results = append(results, errMsg)
			continue
		}

		// 检查工具是否可执行
		invokableTool, ok := targetTool.(tool.InvokableTool)
		if !ok {
			errMsg := fmt.Sprintf("工具 %s 不可执行", toolName)
			log.Printf("%s", errMsg)
			results = append(results, errMsg)
			continue
		}

		// 调用工具
		result, err := invokableTool.InvokableRun(ctx, arguments)
		if err != nil {
			errMsg := fmt.Sprintf("工具 %s 调用失败: %v", toolName, err)
			log.Printf("%s", errMsg)
			results = append(results, errMsg)
			continue
		}

		// 记录工具结果
		toolResult := &schema.Message{
			Role:       schema.Tool,
			Content:    result,
			ToolCallID: tc.ID,
			ToolName:   toolName,
		}
		agent.messages = append(agent.messages, toolResult)
		results = append(results, result)
	}

	return fmt.Sprintf("工具调用完成: %s", results), nil
}
