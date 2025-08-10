package grpclient

import (
	"log"
	"os"
	"testing"
)

var testClient *GrpcClient

func TestMain(m *testing.M) {
	client, err := NewGrpcClient()
	if err != nil {
		log.Fatalf("Could not create gRPC client: %v", err)
	}

	defer func() {
		err := client.Close()
		if err != nil {
			log.Fatalf("Could not close gRPC client: %v", err)
		}
	}()
	testClient = client.(*GrpcClient)

	os.Exit(m.Run())

}
