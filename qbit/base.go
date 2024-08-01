package qbit

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type StatusError struct {
	Status int
}

func (e *StatusError) Error() string {
	return fmt.Sprintf("status code error: %d %s", e.Status, http.StatusText(e.Status))
}

type Client struct {
	client   *http.Client
	Host     string
	Username string
	Password string
	Cookie   string
}

func checkOk(buf []byte) error {
	if !bytes.Equal(buf, []byte("Ok.")) {
		return errors.New("login failed " + string(buf))
	}
	return nil
}
func NewClient(host, username, password string) *Client {
	return &Client{
		client:   &http.Client{},
		Host:     host,
		Username: username,
		Password: password,
	}
}
func (c *Client) post(path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", c.Host+path, body)
	if err != nil {
		return nil, err
	}
	if c.Cookie != "" {
		req.Header.Set("Cookie", c.Cookie)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		_ = resp.Body.Close()
		return nil, &StatusError{Status: resp.StatusCode}
	}
	return resp, err
}

func (c *Client) Login() error {
	path := "/api/v2/auth/login"
	resp, err := c.post(path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err = checkOk(buf); err != nil {
		return err
	}
	c.Cookie = resp.Header.Get("Set-Cookie")
	return nil
}
