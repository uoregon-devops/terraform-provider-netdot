package netdot

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/google/go-querystring/query"
)

type HostCreationData struct {
	Name    string `url:"name"`
	Subnet  string `url:"subnet"`
	Address string `url:"address,omitempty"`
}

type RR struct {
	ID   int    `xml:"id,attr"`
	Name string `xml:"name,attr"`
	Zone string `xml:"zone,attr"`
}

type Ipblock struct {
	ID          int    `xml:"id,attr"`
	Address     string `xml:"address,attr"`
	Description string `xml:"description,attr"`
	Parent      string `xml:"parent,attr"`
}

type HostQueryResponse struct {
	Ipblocks []Ipblock `xml:"Ipblock"`
	RRs      []RR      `xml:"RR"`
}

func (c *Client) CreateHost(name, subnet, address string) (RR, HostQueryResponse, error) {
	hostData := HostCreationData{
		Name:    name,
		Subnet:  subnet,
		Address: address,
	}
	query, err := query.Values(hostData)
	if err != nil {
		return RR{}, HostQueryResponse{}, err
	}

	createRequest, err := c.NewRequest("POST", "/rest/host?"+query.Encode(), nil)
	if err != nil {
		return RR{}, HostQueryResponse{}, err
	}
	createResponse, err := http.DefaultClient.Do(createRequest)
	if err != nil {
		return RR{}, HostQueryResponse{}, err
	}

	if createResponse.StatusCode != 200 {
		return RR{}, HostQueryResponse{}, fmt.Errorf("unexpected status code: %d", createResponse.StatusCode)
	}

	defer createResponse.Body.Close()
	var newHost RR
	bodyBytes, err := io.ReadAll(createResponse.Body)
	if err != nil {
		return RR{}, HostQueryResponse{}, err
	}
	err = xml.Unmarshal(bodyBytes, &newHost)
	if err != nil {
		return RR{}, HostQueryResponse{}, err
	}

	newReq, err := c.NewRequest("GET", fmt.Sprintf("/rest/host?rrid=%d", newHost.ID), nil)
	if err != nil {
		return RR{}, HostQueryResponse{}, err
	}

	newRes, err := http.DefaultClient.Do(newReq)
	if err != nil {
		return RR{}, HostQueryResponse{}, err
	}

	defer newRes.Body.Close()
	var newHostQueryResponse HostQueryResponse
	bodyBytes, err = io.ReadAll(newRes.Body)
	if err != nil {
		return RR{}, HostQueryResponse{}, err
	}
	err = xml.Unmarshal(bodyBytes, &newHostQueryResponse)
	if err != nil {
		return RR{}, HostQueryResponse{}, err
	}

	return newHost, newHostQueryResponse, nil
}

func (c *Client) DeleteHost(id int) error {
	req, err := c.NewRequest("DELETE", fmt.Sprintf("/rest/host?rrid=%d", id), nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) GetIpBlock(subnet string) (Ipblock, error) {
	req, err := c.NewRequest("GET", fmt.Sprintf("/rest/ipblock?address=%s", subnet), nil)
	if err != nil {
		return Ipblock{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Ipblock{}, err
	}

	defer resp.Body.Close()
	var IpBlocks HostQueryResponse
	bodyBytes, err := io.ReadAll(resp.Body)
	//print out the bodyBytes
	// fmt.Println(string(bodyBytes))
	if err != nil {
		return Ipblock{}, err
	}
	err = xml.Unmarshal(bodyBytes, &IpBlocks)
	if err != nil {
		return Ipblock{}, err
	}

	if len(IpBlocks.Ipblocks) == 0 {
		return Ipblock{}, fmt.Errorf("no IP block found")
	}

	return IpBlocks.Ipblocks[0], nil
}

func (c *Client) GetHost(ip string) (RR, error) {
	req, err := c.NewRequest("GET", fmt.Sprintf("/rest/host?address=%s", ip), nil)
	if err != nil {
		return RR{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return RR{}, err
	}

	defer resp.Body.Close()
	var hosts HostQueryResponse
	bodyBytes, err := io.ReadAll(resp.Body)
	//print out the bodyBytes
	// fmt.Println(string(bodyBytes))
	if err != nil {
		return RR{}, err
	}
	err = xml.Unmarshal(bodyBytes, &hosts)
	if err != nil {
		return RR{}, err
	}

	if len(hosts.RRs) == 0 {
		return RR{}, fmt.Errorf("no IP found")
	}

	return hosts.RRs[0], nil
}
