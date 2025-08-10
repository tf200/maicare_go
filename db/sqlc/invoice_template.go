package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Table[T any] struct {
	Name    string
	Columns T
}

type ClientDetailsColumns struct {
	DateOfBirth string
	Filenumber  string
}
type ContractColumns struct {
	FinancingAct    string
	FinancingOption string
}

var (
	TableClientDetails = Table[ClientDetailsColumns]{
		Name: "client_details",
		Columns: ClientDetailsColumns{
			DateOfBirth: "date_of_birth",
			Filenumber:  "filenumber",
		},
	}
	TableContract = Table[ContractColumns]{
		Name: "contract",
		Columns: ContractColumns{
			FinancingAct:    "financing_act",
			FinancingOption: "financing_option",
		},
	}
)

type FetchQueryData struct {
	ClientID   int64
	ContractID int64
	SenderID   int64
}

func (store *Store) FetchInvoiceTemplateItems(ctx context.Context, data FetchQueryData) (map[string]string, error) {
	extraContent := make(map[string]string)
	templItemsIds, err := store.Queries.GetSenderInvoiceTemplate(ctx, data.SenderID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no invoice template found for sender ID %d", data.SenderID)
		}
		return nil, fmt.Errorf("failed to get sender invoice template: %w", err)
	}

	if len(templItemsIds) == 0 {
		return nil, nil
	}

	templItems, err := store.Queries.GetTemplateItemsBySourceTable(ctx, templItemsIds)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no template items found for sender ID %d", data.SenderID)
		}
		return nil, fmt.Errorf("failed to get template items: %w", err)
	}
	if len(templItems) == 0 {
		return nil, nil
	}
	// group items by source table
	groupedItems := make(map[string][]TemplateItem)
	for _, item := range templItems {
		groupedItems[item.SourceTable] = append(groupedItems[item.SourceTable], item)
	}

	for table, items := range groupedItems {
		switch table {
		case TableClientDetails.Name:
			clientDetails, err := store.Queries.GetClientDetails(ctx, data.ClientID)
			if err != nil {
				if err == sql.ErrNoRows {
					return nil, fmt.Errorf("no client details found for client ID %d", data.ClientID)
				}
				return nil, fmt.Errorf("failed to get client details: %w", err)
			}
			for _, item := range items {
				switch item.SourceColumn {
				case TableClientDetails.Columns.DateOfBirth:
					extraContent[item.Description] = clientDetails.DateOfBirth.Time.Format("02-01-2006")
				case TableClientDetails.Columns.Filenumber:
					extraContent[item.Description] = clientDetails.Filenumber
				}
			}
		case TableContract.Name:
			contract, err := store.Queries.GetClientContract(ctx, data.ContractID)
			if err != nil {
				return nil, fmt.Errorf("failed to get contract for client ID %d: %w", data.ContractID, err)
			}
			for _, item := range items {
				switch item.SourceColumn {
				case TableContract.Columns.FinancingAct:
					extraContent[item.Description] = contract.FinancingAct
				case TableContract.Columns.FinancingOption:
					extraContent[item.Description] = contract.FinancingOption
				}
			}
		}

	}

	// If there are no items for the given template items, return nil
	if len(extraContent) == 0 {
		return nil, nil
	}

	return extraContent, nil
}
