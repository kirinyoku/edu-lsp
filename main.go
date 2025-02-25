package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/kirinyoku/lsp/analysis"
	"github.com/kirinyoku/lsp/lsp"
	"github.com/kirinyoku/lsp/rpc"
)

func main() {
	logger := mustGetLogger("/home/kirin/personal/github/lsp/log.txt")
	logger.Println("Starting LSP server")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	state := analysis.NewState()
	writer := os.Stdout

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Failed to decode message: %v\n", err)
			continue
		}

		handleMessage(logger, writer, state, method, contents)
	}
}

func handleMessage(logger *log.Logger, w io.Writer, state analysis.State, method string, contents []byte) {
	const op = "handleMessage"

	logger.Printf("Received message wtih method: %s\n", method)

	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("%s: Failed to unmarshal initialize request: %v\n", op, err)
			return
		}

		logger.Printf("Conected to: %s %s\n", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)

		response := lsp.NewInitializeResponse(request.ID)
		writeResponse(w, response)
	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("%s: Failed to unmarshal didOpen request: %v\n", op, err)
			return
		}

		logger.Printf("Opened: %s\n", request.Params.TextDocument.URI)

		diagnostics := state.OpenDocument(request.Params.TextDocument.URI, request.Params.TextDocument.Text)
		writeResponse(w, lsp.PublishDiagnosticsNotification{
			Notification: lsp.Notification{
				RPC:    "2.0",
				Method: "textDocument/publishDiagnostics",
			},
			Params: lsp.PublishDiagnosticsParams{
				URI:         request.Params.TextDocument.URI,
				Diagnostics: diagnostics,
			},
		})
	case "textDocument/didChange":
		var request lsp.DidChangeTextDocumenttNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("%s: Failed to unmarshal didChange request: %v\n", op, err)
			return
		}

		logger.Printf("Changed: %s\n", request.Params.TextDocument.URI)

		for _, change := range request.Params.ContentChanges {
			diagnostics := state.UpdateDocument(request.Params.TextDocument.URI, change.Text)
			writeResponse(w, lsp.PublishDiagnosticsNotification{
				Notification: lsp.Notification{
					RPC:    "2.0",
					Method: "textDocument/publishDiagnostics",
				},
				Params: lsp.PublishDiagnosticsParams{
					URI:         request.Params.TextDocument.URI,
					Diagnostics: diagnostics,
				},
			})
		}
	case "textDocument/hover":
		var request lsp.HoverRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("%s: Failed to unmarshal hover request: %v\n", op, err)
			return
		}

		response := state.Hover(request.ID, request.Params.TextDocument.URI, request.Params.Position)
		writeResponse(w, response)
	case "textDocument/definition":
		var request lsp.DefinitionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("%s: Failed to unmarshal definition request: %v\n", op, err)
			return
		}

		response := state.Definition(request.ID, request.Params.TextDocument.URI, request.Params.Position)
		writeResponse(w, response)
	case "textDocument/codeAction":
		var request lsp.CodeActionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("%s: Failed to unmarshal code action request: %v\n", op, err)
			return
		}

		response := state.TextDocumentCodeAction(request.ID, request.Params.TextDocument.URI)
		writeResponse(w, response)
	case "textDocument/completion":
		var request lsp.CompletionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("%s: Failed to unmarshal completion request: %v\n", op, err)
			return
		}

		response := state.TextDocumentCompletion(request.ID, request.Params.TextDocument.URI)
		writeResponse(w, response)
	}
}

func writeResponse(w io.Writer, resp any) {
	reply := rpc.MustEncodeMessage(resp)
	w.Write([]byte(reply))
}

func mustGetLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}

	return log.New(logfile, "[lsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
