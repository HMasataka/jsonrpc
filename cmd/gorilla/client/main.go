package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/rpc/v2/json2"
)

type Reply struct {
	Message string
}

func main() {
	url := "http://localhost:8081/rpc"

	message := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "HelloService.SayHello",
		"params":  []interface{}{},
		"id":      1,
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to marshal request:", err)
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
