package grpclient

import (
	context "context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClientInterface interface {
	GenerateCarePlan(ctx context.Context, req *PersonalizedCarePlanRequest) (*PersonalizedCarePlanResponse, error)
	CorrectSpelling(ctx context.Context, req *CorrectSpellingRequest) (*CorrectSpellingResponse, error)
	Close() error
}

type GrpcClient struct {
	// Add any necessary fields here, such as connection or client instances
	conn                *grpc.ClientConn
	carePlanClient      CarePlannerClient
	spellingCheckClient SpellingCorrectionClient
}

func NewGrpcClient(grpcHost string) (GrpcClientInterface, error) {
	conn, err := grpc.NewClient(grpcHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &GrpcClient{
		conn:                conn,
		carePlanClient:      NewCarePlannerClient(conn),
		spellingCheckClient: NewSpellingCorrectionClient(conn),
	}, nil
}
func (c *GrpcClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
