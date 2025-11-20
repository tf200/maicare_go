package db

import (
	"context"

	"maicare_go/util"
)
 



func TestCreateSender (t *testing.T) {









func createRandomSenders(ctx context.Context, qtx *Queries) Sender {
	arg := CreateSenderParams{
		Types:        "main_provider",
		Name:         util.RandomString(5),
		Address:      util.StringPtr("test"),
		PostalCode:   util.StringPtr("test"),
		Place:        util.StringPtr("test"),
		Land:         util.StringPtr("test"),
		Kvknumber:    util.StringPtr("test"),
		Btwnumber:    util.StringPtr("test"),
		PhoneNumber:  util.StringPtr("test"),
		ClientNumber: util.StringPtr("test"),
		EmailAddress: util.StringPtr("test"),
		Contacts:     []byte(`[{"name": "Test Contact", "email": "test@example.com", "phone": "1234567890"}]`),
	}

	sender, err := qtx.CreateSender(context.Background(), arg)
	if err != nil {
		panic(err)
	}
	_, err = qtx.CreateSenderInvoiceTemplate(context.Background(), CreateSenderInvoiceTemplateParams{
		ID:              sender.ID,
		InvoiceTemplate: []int64{1, 2, 3},
	})
	if err != nil {
		panic(err)
	}

	return sender
}
