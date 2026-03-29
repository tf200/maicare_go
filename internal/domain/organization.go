package domain

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	ID                  uuid.UUID
	Name                string
	Street              string
	HouseNumber         string
	HouseNumberAddition *string
	PostalCode          string
	City                string
	Email               *string
	KvkNumber           *string
	BtwNumber           *string
	LocationCount       int64
}

type ListOrganizationsParams struct {
	Limit  int32
	Offset int32
	Search string
}

type OrganizationPage struct {
	Items      []Organization
	TotalCount int64
}

type OrganizationLocation struct {
	ID                  uuid.UUID
	OrganizationID      uuid.UUID
	Name                string
	Street              string
	HouseNumber         string
	HouseNumberAddition *string
	PostalCode          string
	City                string
	Capacity            *int32
	Occupied            int32
	Available           int32
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Shifts              []OrganizationLocationShift
}

type OrganizationLocationShift struct {
	ID         uuid.UUID
	LocationID uuid.UUID
	Slot       int16
	ShiftName  string
	StartTime  string
	EndTime    string
}

type ListOrganizationLocationsParams struct {
	OrganizationID uuid.UUID
	Limit          int32
	Offset         int32
	Search         string
}

type OrganizationLocationPage struct {
	Items      []OrganizationLocation
	TotalCount int64
}
