package multichain

import (
	"fmt"
	"errors"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"encoding/base64"
	//
	"github.com/dghubble/sling"
)

const (
	CONST_ID = "multichain-client"
)

type Response map[string]interface{}

func (r Response) Result() interface{} {
	return r["result"]
}

type Client struct {
	chain string
	httpClient *http.Client
	port string
	endpoints []string
	credentials string
}

func NewClient(chain, host, port, username, password string) *Client {

	credentials := username + ":" + password

	return &Client{
		chain: chain,
		httpClient: &http.Client{},
		port: port,
		endpoints: []string{fmt.Sprintf("http://%s:%s", host, port)},
		credentials: base64.StdEncoding.EncodeToString([]byte(credentials)),
	}
}

// Creates a new temporary config for calling an RPC method on the specified node
func (client *Client) ViaNodes(hosts []int) *Client {

	c := *client

	for _, host := range hosts {

		c.endpoints = append(c.endpoints, fmt.Sprintf("http://%s:%s", host, client.port))

	}

	return &c
}

func (client *Client) debug() string {
	return strings.Join(client.endpoints, " ") + " " + client.credentials
}

func (client *Client) post(msg interface{}) (Response, error) {

	for i, endpoint := range client.endpoints {

		request, err := sling.New().Post(endpoint).BodyJSON(msg).Request()

		request.Header.Add("Authorization", "Basic " + client.credentials)

		resp, err := client.httpClient.Do(request)
		if err != nil {
			if (i + 1) == len(client.endpoints) {
				return nil, err
			}
			continue
		}

		if resp.StatusCode != 200 {
			if (i + 1) == len(client.endpoints) {
				return nil, err
			}
			continue
		}

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		obj := make(Response)

		err = json.Unmarshal(b, &obj)
		if err != nil {
			return nil, err
		}

		if obj["error"] != nil {
			e := obj["error"].(map[string]interface{})
			var s string
			m, ok := msg.(map[string]interface{})
			if ok {
				s = fmt.Sprintf("multichaind/%s: %s", m["method"], e["message"].(string))
			} else {
				s = fmt.Sprintf("multichaind: %s", e["message"].(string))
			}
			return nil, errors.New(s)
		}

		return obj, nil
	}

	return nil, errors.New("PROBABLY NO ENDPOINTS PASSED TO THE REQUEST DISPATCHER")
}
