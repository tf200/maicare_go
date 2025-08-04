package grpclient

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	// Add any necessary fields here, such as connection or client instances
	conn   *grpc.ClientConn
	client CarePlannerClient
}

func NewGrpcClient() (*GrpcClient, error) {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &GrpcClient{
		conn:   conn,
		client: NewCarePlannerClient(conn),
	}, nil
}
func (c *GrpcClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
