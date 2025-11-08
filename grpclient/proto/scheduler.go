package grpclient

import (
	context "context"
	"time"
)

func (c *GrpcClient) AutoGenerateSchedules(ctx context.Context, req *GenerateScheduleRequest) (*GenerateScheduleResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()
	return c.scheduleClient.GenerateSchedule(ctx, req)
}
