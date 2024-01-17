package main

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/venjiang/yomo-llm/llm"
	"github.com/venjiang/yomo-llm/yomo"

	"github.com/invopop/jsonschema"
)

var appID = "aid_test"

// UserModel test user model
type UserModel struct {
	Name string `json:"name" jsonschema:"title=nameA" jsonschema_description:"name schema description" `
	Age  int    `json:"age" jsonschema:"default=0,description=age schema description"`
	// Type model type
	Type string `json:"type,omitempty" jsonschema:"description=type schema description"`
}

func main() {
	// mdoel
	// model := serverless.Model{
	// 	Name: "tom",
	// 	Age:  18,
	// }
	// getJsonSchema(&serverless.Model{})
	// payload, err := json.Marshal(model)
	// if err != nil {
	// 	slog.Error("json marshal", "error", err)
	// 	return
	// }
	// userModel := UserModel{
	// 	Name: "jerry",
	// 	Age:  13,
	// 	Type: "user type",
	// }
	// getJsonSchema(&UserModel{})
	// userPayload, err := json.Marshal(userModel)
	// if err != nil {
	// 	slog.Error("json marshal", err)
	// 	return
	// }
	// normal context
	// ctx := serverless.NewNormalContext(userPayload)
	// // 在原始 handler 上再包一层, 构建一个泛型函数
	// handle[UserModel](ctx)

	// server context
	// 需要在服务端定义好 Model, 服务端需要知道 Model 的结构, 那么就没有意义了eb
	// ctx1 := serverless.NewServerContext(payload)
	// handleServerContext(ctx1)

	// wasm 版本实现上可能比较困难
	// ctx2 := serverless.NewGenericContext[UserModel](userPayload)
	// handleGenericContext(ctx2)

	// ---------- 平台代码 --------------
	sfn := yomo.NewStreamFunction("test", "127.0.0.1:9000")
	// set handler
	sfn.SetHandler(handler)
	// llm
	if llmfn, ok := any(sfn).(llm.LLMFunctionCaller); ok {
		// sfn.SetAppID(appID)
		llmfn.SetDescription("get weather")
		llmfn.SetModel(&UserModel{})
		if err := llm.Register(llmfn); err != nil {
			slog.Error("llm register", "error", err)
			return
		}
	}
	//
	sfn.Connect()
	sfn.Wait()
	llm.Invoke(appID, "hello llm")
}

// ---------------用户代码 app.go BEGIN ----------------
//
//	func Description() string{
//		return "get weather"
//	}
func DataTags() []uint32 {
	return []uint32{0x30}
}

func handler(ctx yomo.Context) {
	// raw bytes
	// data := ctx.Data()
	// slog.Info("handler ", "raw", string(data))
	// model
	var model UserModel
	ctx.ParseModel(&model)
	slog.Info("handler ", "model", model)
	// write
	// ctx.Write(0x33, []byte("hello world"))
}

// ---------------用户代码 app.go END ----------------

func getJsonSchema(v any) ([]byte, error) {
	r := new(jsonschema.Reflector)
	if err := r.AddGoComments("github.com/venjiang/yomo-llm", "./"); err != nil {
		return nil, err
	}
	schema := r.Reflect(v)
	defs := schema.Definitions
	slog.Info("", "schema", schema)
	for k, m := range defs {
		// m, ok := defs["k"]
		// if !ok {
		// 	slog.Error("not found model")
		// 	return nil, fmt.Errorf("not found model")
		// }
		slog.Info(k,
			// "id", m.ID,
			"type", m.Type,
			// "title", m.Title,
			"description", m.Description,
			// "comments", m.Comments,
			"required", m.Required,
			"properties", m.Properties,
		)
		for pair := m.Properties.Oldest(); pair != nil; pair = pair.Next() {
			slog.Info("->",
				"name", pair.Key,
				"type", pair.Value.Type,
				"title", pair.Value.Title,
				"description", pair.Value.Description,
				// "comments", pair.Value.Comments,
				// "required", pair.Value.Required,
			)
		}
	}
	schemaJson, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		slog.Error("json marshal", "error", err)
		return nil, err
	}
	fmt.Println(string(schemaJson))
	return schemaJson, nil
}

func handle[T any](ctx yomo.Context) {
	// data := ctx.Data()
	// 需要用户自己反序列化, 这样感觉意义也不大
	// T 只是标明了入参的类型
	var model T
	err := ctx.ParseModel(&model)
	if err != nil {
		slog.Error("model parse", err)
		return
	}
	// if err := json.Unmarshal(data, &model); err != nil {
	// 	slog.Error("json unmarshal", err)
	// 	return
	// }
	slog.Info("handler ", "model", model)
}

// handlerServerContext 需要在服务端定义好 Model, 服务端需要知道 Model 的结构, 那么就没有意义了
// func handleServerContext(ctx serverless.Context[serverless.Model]) {
// 	data := ctx.Data()
// 	slog.Info("handler ", "model", data)
// }

func handleGenericContext[T any](ctx yomo.GenericContext[T]) {
	data := ctx.Data()
	slog.Info("handler ", "model", data)
}
