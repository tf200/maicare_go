package grpclient

import (
	context "context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClientInterface interface {
	GenerateCarePlan(ctx context.Context, req *PersonalizedCarePlanRequest) (*PersonalizedCarePlanResponse, error)
	Close() error
}

type GrpcClient struct {
	// Add any necessary fields here, such as connection or client instances
	conn   *grpc.ClientConn
	client CarePlannerClient
}

func NewGrpcClient(grpcHost string) (GrpcClientInterface, error) {
	conn, err := grpc.NewClient(grpcHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
