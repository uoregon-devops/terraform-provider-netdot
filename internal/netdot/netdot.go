package netdot

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
)

type netdotHttpError struct {
	statusCode  int
	errorString string
}

func (e netdotHttpError) Error() string {
	return e.errorString
}

type Client struct {
	server      string
	username    string
	password    string
	auth_cookie *http.Cookie
}

type authPayload struct {
	Username    string `xml:"credential_0"`
	Password    string `xml:"credential_1"`
	SessionType int    `xml:"permanent_session"`
}

func (c *Client) newAuthParameters() url.Values {
	return url.Values{
		"destination":       {"index.html"},
		"credential_0":      {c.username},
		"credential_1":      {c.password},
		"permanent_session": {"1"},
	}

}

func NewClient(server, username, password string) *Client {
	return &Client{
		server:   server,
		username: username,
		password: password,
	}
}

func (c *Client) getAuthCookie() (*http.Cookie, error) {
	authPayload := c.newAuthParameters()

	req, err := http.NewRequest("POST", c.server+"/NetdotLogin?"+authPayload.Encode(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "text/xml; version=1.0")
	req.Header.Set("User_Agent", "Netdot::Client::REST")

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	setCookie := resp.Header.Get("Set-Cookie")
	if setCookie == "" {
		return nil, fmt.Errorf("No cookie set")
	}

	cookieParts := strings.Split(setCookie, ";")
	cookieNameValue := strings.Split(cookieParts[0], "=")

	return &http.Cookie{
		Name:  cookieNameValue[0],
		Value: cookieNameValue[1],
		Raw:   setCookie,
	}, nil

}

func (c *Client) Authenticate() error {
	cookie, err := c.getAuthCookie()
	if err != nil {
		return err
	}
	c.auth_cookie = cookie
	return nil
}

func (c *Client) NewRequest(method, endpoint string, body io.Reader) (*http.Request, error) {
	req, _ := http.NewRequest(method, c.server+endpoint, body)
	req.AddCookie(c.auth_cookie)
	req.Header.Set("Accept", "text/xml; version=1.0")
	req.Header.Set("User-Agent", "gonsdb-client")
	return req, nil
}

func (c *Client) Get(endpoint string, v any) (*int, error) {
	req, err := c.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &resp.StatusCode, fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	}

	return nil, xml.NewDecoder(resp.Body).Decode(v)
}

// generic get resource type by id
func (c *Client) GetResourceByID(resourceType string, id int64, resource any) (*int, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid %s id, must be greater than 0", resourceType)
	}

	httpStatusCode, err := c.Get(fmt.Sprintf("/rest/%s/%d", resourceType, id), resource)
	if err != nil {
		return httpStatusCode, err
	}
	return httpStatusCode, nil
}

// generic delete resource type by id
func (c *Client) DeleteResourceByID(resourceType string, id int64, optionalQuery any) error {
	if id <= 0 {
		return fmt.Errorf("invalid %s id, must be greater than 0", resourceType)
	}

	endpoint := fmt.Sprintf("/rest/%s/%d", resourceType, id)

	if optionalQuery != nil {
		param_values, err := query.Values(optionalQuery)
		if err != nil {
			return err
		}
		endpoint = fmt.Sprintf("/rest/%s/%d?%s", resourceType, id, param_values.Encode())
	}

	req, err := c.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func boolToInt(b bool) int64 {
	if b {
		return 1
	}
	return 0
}

func (c *Client) CreateResource(resourceType string, inResource, outResource any) error {
	param_values, err := query.Values(inResource)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("/rest/%s?%s", resourceType, param_values.Encode())

	req, err := c.NewRequest("POST", endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// append body to error
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s, params: %s", resp.StatusCode, string(body), param_values.Encode())
	}

	if err := xml.NewDecoder(resp.Body).Decode(outResource); err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateResource(resourceType string, resourceID int64, inResource, outResource any) error {
	if resourceID <= 0 {
		return fmt.Errorf("invalid resource ID, must be greater than zero")
	}

	param_values, err := query.Values(inResource)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("/rest/%s/%d?%s", resourceType, resourceID, param_values.Encode())

	req, err := c.NewRequest("POST", endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := xml.NewDecoder(resp.Body).Decode(outResource); err != nil {
		return err
	}

	return nil
}
