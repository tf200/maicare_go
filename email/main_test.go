package email

import (
	"log"
	"os"
	"testing"

	"maicare_go/util"
)

var testBrevo *BrevoConf

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatalf("Could not load conf %v", err)
	}
	testBrevo = NewBrevoConf(config.BrevoSenderName, config.BrevoSenderEmail, config.BrevoApiKey)
	os.Exit(m.Run())
}
