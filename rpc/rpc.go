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

// Split is a bufio.SplitFunc implementation for processing LSP-style messages.
// It parses the "Content-Length" header, identifies complete messages, and returns
// them as tokens. This function ensures the content length matches the header value
// and handles partial data appropriately.
//
// Parameters:
// - data: A byte slice containing the input data to split.
// - atEOF: A boolean indicating whether the end of the input has been reached.
//
// Returns:
// - advance: The number of bytes to advance in the input data for the next split.
// - token: A byte slice containing the complete message, or nil if the message is incomplete.
// - err: An error if the "Content-Length" header is malformed or parsing fails.
func Split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	const op = "Split"

	header, content, found := bytes.Cut(data, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return 0, nil, nil
	}

	contentLengthBytes := header[len("Content-Length: "):]
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return 0, nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(content) < contentLength {
		return 0, nil, nil
	}

	totalLength := len(header) + 4 + contentLength
	return totalLength, data[:totalLength], nil
}
