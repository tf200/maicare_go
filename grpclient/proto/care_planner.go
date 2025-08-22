package grpclient

import (
	context "context"
	"time"
)

func (c *GrpcClient) GenerateCarePlan(ctx context.Context, req *PersonalizedCarePlanRequest) (*PersonalizedCarePlanResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()
	return c.carePlanClient.GenerateCarePlan(ctx, req)
}
