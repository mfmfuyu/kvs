package client

import (
	"example.com/kvs/resp"
)

type Client struct {
	writer *resp.Writer
}

func New(writer *resp.Writer) *Client {
	return &Client{
		writer: writer,
	}
}

func (c *Client) Write(value resp.Value) {
	c.writer.Write(value)
}
