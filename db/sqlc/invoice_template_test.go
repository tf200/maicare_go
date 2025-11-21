package db

import (
	"context"
	"testing"

	"maicare_go/util"

	"github.com/stretchr/testify/require"
)

func randomTemplateItem() TemplateItem {
	return TemplateItem{
		ItemTag:      util.RandomString(10),
		Description:  util.RandomString(20),
		SourceTable:  util.RandomString(10),
		SourceColumn: util.RandomString(10),
	}
}

func createRandomTemplateItem(ctx context.Context, t *testing.T, q *Queries) TemplateItem {
	item := randomTemplateItem()
	err := q.db.QueryRow(ctx, "INSERT INTO template_items (item_tag, description, source_table, source_column) VALUES ($1, $2, $3, $4) RETURNING id",
		item.ItemTag, item.Description, item.SourceTable, item.SourceColumn).Scan(&item.ID)
	require.NoError(t, err)
	return item
}

func TestGetAllTemplateItems(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) []TemplateItem
		checks func(t *testing.T, items []TemplateItem, inserted []TemplateItem, err error)
	}{
		{
			name: "successful get all",
			setup: func(ctx context.Context, qtx *Queries) []TemplateItem {
				item1 := createRandomTemplateItem(ctx, t, qtx)
				item2 := createRandomTemplateItem(ctx, t, qtx)
				return []TemplateItem{item1, item2}
			},
			checks: func(t *testing.T, items []TemplateItem, inserted []TemplateItem, err error) {
				require.NoError(t, err)
				require.Len(t, items, 2)
				// Check that all inserted items are in the result
				for _, ins := range inserted {
					found := false
					for _, item := range items {
						if item.ID == ins.ID {
							require.Equal(t, ins.ItemTag, item.ItemTag)
							require.Equal(t, ins.Description, item.Description)
							require.Equal(t, ins.SourceTable, item.SourceTable)
							require.Equal(t, ins.SourceColumn, item.SourceColumn)
							found = true
							break
						}
					}
					require.True(t, found, "Inserted item not found in results")
				}
			},
		},
		{
			name: "no items",
			setup: func(ctx context.Context, qtx *Queries) []TemplateItem {
				return []TemplateItem{}
			},
			checks: func(t *testing.T, items []TemplateItem, inserted []TemplateItem, err error) {
				require.NoError(t, err)
				require.Len(t, items, 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			inserted := tt.setup(ctx, qtx)
			items, err := qtx.GetAllTemplateItems(ctx)
			tt.checks(t, items, inserted, err)
		})
	}
}

func TestGetTemplateItemsByIds(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) ([]int64, []int64)
		checks func(t *testing.T, ids []int64, expected []int64, err error)
	}{
		{
			name: "successful get by ids",
			setup: func(ctx context.Context, qtx *Queries) ([]int64, []int64) {
				item1 := createRandomTemplateItem(ctx, t, qtx)
				item2 := createRandomTemplateItem(ctx, t, qtx)
				inputIds := []int64{item1.ID, item2.ID}
				expected := []int64{item1.ID, item2.ID}
				return inputIds, expected
			},
			checks: func(t *testing.T, ids []int64, expected []int64, err error) {
				require.NoError(t, err)
				require.Len(t, ids, len(expected))
				require.ElementsMatch(t, ids, expected)
			},
		},
		{
			name: "some ids not found",
			setup: func(ctx context.Context, qtx *Queries) ([]int64, []int64) {
				item1 := createRandomTemplateItem(ctx, t, qtx)
				inputIds := []int64{item1.ID, 99999} // 99999 doesn't exist
				expected := []int64{item1.ID}
				return inputIds, expected
			},
			checks: func(t *testing.T, ids []int64, expected []int64, err error) {
				require.NoError(t, err)
				require.Len(t, ids, len(expected))
				require.ElementsMatch(t, ids, expected)
			},
		},
		{
			name: "no ids found",
			setup: func(ctx context.Context, qtx *Queries) ([]int64, []int64) {
				inputIds := []int64{99999, 88888}
				expected := []int64{}
				return inputIds, expected
			},
			checks: func(t *testing.T, ids []int64, expected []int64, err error) {
				require.NoError(t, err)
				require.Len(t, ids, 0)
			},
		},
		{
			name: "empty ids",
			setup: func(ctx context.Context, qtx *Queries) ([]int64, []int64) {
				inputIds := []int64{}
				expected := []int64{}
				return inputIds, expected
			},
			checks: func(t *testing.T, ids []int64, expected []int64, err error) {
				require.NoError(t, err)
				require.Len(t, ids, 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			inputIds, expected := tt.setup(ctx, qtx)
			ids, err := qtx.GetTemplateItemsByIds(ctx, inputIds)
			tt.checks(t, ids, expected, err)
		})
	}
}

func TestGetTemplateItemsBySourceTable(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) ([]int64, []TemplateItem)
		checks func(t *testing.T, items []TemplateItem, expected []TemplateItem, err error)
	}{
		{
			name: "successful get by source table",
			setup: func(ctx context.Context, qtx *Queries) ([]int64, []TemplateItem) {
				item1 := createRandomTemplateItem(ctx, t, qtx)
				item2 := createRandomTemplateItem(ctx, t, qtx)
				item3 := createRandomTemplateItem(ctx, t, qtx)
				inputIds := []int64{item1.ID, item2.ID, item3.ID}
				expected := []TemplateItem{item1, item2, item3}
				return inputIds, expected
			},
			checks: func(t *testing.T, items []TemplateItem, expected []TemplateItem, err error) {
				require.NoError(t, err)
				require.Len(t, items, len(expected))
				// Since ordered by source_table, we need to check the order
				// For simplicity, just check all items are present
				for _, exp := range expected {
					found := false
					for _, item := range items {
						if item.ID == exp.ID {
							require.Equal(t, exp.ItemTag, item.ItemTag)
							require.Equal(t, exp.Description, item.Description)
							require.Equal(t, exp.SourceTable, item.SourceTable)
							require.Equal(t, exp.SourceColumn, item.SourceColumn)
							found = true
							break
						}
					}
					require.True(t, found, "Expected item not found in results")
				}
			},
		},
		{
			name: "some ids not found",
			setup: func(ctx context.Context, qtx *Queries) ([]int64, []TemplateItem) {
				item1 := createRandomTemplateItem(ctx, t, qtx)
				inputIds := []int64{item1.ID, 99999}
				expected := []TemplateItem{item1}
				return inputIds, expected
			},
			checks: func(t *testing.T, items []TemplateItem, expected []TemplateItem, err error) {
				require.NoError(t, err)
				require.Len(t, items, len(expected))
				require.Equal(t, expected[0].ID, items[0].ID)
				require.Equal(t, expected[0].ItemTag, items[0].ItemTag)
				require.Equal(t, expected[0].Description, items[0].Description)
				require.Equal(t, expected[0].SourceTable, items[0].SourceTable)
				require.Equal(t, expected[0].SourceColumn, items[0].SourceColumn)
			},
		},
		{
			name: "no ids found",
			setup: func(ctx context.Context, qtx *Queries) ([]int64, []TemplateItem) {
				inputIds := []int64{99999}
				expected := []TemplateItem{}
				return inputIds, expected
			},
			checks: func(t *testing.T, items []TemplateItem, expected []TemplateItem, err error) {
				require.NoError(t, err)
				require.Len(t, items, 0)
			},
		},
		{
			name: "empty ids",
			setup: func(ctx context.Context, qtx *Queries) ([]int64, []TemplateItem) {
				inputIds := []int64{}
				expected := []TemplateItem{}
				return inputIds, expected
			},
			checks: func(t *testing.T, items []TemplateItem, expected []TemplateItem, err error) {
				require.NoError(t, err)
				require.Len(t, items, 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "Failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			inputIds, expected := tt.setup(ctx, qtx)
			items, err := qtx.GetTemplateItemsBySourceTable(ctx, inputIds)
			tt.checks(t, items, expected, err)
		})
	}
}
