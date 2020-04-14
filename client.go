package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	pb "gorpc/proto"
	"time"
)

type Client struct {
	client pb.DemoClient
}

func NewClient(conf GRPCConfig) (*Client, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", conf.Host, conf.Port), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("did not connect: %v", err)
	}

	client := Client{
		client: pb.NewDemoClient(conn),
	}

	return &client, nil
}

func (c *Client) ServerInfo() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.client.ServerInfo(ctx, &pb.EmptyRequest{})
	if err != nil {
		return "", fmt.Errorf("could not make request: %v", err)
	}

	return r.Info, nil
}
