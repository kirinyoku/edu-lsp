package rpc_test

import (
	"testing"

	"github.com/kirinyoku/lsp/rpc"
)

type EncodingExample struct {
	Testing bool `json:"testing"`
}

func TestMustEncodeMessage(t *testing.T) {
	expected := "Content-Length: 16\r\n\r\n{\"testing\":true}"
	actual := rpc.MustEncodeMessage(EncodingExample{Testing: true})

	if expected != actual {
		t.Fatalf("expected %q, got %q", expected, actual)
	}
}

func TestDecodeMessage(t *testing.T) {
	incomingMessage := "Content-Length: 15\r\n\r\n{\"method\":\"hi\"}"
	method, content, err := rpc.DecodeMessage([]byte(incomingMessage))
	if err != nil {
		t.Fatal(err)
	}

	contentLength := len(content)

	if contentLength != 15 {
		t.Fatalf("expected content length 15, got %d", contentLength)
	}

	if method != "hi" {
		t.Fatalf("expected method \"hi\", got %q", method)
	}
}
