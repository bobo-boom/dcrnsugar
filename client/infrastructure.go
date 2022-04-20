package client

import (
	"errors"
	"fmt"
	"github.com/bobo-boom/dcrnsugar/config"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	// ErrInvalidAuth is an error to describe the condition where the client
	// is either unable to authenticate or the specified endpoint is
	// incorrect.
	ErrInvalidAuth = errors.New("authentication failure")

	// ErrInvalidEndpoint is an error to describe the condition where the
	// websocket handshake failed with the specified endpoint.
	ErrInvalidEndpoint = errors.New("the endpoint either does not support " +
		"websockets or does not exist")

	// ErrClientNotConnected is an error to describe the condition where a
	// websocket client has been created, but the connection was never
	// established.  This condition differs from ErrClientDisconnect, which
	// represents an established connection that was lost.
	ErrClientNotConnected = errors.New("the client was never connected")

	// ErrClientDisconnect is an error to describe the condition where the
	// client has been disconnected from the RPC server.  When the
	// DisableAutoReconnect option is not set, any outstanding futures
	// when a client disconnect occurs will return this error as will
	// any new requests.
	ErrClientDisconnect = errors.New("the client has been disconnected")

	// ErrClientShutdown is an error to describe the condition where the
	// client is either already shutdown, or in the process of shutting
	// down.  Any outstanding futures when a client shutdown occurs will
	// return this error as will any new requests.
	ErrClientShutdown = errors.New("the client has been shutdown")

	// ErrNotWebsocketClient is an error to describe the condition of
	// calling a Client method intended for a websocket client when the
	// client has been configured to run in HTTP POST mode instead.
	ErrNotWebsocketClient = errors.New("client is not configured for " +
		"websockets")

	// ErrClientAlreadyConnected is an error to describe the condition where
	// a new client connection cannot be established due to a websocket
	// client having already connected to the RPC server.
	ErrClientAlreadyConnected = errors.New("websocket client has already " +
		"connected")

	// ErrRequestCanceled is an error to describe the condition where
	// a request was canceled by the caller by terminating the passed
	// context.
	ErrRequestCanceled = errors.New("request was canceled by the caller")
)

const (
	// sendBufferSize is the number of elements the websocket send channel
	// can queue before blocking.
	sendBufferSize = 50

	// sendPostBufferSize is the number of elements the HTTP POST send
	// channel can queue before blocking.
	sendPostBufferSize = 100

	// connectionRetryInterval is the amount of time to wait in between
	// retries when automatically reconnecting to an RPC server.
	connectionRetryInterval = time.Second * 5

	// pingInterval is the amount of time between ping messages sent to
	// the server.
	pingInterval = time.Second * 10

	//url
	urlPrefix = "/insight/api/addr/"
)

type RequestDetails struct {
	HttpRequest *http.Request
}

// response is the raw bytes of a JSON-RPC result, or the error if the response
// error object was non-null.
type response struct {
	result []byte
	err    error
}

// jsonRequest holds information about a json request that is used to properly
// detect, interpret, and deliver a reply to it.
type jsonRequest struct {
	id             uint64
	method         string
	cmd            interface{}
	marshalledJSON []byte
	responseChan   chan *response
}

type Client struct {
	id uint64 // atomic, so must stay 64-bit aligned

	// config holds the connection configuration associated with this
	// client.
	config *config.Config

	// httpClient is the underlying HTTP client to use when running in HTTP
	// POST mode.
	httpClient *http.Client

	// mtx is a mutex to protect access to connection related fields.
	mtx sync.Mutex

	wg sync.WaitGroup
}

// String implements fmt.Stringer by returning the URL of the RPC server the
// client makes requests to.
func (c *Client) String() string {
	var u url.URL
	if c.config.EnableSSL {
		u.Scheme = "https"
	}
	u.Scheme = "http"

	u.Host = c.config.ServerHost
	u.Path = ""
	return u.String()
}

func newHttpClient(config *config.Config) (*http.Client, error) {

	client := http.Client{}

	return &client, nil
}
func (c *Client) SendRequest(details *RequestDetails) (respData []byte, err error) {

	log.Printf("Sending request........\n")
	httpResponse, err := c.httpClient.Do(details.HttpRequest)
	if err != nil {
		return nil, err
	}

	// Read the raw bytes and close the response.
	respBytes, err := io.ReadAll(httpResponse.Body)
	httpResponse.Body.Close()
	if err != nil {
		err = fmt.Errorf("error  reply: %v", err)
		return nil, err
	}
	return respBytes, err

}

func (c *Client) GenerateGetBalanceUrl(address string) string {

	return c.String() + urlPrefix + address
}

func New(config *config.Config) (*Client, error) {

	var httpClient *http.Client

	httpClient, err := newHttpClient(config)
	if err != nil {
		return nil, err
	}

	client := &Client{
		config:     config,
		httpClient: httpClient,
	}

	return client, nil
}
