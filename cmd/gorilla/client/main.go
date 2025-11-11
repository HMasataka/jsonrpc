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

type MathArgs struct {
	A int `json:"a"`
	B int `json:"b"`
}

type MathReply struct {
	Result int `json:"result"`
}

func callJSONRPC(url, method string, args, result any) error {
	jsonData, err := json2.EncodeClientRequest(method, args)
	if err != nil {
		return fmt.Errorf("failed to encode request for %s: %w", method, err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to connect to server for %s: %w", method, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response for %s: %w", method, err)
	}

	if err := json2.DecodeClientResponse(bytes.NewReader(body), result); err != nil {
		return fmt.Errorf("failed to decode response for %s: %w", method, err)
	}

	return nil
}

func callSayHello(url string) error {
	args := &Args{}
	var result Reply

	if err := callJSONRPC(url, "HelloService.SayHello", args, &result); err != nil {
		return err
	}

	fmt.Printf("SayHello result: %s\n", result.Message)

	return nil
}

func callAddNumbers(url string, a, b int) error {
	args := &MathArgs{A: a, B: b}
	var result MathReply

	if err := callJSONRPC(url, "HelloService.AddNumbers", args, &result); err != nil {
		return err
	}

	fmt.Printf("AddNumbers result: %d + %d = %d\n", a, b, result.Result)

	return nil
}

func main() {
	url := "http://localhost:8081/"

	if err := callSayHello(url); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}

	if err := callAddNumbers(url, 15, 25); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
}
