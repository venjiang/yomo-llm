package yomo

import (
	"encoding/json"
	"log/slog"
)

type Context interface {
	Data() []byte
	Tag() uint32
	Write(tag uint32, data []byte) error
	// TODO: 这个是需要由用户提供反序列化方法
	ParseModel(model any) error
}

type GenericContext[T any] interface {
	Data() *T
	Tag() uint32
	Write(tag uint32, data []byte) error
}

// Model test model
type Model struct {
	Name string `json:"name,omitempty" jsonschema:"title=this is title,description=this is description"`
	// this is age comment
	Age int `json:"age" ai:"age"`
}

// normal context
type NormalContext struct {
	payload []byte
}

func NewNormalContext(payload []byte) *NormalContext {
	return &NormalContext{payload: payload}
}

func (c *NormalContext) Data() []byte {
	return c.payload
}

func (c *NormalContext) Tag() uint32 {
	return 0x29
}

func (c *NormalContext) Write(tag uint32, data []byte) error {
	slog.Info("NormalContext write", "tag", tag, "data", string(data))
	return nil
}

func (c *NormalContext) ParseModel(model any) error {
	err := json.Unmarshal(c.payload, model)
	return err
}

// // ServerContext server context
// type serverContext struct {
// 	payload []byte
// }
//
// func NewServerContext(payload []byte) Context {
// 	return &serverContext{payload: payload}
// }
//
// func (c *serverContext) Tag() uint32 {
// 	return 0x30
// }
//
// func (c *serverContext) Data() *Model {
// 	var data Model
// 	// json unmarshal payload to data
// 	if err := json.Unmarshal(c.payload, &data); err != nil {
// 		return nil
// 	}
// 	return &data
// }
//
// func (c *serverContext) Write(tag uint32, data []byte) error {
// 	slog.Info("ServerContext write", "tag", tag, "data", string(data))
// 	return nil
// }

// GenericContext generic context
type genericContext[T any] struct {
	payload []byte
}

// func NewGenericContext[T any](payload []byte) Context[T] {
func NewGenericContext[T any](payload []byte) GenericContext[T] {
	return &genericContext[T]{payload: payload}
}

func (c *genericContext[T]) Tag() uint32 {
	return 0x31
}

func (c *genericContext[T]) Data() *T {
	var data T
	// json unmarshal payload to data
	if err := json.Unmarshal(c.payload, &data); err != nil {
		return nil
	}
	return &data
}

func (c *genericContext[T]) Write(tag uint32, data []byte) error {
	slog.Info("GenericContext write", "tag", tag, "data", string(data))
	return nil
}
