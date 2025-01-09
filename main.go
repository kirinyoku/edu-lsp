package main

import (
	"bufio"
	"log"
	"os"

	"github.com/kirinyoku/lsp/rpc"
)

func main() {
	logger := mustGetLogger("/home/kirin/personal/github/lsp/log.txt")
	logger.Println("Starting LSP server")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	for scanner.Scan() {
		msg := scanner.Text()
		handleMessage(logger, msg)
	}
}

func handleMessage(logger *log.Logger, msg any) {
	// TODO: Implement message handling
	logger.Println(msg)
}

func mustGetLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}

	return log.New(logfile, "[lsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
