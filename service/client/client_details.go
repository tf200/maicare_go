package clientp

import (
	"context"
	"encoding/json"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/pagination"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (s *clientService) CreateClientDetails(req CreateClientDetailsRequest, ctx context.Context) (*CreateClientDetailsResponse, error) {
	parsedDateOfBirth, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateClientDetails",
			"Failed to parse date of birth", zap.Error(err))
		return nil, fmt.Errorf("failed to parse date of birth")
	}
	AddressesJSON, err := json.Marshal(req.Addresses)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateClientDetails",
			"Failed to marshal addresses", zap.Error(err))
		return nil, fmt.Errorf("failed to marshal addresses")
	}

	client, err := s.Store.CreateClientDetails(ctx, db.CreateClientDetailsParams{
		FirstName:                  req.FirstName,
		LastName:                   req.LastName,
		DateOfBirth:                pgtype.Date{Time: parsedDateOfBirth, Valid: true},
		Identity:                   true,
		Bsn:                        req.Bsn,
		BsnVerifiedBy:              req.BsnVerifiedBy,
		Source:                     req.Source,
		Birthplace:                 req.Birthplace,
		Email:                      req.Email,
		PhoneNumber:                req.PhoneNumber,
		Organisation:               req.Organisation,
		Departement:                req.Departement,
		Infix:                      req.Infix,
		Gender:                     req.Gender,
		Filenumber:                 req.Filenumber,
		SenderID:                   req.SenderID,
		LocationID:                 req.LocationID,
		Addresses:                  AddressesJSON,
		LegalMeasure:               req.LegalMeasure,
		EducationCurrentlyEnrolled: req.EducationCurrentlyEnrolled,
		EducationInstitution:       req.EducationInstitution,
		EducationMentorName:        req.EducationMentorName,
		EducationMentorPhone:       req.EducationMentorPhone,
		EducationMentorEmail:       req.EducationMentorEmail,
		EducationAdditionalNotes:   req.EducationAdditionalNotes,
		EducationLevel:             req.EducationLevel,
		WorkCurrentlyEmployed:      req.WorkCurrentlyEmployed,
		WorkCurrentEmployer:        req.WorkCurrentEmployer,
		WorkCurrentEmployerPhone:   req.WorkCurrentEmployerPhone,
		WorkCurrentEmployerEmail:   req.WorkCurrentEmployerEmail,
		WorkCurrentPosition:        req.WorkCurrentPosition,
		WorkStartDate:              pgtype.Date{Time: req.WorkStartDate, Valid: true},
		WorkAdditionalNotes:        req.WorkAdditionalNotes,
		LivingSituation:            req.LivingSituation,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateClientDetails",
			"Failed to create client details", zap.Error(err))
		return nil, fmt.Errorf("failed to create client details")
	}

	var addresses []Address
	err = json.Unmarshal(client.Addresses, &addresses)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateClientDetails",
			"Failed to unmarshal addresses", zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal addresses")
	}

	result := &CreateClientDetailsResponse{
		ID:                         client.ID,
		FirstName:                  client.FirstName,
		LastName:                   client.LastName,
		DateOfBirth:                client.DateOfBirth.Time,
		Identity:                   client.Identity,
		Status:                     client.Status,
		Bsn:                        client.Bsn,
		Source:                     client.Source,
		Birthplace:                 client.Birthplace,
		Email:                      client.Email,
		PhoneNumber:                client.PhoneNumber,
		Organisation:               client.Organisation,
		Departement:                client.Departement,
		Gender:                     client.Gender,
		Filenumber:                 client.Filenumber,
		ProfilePicture:             client.ProfilePicture,
		Infix:                      client.Infix,
		Created:                    client.CreatedAt.Time,
		SenderID:                   client.SenderID,
		LocationID:                 client.LocationID,
		DepartureReason:            client.DepartureReason,
		DepartureReport:            client.DepartureReport,
		Addresses:                  addresses,
		LegalMeasure:               client.LegalMeasure,
		EducationCurrentlyEnrolled: client.EducationCurrentlyEnrolled,
		EducationInstitution:       client.EducationInstitution,
		EducationMentorName:        client.EducationMentorName,
		EducationMentorPhone:       client.EducationMentorPhone,
		EducationMentorEmail:       client.EducationMentorEmail,
		EducationAdditionalNotes:   client.EducationAdditionalNotes,
		EducationLevel:             client.EducationLevel,
		WorkCurrentlyEmployed:      client.WorkCurrentlyEmployed,
		WorkCurrentEmployer:        client.WorkCurrentEmployer,
		WorkCurrentEmployerPhone:   client.WorkCurrentEmployerPhone,
		WorkCurrentEmployerEmail:   client.WorkCurrentEmployerEmail,
		WorkCurrentPosition:        client.WorkCurrentPosition,
		WorkStartDate:              client.WorkStartDate.Time,
		WorkAdditionalNotes:        client.WorkAdditionalNotes,
		LivingSituation:            client.LivingSituation,
		LivingSituationNotes:       client.LivingSituationNotes,
	}

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "CreateClientDetails",
		"Successfully created client details", zap.Int64("ClientID", client.ID))
	return result, nil
}

func (s *clientService) ListClientDetails(ctx *gin.Context, req ListClientsApiParams) (*pagination.Response[ListClientsApiResponse], error) {
	params := req.GetParams()

	clients, err := s.Store.ListClientDetails(ctx, db.ListClientDetailsParams{
		Limit:      params.Limit,
		Offset:     params.Offset,
		Status:     req.Status,
		LocationID: req.LocationID,
		Search:     req.Search,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "ListClientDetails",
			"Failed to list client details", zap.Error(err))
		return nil, fmt.Errorf("failed to list client details")
	}
	if len(clients) == 0 {
		s.Logger.LogBusinessEvent(logger.LogLevelInfo, "ListClientDetails",
			"No clients found")
		pagObj := pagination.NewResponse(ctx, req.Request, []ListClientsApiResponse{}, 0)
		return &pagObj, nil
	}

	totalCount := clients[0].TotalCount

	clientList := make([]ListClientsApiResponse, len(clients))
	for i, client := range clients {
		var addresses []Address
		err = json.Unmarshal(client.Addresses, &addresses)
		if err != nil {
			s.Logger.LogBusinessEvent(logger.LogLevelError, "ListClientDetails",
				"Failed to unmarshal addresses", zap.Error(err))
			return nil, fmt.Errorf("failed to unmarshal addresses")
		}
		clientList[i] = ListClientsApiResponse{
			ID:                    client.ID,
			FirstName:             client.FirstName,
			LastName:              client.LastName,
			DateOfBirth:           client.DateOfBirth.Time,
			Identity:              client.Identity,
			Status:                client.Status,
			Bsn:                   client.Bsn,
			Source:                client.Source,
			Birthplace:            client.Birthplace,
			Email:                 client.Email,
			PhoneNumber:           client.PhoneNumber,
			Organisation:          client.Organisation,
			Departement:           client.Departement,
			Gender:                client.Gender,
			Filenumber:            client.Filenumber,
			ProfilePicture:        s.GenerateResponsePresignedURL(client.ProfilePicture, ctx),
			Infix:                 client.Infix,
			CreatedAt:             client.CreatedAt.Time,
			SenderID:              client.SenderID,
			LocationID:            client.LocationID,
			DepartureReason:       client.DepartureReason,
			DepartureReport:       client.DepartureReport,
			Addresses:             addresses,
			LegalMeasure:          client.LegalMeasure,
			HasUntakenMedications: client.HasUntakenMedications,
		}
	}
	pagObj := pagination.NewResponse(ctx, req.Request, clientList, totalCount)
	return &pagObj, nil
}

func (s *clientService) GetClientsCount(ctx context.Context) (*GetClientsCountResponse, error) {
	count, err := s.Store.GetClientCounts(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GetClientsCount",
			"Failed to get clients count", zap.Error(err))
		return nil, fmt.Errorf("failed to get clients count")
	}
	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "GetClientsCount",
		"Successfully retrieved clients count")
	return &GetClientsCountResponse{
		TotalClients:         count.TotalClients,
		ClientsInCare:        count.ClientsInCare,
		ClientsOnWaitingList: count.ClientsOnWaitingList,
		ClientsOutOfCare:     count.ClientsOutOfCare,
	}, nil
}

func (s *clientService) GetClientDetails(ctx context.Context, clientID int64) (*GetClientApiResponse, error) {
	client, err := s.Store.GetClientDetails(ctx, clientID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GetClientDetails",
			"Failed to get client details", zap.Error(err), zap.Int64("ClientID", clientID))
		return nil, fmt.Errorf("failed to get client details")
	}

	return &GetClientApiResponse{
		ID:                         client.ID,
		FirstName:                  client.FirstName,
		LastName:                   client.LastName,
		DateOfBirth:                client.DateOfBirth.Time,
		Identity:                   client.Identity,
		Status:                     client.Status,
		Bsn:                        client.Bsn,
		BsnVerifiedBy:              client.BsnVerifiedBy,
		BsnVerifiedByFirstName:     client.BsnVerifiedByFirstName,
		BsnVerifiedByLastName:      client.BsnVerifiedByLastName,
		Source:                     client.Source,
		Birthplace:                 client.Birthplace,
		Email:                      client.Email,
		PhoneNumber:                client.PhoneNumber,
		Organisation:               client.Organisation,
		Departement:                client.Departement,
		Gender:                     client.Gender,
		Filenumber:                 client.Filenumber,
		ProfilePicture:             s.GenerateResponsePresignedURL(client.ProfilePicture, ctx),
		Infix:                      client.Infix,
		CreatedAt:                  client.CreatedAt.Time,
		SenderID:                   client.SenderID,
		LocationID:                 client.LocationID,
		DepartureReason:            client.DepartureReason,
		DepartureReport:            client.DepartureReport,
		LegalMeasure:               client.LegalMeasure,
		HasUntakenMedications:      client.HasUntakenMedications,
		EducationCurrentlyEnrolled: client.EducationCurrentlyEnrolled,
		EducationInstitution:       client.EducationInstitution,
		EducationMentorName:        client.EducationMentorName,
		EducationMentorEmail:       client.EducationMentorEmail,
		EducationMentorPhone:       client.EducationMentorPhone,
		EducationAdditionalNotes:   client.EducationAdditionalNotes,
		EducationLevel:             client.EducationLevel,
		WorkCurrentlyEmployed:      client.WorkCurrentlyEmployed,
		WorkCurrentEmployer:        client.WorkCurrentEmployer,
		WorkCurrentEmployerPhone:   client.WorkCurrentEmployerPhone,
		WorkCurrentEmployerEmail:   client.WorkCurrentEmployerEmail,
		WorkCurrentPosition:        client.WorkCurrentPosition,
		WorkStartDate:              client.WorkStartDate.Time,
		WorkAdditionalNotes:        client.WorkAdditionalNotes,
		LivingSituation:            client.LivingSituation,
		LivingSituationNotes:       client.LivingSituationNotes,
	}, nil

}

func (s *clientService) GetClientAddresses(ctx context.Context, clientID int64) (*GetClientAddressesApiResponse, error) {
	address, err := s.Store.GetClientAddresses(ctx, clientID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GetClientAddresses",
			"Failed to get client addresses", zap.Error(err), zap.Int64("ClientID", clientID))
		return nil, fmt.Errorf("failed to get client addresses")
	}

	var addresses []Address
	err = json.Unmarshal(address, &addresses)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GetClientAddresses",
			"Failed to unmarshal addresses", zap.Error(err), zap.Int64("ClientID", clientID))
		return nil, fmt.Errorf("failed to unmarshal addresses")
	}

	return &GetClientAddressesApiResponse{
		Addresses: addresses,
	}, nil
}

func (s *clientService) UpdateClientDetails(ctx context.Context, req UpdateClientDetailsRequest, clientID int64) (*UpdateClientDetailsResponse, error) {
	client, err := s.Store.UpdateClientDetails(ctx, db.UpdateClientDetailsParams{
		ID:                         clientID,
		FirstName:                  req.FirstName,
		LastName:                   req.LastName,
		DateOfBirth:                pgtype.Date{Time: req.DateOfBirth, Valid: true},
		Identity:                   req.Identity,
		Bsn:                        req.Bsn,
		BsnVerifiedBy:              req.BsnVerifiedBy,
		Source:                     req.Source,
		Birthplace:                 req.Birthplace,
		Email:                      req.Email,
		PhoneNumber:                req.PhoneNumber,
		Organisation:               req.Organisation,
		Departement:                req.Departement,
		Gender:                     req.Gender,
		Filenumber:                 req.Filenumber,
		ProfilePicture:             req.ProfilePicture,
		Infix:                      req.Infix,
		SenderID:                   req.SenderID,
		LocationID:                 req.LocationID,
		DepartureReason:            req.DepartureReason,
		DepartureReport:            req.DepartureReport,
		LegalMeasure:               req.LegalMeasure,
		EducationCurrentlyEnrolled: req.EducationCurrentlyEnrolled,
		EducationInstitution:       req.EducationInstitution,
		EducationMentorName:        req.EducationMentorName,
		EducationMentorPhone:       req.EducationMentorPhone,
		EducationMentorEmail:       req.EducationMentorEmail,
		EducationAdditionalNotes:   req.EducationAdditionalNotes,
		EducationLevel:             req.EducationLevel,
		WorkCurrentlyEmployed:      req.WorkCurrentlyEmployed,
		WorkCurrentEmployer:        req.WorkCurrentEmployer,
		WorkCurrentEmployerPhone:   req.WorkCurrentEmployerPhone,
		WorkCurrentEmployerEmail:   req.WorkCurrentEmployerEmail,
		WorkCurrentPosition:        req.WorkCurrentPosition,
		WorkStartDate:              pgtype.Date{Time: req.WorkStartDate, Valid: true},
		WorkAdditionalNotes:        req.WorkAdditionalNotes,
		LivingSituation:            req.LivingSituation,
		LivingSituationNotes:       req.LivingSituationNotes,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "UpdateClientDetails",
			"Failed to update client details", zap.Error(err), zap.Int64("ClientID", clientID))
		return nil, fmt.Errorf("failed to update client details")
	}

	var addresses []Address
	err = json.Unmarshal(client.Addresses, &addresses)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "UpdateClientDetails",
			"Failed to unmarshal addresses", zap.Error(err), zap.Int64("ClientID", clientID))
		return nil, fmt.Errorf("failed to unmarshal addresses")
	}

	result := &UpdateClientDetailsResponse{
		ID:                    client.ID,
		FirstName:             client.FirstName,
		LastName:              client.LastName,
		DateOfBirth:           client.DateOfBirth.Time,
		Identity:              client.Identity,
		Status:                client.Status,
		Bsn:                   client.Bsn,
		BsnVerifiedBy:         client.BsnVerifiedBy,
		Source:                client.Source,
		Birthplace:            client.Birthplace,
		Email:                 client.Email,
		PhoneNumber:           client.PhoneNumber,
		Organisation:          client.Organisation,
		Departement:           client.Departement,
		Gender:                client.Gender,
		Filenumber:            client.Filenumber,
		ProfilePicture:        client.ProfilePicture,
		Infix:                 client.Infix,
		SenderID:              client.SenderID,
		LocationID:            client.LocationID,
		DepartureReason:       client.DepartureReason,
		DepartureReport:       client.DepartureReport,
		Addresses:             addresses,
		LegalMeasure:          client.LegalMeasure,
		HasUntakenMedications: client.HasUntakenMedications,
	}
	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "UpdateClientDetails",
		"Successfully updated client details", zap.Int64("ClientID", client.ID))
	return result, nil
}
