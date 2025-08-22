package grpclient

import (
	context "context"
	"time"
)

func (c *GrpcClient) CorrectSpelling(ctx context.Context, req *CorrectSpellingRequest) (*CorrectSpellingResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()
	return c.spellingCheckClient.CorrectSpelling(ctx, req)
}
