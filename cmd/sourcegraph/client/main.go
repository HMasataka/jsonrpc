package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	reqBody := bytes.NewBuffer([]byte(`{"jsonrpc":"2.0","method":"sayHello","params":null,"id":1}`))
	resp, err := http.Post("http://localhost:8080/rpc", "application/json", reqBody)
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

	fmt.Println(string(body))
}
