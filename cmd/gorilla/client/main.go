package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/rpc/v2/json2"
)

type Args struct{}

type Reply struct {
	Message string
}

func main() {
	url := "http://localhost:8081/rpc"

	args := &Args{}
	jsonData, err := json2.EncodeClientRequest("HelloService.SayHello", args)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to encode request:", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to server:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read response:", err)
		return
	}

	var result Reply
	if err := json2.DecodeClientResponse(bytes.NewReader(body), &result); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to decode response:", err)
		return
	}

	fmt.Println(result.Message)
}
