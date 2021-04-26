package etherclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type RPCRequest struct {
	Method  string          `json:"method"`
	Params  interface{}     `json:"params,omitempty"`
	ID      json.RawMessage `json:"id"`
	JSONRPC string          `json:"jsonrpc"`
}

type RPCRequests []*RPCRequest

type RPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	Result  interface{}     `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
	ID      json.RawMessage `json:"id"`
}

type RPCResponses []*RPCResponse

type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type HTTPError struct {
	Code int
	err  error
}

type rpcClient struct {
	endpoint      string
	httpClient    *http.Client
	customHeaders map[string]string
}

type RPCClientOpts struct {
	HTTPClient    *http.Client
	CustomHeaders map[string]string
}

type RPCClient interface {
	CallRaw(request *RPCRequest) (*RPCResponse, error)
	CallBatchRaw(requests RPCRequests) (RPCResponses, error)
}

func (e *HTTPError) Error() string {
	return e.err.Error()
}

func (e *RPCError) Error() string {
	return strconv.Itoa(e.Code) + ":" + e.Message
}

func NewHTTPClient(endpoint string) RPCClient {
	return NewClientWithOpts(endpoint, nil)
}

func NewClientWithOpts(endpoint string, opts *RPCClientOpts) RPCClient {
	rpcClient := &rpcClient{
		endpoint:      endpoint,
		httpClient:    &http.Client{},
		customHeaders: make(map[string]string),
	}

	if opts == nil {
		return rpcClient
	}

	if opts.HTTPClient != nil {
		rpcClient.httpClient = opts.HTTPClient
	}

	if opts.CustomHeaders != nil {
		for k, v := range opts.CustomHeaders {
			rpcClient.customHeaders[k] = v
		}
	}
	return rpcClient
}

func (client *rpcClient) newRequest(req interface{}) (*http.Request, error) {

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", client.endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	for k, v := range client.customHeaders {
		request.Header.Set(k, v)
	}

	return request, nil
}

func (client *rpcClient) doCall(RPCRequest *RPCRequest) (*RPCResponse, error) {

	httpRequest, err := client.newRequest(RPCRequest)
	if err != nil {
		return nil, fmt.Errorf("rpc call %v() on %v: %v", RPCRequest.Method, httpRequest.URL.String(), err.Error())
	}
	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("rpc call %v() on %v: %v", RPCRequest.Method, httpRequest.URL.String(), err.Error())
	}
	defer httpResponse.Body.Close()

	var rpcResponse *RPCResponse
	decoder := json.NewDecoder(httpResponse.Body)
	decoder.DisallowUnknownFields()
	decoder.UseNumber()
	err = decoder.Decode(&rpcResponse)

	if err != nil {
		if httpResponse.StatusCode >= 400 {
			return nil, &HTTPError{
				Code: httpResponse.StatusCode,
				err:  fmt.Errorf("rpc call %v() on %v status code: %v. could not decode body to rpc response: %v", RPCRequest.Method, httpRequest.URL.String(), httpResponse.StatusCode, err.Error()),
			}
		}
		return nil, fmt.Errorf("rpc call %v() on %v status code: %v. could not decode body to rpc response: %v", RPCRequest.Method, httpRequest.URL.String(), httpResponse.StatusCode, err.Error())
	}

	if rpcResponse == nil {
		if httpResponse.StatusCode >= 400 {
			return nil, &HTTPError{
				Code: httpResponse.StatusCode,
				err:  fmt.Errorf("rpc call %v() on %v status code: %v. rpc response missing", RPCRequest.Method, httpRequest.URL.String(), httpResponse.StatusCode),
			}
		}
		return nil, fmt.Errorf("rpc call %v() on %v status code: %v. rpc response missing", RPCRequest.Method, httpRequest.URL.String(), httpResponse.StatusCode)
	}

	return rpcResponse, nil
}

func (client *rpcClient) doBatchCall(rpcRequest []*RPCRequest) ([]*RPCResponse, error) {
	httpRequest, err := client.newRequest(rpcRequest)
	if err != nil {
		return nil, fmt.Errorf("rpc batch call on %v: %v", httpRequest.URL.String(), err.Error())
	}
	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("rpc batch call on %v: %v", httpRequest.URL.String(), err.Error())
	}
	defer httpResponse.Body.Close()

	var rpcResponse RPCResponses
	decoder := json.NewDecoder(httpResponse.Body)
	decoder.DisallowUnknownFields()
	decoder.UseNumber()
	err = decoder.Decode(&rpcResponse)

	if err != nil {
		if httpResponse.StatusCode >= 400 {
			return nil, &HTTPError{
				Code: httpResponse.StatusCode,
				err:  fmt.Errorf("rpc batch call on %v status code: %v. could not decode body to rpc response: %v", httpRequest.URL.String(), httpResponse.StatusCode, err.Error()),
			}
		}
		return nil, fmt.Errorf("rpc batch call on %v status code: %v. could not decode body to rpc response: %v", httpRequest.URL.String(), httpResponse.StatusCode, err.Error())
	}

	if rpcResponse == nil || len(rpcResponse) == 0 {
		if httpResponse.StatusCode >= 400 {
			return nil, &HTTPError{
				Code: httpResponse.StatusCode,
				err:  fmt.Errorf("rpc batch call on %v status code: %v. rpc response missing", httpRequest.URL.String(), httpResponse.StatusCode),
			}
		}
		return nil, fmt.Errorf("rpc batch call on %v status code: %v. rpc response missing", httpRequest.URL.String(), httpResponse.StatusCode)
	}

	return rpcResponse, nil
}

func (client *rpcClient) CallRaw(request *RPCRequest) (*RPCResponse, error) {

	return client.doCall(request)
}

func (client *rpcClient) CallBatchRaw(requests RPCRequests) (RPCResponses, error) {
	if len(requests) == 0 {
		return nil, errors.New("empty request list")
	}

	return client.doBatchCall(requests)
}
