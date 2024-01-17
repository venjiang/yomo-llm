package llm

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	"github.com/invopop/jsonschema"
)

// ToolCall is the tool call in Request and Response
type ToolCall struct {
	ID       string    `json:"id,omitempty"` // present in Response only
	Type     string    `json:"type"`
	Function *Function `json:"function"`
}

// Function is the function in Request and Response
type Function struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Parameters  *FunctionParameters `json:"parameters,omitempty"` // chatCompletionFunctionParameters
	Arguments   string              `json:"arguments,omitempty"`
}

// FunctionParameters defines the parameters the functions accepts.
// from API doc: "The parameters the functions accepts, described as a JSON Schema object. See the [guide](/docs/guides/gpt/function-calling) for examples, and the [JSON Schema reference](https://json-schema.org/understanding-json-schema/) for documentation about the format."
// INFO: 对应 model
type FunctionParameters struct {
	Type       string                        `json:"type"`
	Properties map[string]*ParameterProperty `json:"properties"`
	Required   []string                      `json:"required"`
}

// ParameterProperty defines the property of the parameters
type ParameterProperty struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
}

// LLMFunctionCaller is the interface for calling function
type LLMFunctionCaller interface {
	AppID() string

	Name() string

	SetDescription(string)

	Description() string

	SetModel(any)

	Model() any
}

var (
	tools = make(map[string][]ToolCall)
	mu    sync.Mutex
)

func Register(fnc LLMFunctionCaller) error {
	// TODO: register function to LLM
	mu.Lock()
	defer mu.Unlock()
	// get old tool calls
	toolCalls := tools[fnc.AppID()]
	params, err := parseFunctionParameters(fnc.Model())
	if err != nil {
		slog.Error("parse function parameters", "error", err)
		return err
	}
	fn := &Function{
		Name:        fnc.Name(),
		Description: fnc.Description(),
		Parameters:  params,
	}
	toolCalls = append(toolCalls, ToolCall{
		Type:     "function",
		Function: fn,
	})
	paramsJson, _ := json.Marshal(params)
	slog.Info("register LLM function",
		"name", fnc.Name(),
		"description", fnc.Description(),
		"parameters", string(paramsJson),
	)
	// set new tool calls
	tools[fnc.AppID()] = toolCalls
	return nil
}

func parseFunctionParameters(v any) (*FunctionParameters, error) {
	r := new(jsonschema.Reflector)
	// if err := r.AddGoComments("go-generic", "./"); err != nil {
	// 	return nil, err
	// }
	schema := r.Reflect(v)
	defs := schema.Definitions
	functionParameters := &FunctionParameters{}
	properties := make(map[string]*ParameterProperty)
	// INFO: 只有一个 model
	if len(defs) < 1 {
		slog.Warn("not found model")
		return nil, nil
	}
	for k, m := range defs {
		// INFO: k 就是模型类型名称
		slog.Info(k,
			// "id", m.ID,
			"type", m.Type,
			// "title", m.Title,
			"description", m.Description,
			// "comments", m.Comments,
			"required", m.Required,
		)
		// type
		functionParameters.Type = m.Type
		// required
		functionParameters.Required = m.Required
		// properties
		for pair := m.Properties.Oldest(); pair != nil; pair = pair.Next() {
			slog.Info("property",
				"name", pair.Key,
				"type", pair.Value.Type,
				"title", pair.Value.Title,
				"description", pair.Value.Description,
				// "comments", pair.Value.Comments,
				// "required", pair.Value.Required,
			)
			properties[pair.Key] = &ParameterProperty{
				Type:        pair.Value.Type,
				Description: pair.Value.Description,
				// Enum:        []string{},
			}
		}
		// INFO: 只取一个 model
		break
	}
	functionParameters.Properties = properties

	return functionParameters, nil
}

func Invoke(appID string, msg string) {
	// tools
	toolCalls := tools[appID]
	toolCallsJson, _ := json.Marshal(toolCalls)
	// llm request
	llmRequestFormat := `
	{
		"messages": [
				{"role": "system", "content": "You are a very helpful assistant. Your job is to choose the best possible action to solve the user question or task. If you don't know the answer, stop the conversation by saying "no func call"."},
				{"role": "user", "content": "%s"}
		],
		"tools": %s
	}`
	llmRequest := fmt.Sprintf(llmRequestFormat, msg, toolCallsJson)
	// TODO: 调用 LLM Completions API 请求
	// TODO: 调用 yomo source 发送数据到 zipperud
	fmt.Printf("Invoke LLM function: %s \napp_id=%s\n", llmRequest, appID)
}

// TODO: 需要实例化一个 LLMBridge 实例
// 1. 实例用于存储 Tools
// 2. 配置 LLM 所需参数
// 3. yomo source 实例, 用于发送数据到 zipper 以转发给 sfn
// 4. 如有必要,提供一个 receiver sfn 实例, 用于接收 ai sfn 处理结果
