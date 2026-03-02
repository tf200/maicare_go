package main

import (
	"context"
	"fmt"
	"strings"

	db "maicare_go/db/sqlc"

	"github.com/brianvoe/gofakeit/v7"
)

func (s *Seeder) SeedOrganisations(ctx context.Context, count int) error {
	for i := range count {
		if i == 0 || i+1 == count {
			fmt.Printf("[seed] organisations: %d/%d\n", i+1, count)
		}
		name := fmt.Sprintf("%s Care %d", gofakeit.Company(), i+1)
		phone := fakePhone()
		email := strings.ToLower(fmt.Sprintf("info.%d@%s.nl", i+1, sanitizeDomain(name)))
		kvk := fmt.Sprintf("%08d", gofakeit.Number(10000000, 99999999))
		btw := fmt.Sprintf("NL%09dB01", gofakeit.Number(100000000, 999999999))

		created, err := s.store.CreateOrganisation(ctx, db.CreateOrganisationParams{
			Name:                name,
			Street:              gofakeit.StreetName(),
			HouseNumber:         fmt.Sprintf("%d", gofakeit.Number(1, 350)),
			HouseNumberAddition: nullableString(strings.ToUpper(gofakeit.LetterN(1)), 0.8),
			PostalCode:          fakePostalCodeNL(),
			City:                gofakeit.City(),
			PhoneNumber:         &phone,
			Email:               &email,
			KvkNumber:           &kvk,
			BtwNumber:           &btw,
		})
		if err != nil {
			return fmt.Errorf("create organisation %d: %w", i+1, err)
		}

		s.data.OrganisationIDs = append(s.data.OrganisationIDs, created.ID)
	}

	return nil
}

func (s *Seeder) SeedLocations(ctx context.Context, perOrganisation int) error {
	if len(s.data.OrganisationIDs) == 0 {
		return fmt.Errorf("no organisations available; seed organisations first")
	}

	for _, organisationID := range s.data.OrganisationIDs {
		for i := range perOrganisation {
			if i == 0 {
				fmt.Printf("[seed] locations for organisation %s\n", organisationID)
			}
			capacity := int32(gofakeit.Number(6, 60))
			created, err := s.store.CreateLocation(ctx, db.CreateLocationParams{
				OrganisationID:      organisationID,
				Name:                fmt.Sprintf("%s %d", oneOf([]string{"Main", "North", "South", "West", "East"}), i+1),
				Street:              gofakeit.StreetName(),
				HouseNumber:         fmt.Sprintf("%d", gofakeit.Number(1, 350)),
				HouseNumberAddition: nullableString(strings.ToUpper(gofakeit.LetterN(1)), 0.85),
				PostalCode:          fakePostalCodeNL(),
				City:                gofakeit.City(),
				Capacity:            &capacity,
			})
			if err != nil {
				return fmt.Errorf("create location for organisation %s: %w", organisationID, err)
			}

			s.data.LocationIDs = append(s.data.LocationIDs, created.ID)
		}
	}

	return nil
}

func (s *Seeder) SeedSenders(ctx context.Context, count int) error {
	types := []db.SenderTypesEnum{
		db.SenderTypesEnumMainProvider,
		db.SenderTypesEnumLocalAuthority,
		db.SenderTypesEnumParticularParty,
		db.SenderTypesEnumHealthcareInstitution,
	}

	for i := range count {
		if (i+1)%20 == 0 || i == 0 || i+1 == count {
			fmt.Printf("[seed] senders: %d/%d\n", i+1, count)
		}
		name := fmt.Sprintf("%s Sender %d", gofakeit.Company(), i+1)
		email := strings.ToLower(fmt.Sprintf("contact.%d@%s.nl", i+1, sanitizeDomain(name)))
		street := gofakeit.StreetName()
		houseNumber := fmt.Sprintf("%d", gofakeit.Number(1, 350))
		houseNumberAddition := nullableString(strings.ToUpper(gofakeit.LetterN(1)), 0.85)
		postalCode := fakePostalCodeNL()
		city := gofakeit.City()
		land := "Netherlands"
		kvk := fmt.Sprintf("%08d", gofakeit.Number(10000000, 99999999))
		btw := fmt.Sprintf("NL%09dB01", gofakeit.Number(100000000, 999999999))
		phone := fakePhone()
		clientNumber := fmt.Sprintf("CL-%06d", gofakeit.Number(1, 999999))

		created, err := s.store.CreateSender(ctx, db.CreateSenderParams{
			Types:               oneOf(types),
			Name:                name,
			Street:              &street,
			HouseNumber:         &houseNumber,
			HouseNumberAddition: houseNumberAddition,
			PostalCode:          &postalCode,
			City:                &city,
			Land:                &land,
			Kvknumber:           &kvk,
			Btwnumber:           &btw,
			PhoneNumber:         &phone,
			ClientNumber:        &clientNumber,
			EmailAddress:        &email,
			Contacts:            []byte("[]"),
		})
		if err != nil {
			return fmt.Errorf("create sender %d: %w", i+1, err)
		}

		s.data.SenderIDs = append(s.data.SenderIDs, created.ID)
	}

	return nil
}
