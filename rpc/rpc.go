package rpc

import (
	"encoding/json"
	"net/http"
)

type RPCRequest struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

type RPCResponse struct {
	Result interface{} `json:"result"`
	Error  string      `json:"error"`
}

type Rpchandler struct {
	AllRpc map[string]func(json.RawMessage) (interface{}, error)
}

func NewRpchandler() *Rpchandler {
	return &Rpchandler{
		AllRpc: make(map[string]func(json.RawMessage) (interface{}, error)),
	}
}

func (h *Rpchandler) AddMethod(method string, handler func(json.RawMessage) (interface{}, error)) {
	h.AllRpc[method] = handler
}

func (h *Rpchandler) HandleRPC(method string, params json.RawMessage) (interface{}, error) {
	if handler, ok := h.AllRpc[method]; ok {
		return handler(params)
	}
	return nil, http.ErrNotSupported
}
