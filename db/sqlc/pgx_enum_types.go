package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

var pgxEnumTypes = []string{
	"intake_participants_enum",
	"_intake_participants_enum",
	"informed_party_enum",
	"_informed_party_enum",
	"incident_cause_category_enum",
	"_incident_cause_category_enum",
	"incident_follow_up_action_enum",
	"_incident_follow_up_action_enum",
	"client_status_enum",
	"_client_status_enum",
}

func RegisterEnumTypes(ctx context.Context, conn *pgx.Conn) error {
	for _, typeName := range pgxEnumTypes {
		loadedType, err := conn.LoadType(ctx, typeName)
		if err != nil {
			return fmt.Errorf("load enum type %s: %w", typeName, err)
		}

		conn.TypeMap().RegisterType(loadedType)
	}

	return nil
}
