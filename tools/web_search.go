package tools

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

// SearchAPIKey 从环境变量获取搜索 API Key
var SearchAPIKey = "SEARCH_API_KEY"
var searchURL = "https://www.searchapi.io/api/v1/search"

type WebSearchTool struct {
}

type WebSearchInput struct {
	Query string `json:"query" jsonschema:"required" jsonschema_description:"搜索关键词"`
}

func NewWebSearchTool(toolName string, toolDesc string) tool.InvokableTool {
	webSearchTool, err := utils.InferTool[WebSearchInput, string](
		toolName,
		toolDesc,
		func(ctx context.Context, input WebSearchInput) (string, error) {
			return webSearch(input.Query), nil
		},
	)
	if err != nil {
		panic(err)
	}
	return webSearchTool
}

func webSearch(query string) string {
	// 构造查询字符串
	params := url.Values{}
	params.Set("q", query)
	params.Set("engine", "baidu")
	params.Set("api_key", SearchAPIKey)
	// 构造请求
	req, err := http.NewRequest(http.MethodGet, searchURL, nil)
	if err != nil {
		return fmt.Sprintf("Error occurred: %v", err)
	}
	req.URL.RawQuery = params.Encode()
	// 发送请求
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Sprintf("Error occurred: %v", err)
	}
	defer resp.Body.Close()
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("Error occurred: %v", err)
	}
	return string(result)
}
