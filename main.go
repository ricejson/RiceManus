package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/ricejson/rice-manus/agent/base"
	"github.com/ricejson/rice-manus/agent/react"
	"github.com/ricejson/rice-manus/agent/tool"
	"github.com/ricejson/rice-manus/tools"
)

var APIKey = os.Getenv("V3_API_KEY")

func main() {
	registry := tools.NewToolRegistry()
	registry.RegistryTool(tools.NewWebSearchTool("web_search", "使用百度搜索引擎搜索"))
	registry.RegistryTool(tools.NewWebSearchTool("terminate", "Terminate the interaction when the request is met OR if the assistant cannot proceed further with the task.  \n            \"When you have finished all the tasks, call this tool to end the work.  "))
	tools := registry.GetTools()
	// 构建智能体
	model, err := openai.NewChatModel(context.Background(), &openai.ChatModelConfig{
		BaseURL: "https://api.vveai.com/v1",
		Model:   "o3-mini",
		APIKey:  APIKey,
	})
	if err != nil {
		log.Fatalf(err.Error())
	}
	agent := base.NewBaseAgent(model, "rice-manus", 10, react.NewReActAgent(tool.NewToolAgent(tools)))
	res, err := agent.Run("帮我从百度上获取编程导航的一些信息")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(res)
}
