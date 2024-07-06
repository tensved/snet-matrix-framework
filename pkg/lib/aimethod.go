package lib

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
)

// Constraint defines an interface for value validation.
type Constraint interface {
	Check(value interface{}) error
}

// MInput represents metadata for an input parameter
type MInput struct {
	Name        string       // Name of the input
	Description string       // Description of the input
	Default     interface{}  // Default value of the input
	Required    bool         // Indicates if the input is required. Default: false
	MimeType    string       //  MIME type of the input
	Constraints []Constraint // Constraints applied to the input
}

// InputsBuilder constructs a slice of MInput from given parameters.
func InputsBuilder(inputs ...MInput) []MInput {
	return inputs
}

// OutputsBuilder creates a slice of MOutput from the provided outputs
func OutputsBuilder(outputs ...MOutput) []MOutput {
	return outputs
}

// MOutput represents output metadata for a method
type MOutput struct {
	Name        string       // Name of the output
	Description string       // Description of the output
	MimeType    string       // Mime type of the output
	Constraints []Constraint // Constraints applied to the output
}

type AIMethod struct {
	Name        string
	Description string

	Inputs  []MInput
	Outputs []MOutput

	Handler MHandler
}

func (m *AIMethod) CheckData(data map[string]any) error {
	for key, val := range data {
		for _, input := range m.Inputs {
			if input.Name == key {
				for _, constraint := range input.Constraints {
					if err := constraint.Check(val); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

// MethodOpts represents options for creating a new AIMethod instance
type MethodOpts struct {
	Description string
	Inputs      []MInput
	Outputs     []MOutput
}

// NewAIMethod creates a new AIMethod instance with specified options
func NewAIMethod(name string, opts MethodOpts, handler MHandler) *AIMethod {
	return &AIMethod{
		Name:        name,
		Description: opts.Description,
		Inputs:      opts.Inputs,
		Outputs:     opts.Outputs,
		Handler:     handler,
	}
}

// MHandler defines the interface for handling AI method calls
type MHandler interface {
	Call(c *MContext)
}

// HTTPHandler represents a default handler for making HTTP requests
type HTTPHandler struct {
	Endpoint string
	Method   string
}

// Call executes the HTTP request using the HTTPHandler
func (h *HTTPHandler) Call(c *MContext) {
	client := &http.Client{}

	data, err := json.Marshal(c.Params)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal params")
		return
	}

	req, err := http.NewRequest(h.Method, h.Endpoint, bytes.NewReader(data))
	if err != nil {
		log.Error().Err(err).Msg("failed to create request")
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("failed to make request")
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read response body")
		return
	}

	c.Result <- string(body)
}

// WSHandler represents a default handler for WebSocket connections
type WSHandler struct {
	Endpoint string
}

// Call executes the WebSocket connection using the WSHandler.
func (h *WSHandler) Call(c *MContext) {
	conn, _, err := websocket.DefaultDialer.Dial(h.Endpoint, nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to websocket")
		return
	}

	data, err := json.Marshal(c.Params)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal params")
		return
	}

	err = conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Error().Err(err).Msg("failed to write message to websocket")
		return
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Error().Err(err).Msg("failed to read message from websocket")
			return
		}

		c.Result <- string(message)
		return
	}
}
