package yomo

type AsyncHandler func(ctx Context)

// StreamFunction defines serverless streaming functions.
type StreamFunction interface {
	// SetObserveDataTags set the data tag list that will be observed
	// SetObserveDataTags(tag ...uint32)
	// Init will initialize the stream function
	// Init(fn func() error) error
	// SetHandler set the handler function, which accept the raw bytes data and return the tag & response
	SetHandler(fn AsyncHandler) error
	// SetErrorHandler set the error handler function when server error occurs
	// SetErrorHandler(fn func(err error))
	// SetPipeHandler set the pipe handler function
	// SetPipeHandler(fn core.PipeHandler) error
	// Connect create a connection to the zipper
	Connect() error
	// // Close will close the connection
	// Close() error
	// // Wait waits sfn to finish.
	Wait()
}

type sfn struct {
	name        string
	zipperAddr  string
	description string
	model       any
	fn          AsyncHandler
}

func NewStreamFunction(name string, zipperAddr string) StreamFunction {
	return &sfn{
		name:       name,
		zipperAddr: zipperAddr,
	}
}

// function
// func SetHandler(fn AsyncHandler[T]) error {

// method
func (s *sfn) SetHandler(fn AsyncHandler) error {
	s.fn = fn
	return nil
}

func (s *sfn) Connect() error {
	return nil
}

func (s *sfn) Wait() {
	// INFO: 仅用于测试
	if s.fn != nil {
		ctx := NewNormalContext([]byte(`{"name":"test","age":18,"type":"test"}`))
		s.fn(ctx)
	}
}

// ==================LLM=====================
func (s *sfn) AppID() string {
	// TODO: 来源于用户配置
	return "aid_test"
}

func (s *sfn) Name() string {
	return s.name
}

func (s *sfn) SetDescription(v string) {
	s.description = v
}

func (s *sfn) Description() string {
	return s.description
}

func (s *sfn) SetModel(v any) {
	s.model = v
}

func (s *sfn) Model() any {
	return s.model
}
