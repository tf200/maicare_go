package bucket

import (
	"context"
	"log"
	"maicare_go/util"
	"os"
	"testing"
)

var testBucketClient ObjectStorageInterface

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatalf("Could not load conf %v", err)
	}
	testBucketClient, err = NewObjectStorageClient(context.Background(), config)
	if err != nil {
		log.Fatalf("Could not create bucket client %v", err)
	}

	os.Exit(m.Run())
}
