package db

type Table struct {
	Name    string
	Columns [2]string
}

var (
	TablelientDetails = Table{
		Name:    "client_details",
		Columns: [2]string{"date_of_birth", "uuid_456"},
	}
)

// func (store *Store) FetchInvoiceTemplateItems(ctx context.Context, ids []int64) {
// 	var extraContent map[string]string
// 	templItems, err := store.Queries.GetInvoiceTemplateItems(ctx, ids)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, item := range templItems {
// 		if item.SourceTable == string(TablelientDetails.Name) {
// 			// get all columns from client tablle by client id will be passed in later
// 			clientDt, err := store.Queries.GetClientDetails(ctx, clientID)
// 			if err != nil {
// 				return nil, err
// 			}
// 			// keep only the all the columns  that matches the item source column
// 			val := reflect.ValueOf(clientDt)
// 			// get the filed by json tag
// 			found := false
// 			for i := 0; i < val.NumField(); i++ {
// 				field := val.Type().Field(i)
// 				fieldValue := val.Field(i)
// 				jsonTag := field.Tag.Get("json")
// 				jsonFieldName := strings.Split(jsonTag, ",")[0]
// 				if jsonFieldName != "" && jsonFieldName == item.SourceColumn {
// 					if fieldValue.Type() == reflect.TypeOf(pgtype.Date{}) {
// 						pgDate := fieldValue.Interface().(pgtype.Date)
// 						if pgDate.Valid {
// 							extraContent[item.Description] = pgDate.Time.Format("2006-01-02")
// 						} else {
// 							extraContent[item.Description] = ""
// 						}
// 					} else if fieldValue.CanInterface() {
// 						extraContent[item.Description] = fmt.Sprintf("%v", fieldValue.Interface())
// 					}
// 					found = true
// 					break

// 				}

// 			}
// 			if !found {
// 				extraContent[item.Description] = ""
// 			}

// 		}
// 	}

// }
