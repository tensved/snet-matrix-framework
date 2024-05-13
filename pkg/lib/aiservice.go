// Package lib provides functionality for managing AI services.
package lib

import "fmt"

// AIService represents an AI service.
type AIService struct {
	Name        string               // Name of the service
	Description string               // Description of the service
	Type        string               // Type of service: http, websocket, or grpc
	DataType    string               // Type of data: json, yaml, bytes
	Contact     string               // Developer's email
	BaseURL     string               // Endpoint for the service
	Methods     map[string]*AIMethod // List of methods that can be called on the service
}

// CreateMethod adds a new method to the AI service.
func (s *AIService) CreateMethod(method *AIMethod) {
	s.Methods[method.Name] = method
}

func (s *AIService) GetMethod(methodName string) *AIMethod {
	return s.Methods[methodName]
}

// CallMethod calls a method on the AI service.
func (s *AIService) CallMethod(methodName string, data map[string]interface{}) (interface{}, error) {
	method, ok := s.Methods[methodName]
	if !ok {
		return nil, fmt.Errorf("method %s not found", methodName)
	}

	// marshalData to input

	if err := method.CheckData(data); err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("invalid data")
	}

	c := NewMContext(WithParams(data))

	go method.Handler.Call(c)

	result := <-c.Result

	return result, nil
}

// IAIService defines the interface for an AI service.
type IAIService interface {
	CreateMethod(method AIMethod)
	GetMethod(methodName string) *AIMethod
	CallMethod(methodName string, data map[string]interface{}) error
}

// AIServiceOpts contains options for configuring an AI service.
type AIServiceOpts struct {
	Description string // Description of the service
	DataType    string // Type of data: json, yaml, bytes
	Contact     string // Developer's email
	BaseURL     string // Endpoint for the service
}

// AIServiceOptsFunc defines a function to apply options to an AI service.
type AIServiceOptsFunc func(service *AIService)

// WithOpts returns a function that applies the specified options to an AI service.
func WithOpts(options AIServiceOpts) AIServiceOptsFunc {
	return func(service *AIService) {
		service.Description = options.Description
		service.DataType = options.DataType
		service.Contact = options.Contact
		service.BaseURL = options.BaseURL
	}
}

// NewAIService creates a new AI service with the given name, service type and options.
func NewAIService(name, serviceType string, Opts ...AIServiceOptsFunc) *AIService {
	aiSvc := AIService{
		Name:    name,
		Type:    serviceType,
		Methods: make(map[string]*AIMethod),
	}

	for _, opt := range Opts {
		opt(&aiSvc)
	}

	return &aiSvc
}

// NewAIServiceHTTP creates a new HTTP AI service with the given name and options.
func NewAIServiceHTTP(name string, Opts ...AIServiceOptsFunc) *AIService {
	return NewAIService(name, "http", Opts...)
}

// NewAIServiceWS creates a new WebSocket AI service with the given name and options.
func NewAIServiceWS(name string, Opts ...AIServiceOptsFunc) *AIService {
	return NewAIService(name, "websocket", Opts...)
}
