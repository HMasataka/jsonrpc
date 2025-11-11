package main

import (
	"context"
	"log"
	"net/http"

	"github.com/sourcegraph/jsonrpc2"
)

type myHandler struct{}

func (h *myHandler) Handle(ctx context.Context, c *jsonrpc2.Conn, r *jsonrpc2.Request) {
	switch r.Method {
	case "sayHello":
		if err := c.Reply(ctx, r.ID, "hello world"); err != nil {
			log.Println(err)
			return
		}
	default:
		err := &jsonrpc2.Error{Code: jsonrpc2.CodeMethodNotFound, Message: "Method not found"}
		if err := c.ReplyWithError(ctx, r.ID, err); err != nil {
			log.Println(err)
			return
		}
	}
}

func main() {
	http.HandleFunc("/rpc", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		conn := jsonrpc2.NewConn(ctx, jsonrpc2.NewBufferedStream(
			&httpReadWriteCloser{w: w, r: r},
			jsonrpc2.PlainObjectCodec{},
		), &myHandler{})
		<-conn.DisconnectNotify()
	})

	log.Println("Starting sourcegraph JSON-RPC server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

type httpReadWriteCloser struct {
	w http.ResponseWriter
	r *http.Request
}

func (h *httpReadWriteCloser) Read(p []byte) (n int, err error) {
	return h.r.Body.Read(p)
}

func (h *httpReadWriteCloser) Write(p []byte) (n int, err error) {
	return h.w.Write(p)
}

func (h *httpReadWriteCloser) Close() error {
	return h.r.Body.Close()
}
