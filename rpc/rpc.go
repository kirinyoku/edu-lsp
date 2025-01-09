package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

type BaseMessage struct {
	Method string `json:"method"`
}

// MustEncodeMessage encodes the given message as a JSON string and formats it
// into an LSP-style message with a "Content-Length" header.
// If JSON encoding fails, the function panics.
//
// Parameters:
// - msg: The message to encode. Can be any type that is JSON serializable.
//
// Returns:
// - A formatted string containing the "Content-Length" header and the message body.
func MustEncodeMessage(msg any) string {
	content, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
}

// DecodeMessage extracts and parses an LSP-style message from the given byte slice.
// It validates the presence of the "Content-Length" header, decodes the message content,
// and unmarshals it into a BaseMessage struct to retrieve the method.
//
// Parameters:
// - msg: A byte slice containing the LSP message to decode.
//
// Returns:
// - A string representing the method field from the BaseMessage.
// - A byte slice containing the message content limited to the content length.
// - An error if the header is malformed, the content length is invalid, or unmarshaling fails.
func DecodeMessage(msg []byte) (string, []byte, error) {
	const op = "DecodeMessage"

	header, content, found := bytes.Cut(msg, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return "", nil, fmt.Errorf("%s: did not find separator", op)
	}

	contentLengthBytes := header[len("Content-Length: "):]
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return "", nil, fmt.Errorf("%s: %w", op, err)
	}

	_ = content

	var baseMessage BaseMessage
	if err := json.Unmarshal(content[:contentLength], &baseMessage); err != nil {
		return "", nil, fmt.Errorf("%s: %w", op, err)
	}

	return baseMessage.Method, content[:contentLength], nil
}
