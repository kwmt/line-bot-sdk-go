// Copyright 2016 LINE Corporation
//
// LINE Corporation licenses this file to you under the Apache License,
// version 2.0 (the "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at:
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package linebot

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"path"
	"log"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

// APIEndpoint constants
const (
	APIEndpointBase = "https://api.line.me"

	APIEndpointPushMessage       = "/v2/bot/message/push"
	APIEndpointReplyMessage      = "/v2/bot/message/reply"
	APIEndpointMulticast         = "/v2/bot/message/multicast"
	APIEndpointGetMessageContent = "/v2/bot/message/%s/content"
	APIEndpointLeaveGroup        = "/v2/bot/group/%s/leave"
	APIEndpointLeaveRoom         = "/v2/bot/room/%s/leave"
	APIEndpointGetProfile        = "/v2/bot/profile/%s"
)

// Client type
type Client struct {
	ChannelSecret string
	ChannelToken  string
	EndpointBase  *url.URL     // default APIEndpointBase
	HttpClient    *http.Client // default http.DefaultClient
}

// ClientOption type
type ClientOption func(*Client) error

// New returns a new bot client instance.
func New(channelSecret, channelToken string, options ...ClientOption) (*Client, error) {
	if channelSecret == "" {
		return nil, errors.New("missing channel secret")
	}
	if channelToken == "" {
		return nil, errors.New("missing channel access token")
	}
	c := &Client{
		ChannelSecret: channelSecret,
		ChannelToken:  channelToken,
		HttpClient:    http.DefaultClient,
	}
	for _, option := range options {
		err := option(c)
		if err != nil {
			return nil, err
		}
	}
	if c.EndpointBase == nil {
		u, err := url.ParseRequestURI(APIEndpointBase)
		if err != nil {
			return nil, err
		}
		c.EndpointBase = u
	}
	return c, nil
}

// WithHTTPClient function
func WithHTTPClient(c *http.Client) ClientOption {
	return func(client *Client) error {
		client.HttpClient = c
		return nil
	}
}

// WithEndpointBase function
func WithEndpointBase(endpointBase string) ClientOption {
	return func(client *Client) error {
		u, err := url.ParseRequestURI(endpointBase)
		if err != nil {
			return err
		}
		client.EndpointBase = u
		return nil
	}
}

func (client *Client) url(endpoint string) string {
	u := *client.EndpointBase
	u.Path = path.Join(u.Path, endpoint)
	return u.String()
}

func (client *Client) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+client.ChannelToken)
	req.Header.Set("User-Agent", "LINE-BotSDK-Go/"+version)
	if ctx != nil {
		return ctxhttp.Do(ctx, client.HttpClient, req)
	}
	return client.HttpClient.Do(req)

}

func (client *Client) get(ctx context.Context, endpoint string) (*http.Response, error) {
	req, err := http.NewRequest("GET", client.url(endpoint), nil)
	if err != nil {
		return nil, err
	}
	return client.do(ctx, req)
}

func (client *Client) post(ctx context.Context, endpoint string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", client.url(endpoint), body)
	if err != nil {
		return nil, err
	}
	log.Println(req)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	return client.do(ctx, req)
}
