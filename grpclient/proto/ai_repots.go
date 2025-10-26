package grpclient

import (
	context "context"
	"time"
)

func (c *GrpcClient) GenerateAutoReports(ctx context.Context, req *PastReports) (*GeneratedReports, error) {
	ctx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()
	return c.reportsClient.GenerateAutoReport(ctx, req)
}
