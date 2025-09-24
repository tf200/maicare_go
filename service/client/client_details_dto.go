package clientp

import (
	"maicare_go/pagination"
	"time"
)

// Address represents a client address
type Address struct {
	BelongsTo   *string `json:"belongs_to"`
	Address     *string `json:"address"`
	City        *string `json:"city"`
	ZipCode     *string `json:"zip_code"`
	PhoneNumber *string `json:"phone_number"`
}

// CreateClientDetailsRequest represents a request to create a new client
type CreateClientDetailsRequest struct {
	FirstName                  string    `json:"first_name" binding:"required"`
	LastName                   string    `json:"last_name" binding:"required"`
	Email                      string    `json:"email" binding:"required,email"`
	Organisation               *string   `json:"organisation"`
	LocationID                 *int64    `json:"location_id"`
	LegalMeasure               *string   `json:"legal_measure"`
	Birthplace                 *string   `json:"birthplace"`
	Departement                *string   `json:"departement"`
	Gender                     string    `json:"gender" binding:"required oneof=male female other"`
	Filenumber                 string    `json:"filenumber" binding:"required"`
	DateOfBirth                string    `json:"date_of_birth" binding:"required" time_format:"2006-01-02"`
	PhoneNumber                *string   `json:"phone_number" binding:"required"`
	SenderID                   *int64    `json:"sender_id" binding:"required"`
	Infix                      *string   `json:"infix"`
	Source                     *string   `json:"source" binding:"required"`
	Bsn                        *string   `json:"bsn"`
	BsnVerifiedBy              *int64    `json:"bsn_verified_by"` // needs to be checked
	Addresses                  []Address `json:"addresses"`
	EducationCurrentlyEnrolled bool      `json:"education_currently_enrolled"`
	EducationInstitution       *string   `json:"education_institution"`
	EducationMentorName        *string   `json:"education_mentor_name"`
	EducationMentorPhone       *string   `json:"education_mentor_phone"`
	EducationMentorEmail       *string   `json:"education_mentor_email"`
	EducationAdditionalNotes   *string   `json:"education_additional_notes"`
	EducationLevel             *string   `json:"education_level" binding:"oneof=primary secondary higher none"`
	WorkCurrentlyEmployed      bool      `json:"work_currently_employed"`
	WorkCurrentEmployer        *string   `json:"work_current_employer"`
	WorkCurrentEmployerPhone   *string   `json:"work_employer_phone"`
	WorkCurrentEmployerEmail   *string   `json:"work_employer_email"`
	WorkCurrentPosition        *string   `json:"work_current_position"`
	WorkStartDate              time.Time `json:"work_start_date"`
	WorkAdditionalNotes        *string   `json:"work_additional_notes"`
	LivingSituation            *string   `json:"living_situation" binding:"oneof=home foster_care youth_care_institution other"`
	LivingSituationNotes       *string   `json:"living_situation_notes"`
}

// CreateClientDetailsResponse represents a response to a create client request
type CreateClientDetailsResponse struct {
	ID                         int64     `json:"id"`
	FirstName                  string    `json:"first_name"`
	LastName                   string    `json:"last_name"`
	DateOfBirth                time.Time `json:"date_of_birth"`
	Identity                   bool      `json:"identity"`
	Status                     *string   `json:"status"`
	Bsn                        *string   `json:"bsn"`
	BsnVerifiedBy              *int64    `json:"bsn_verified_by"` // needs to be checked
	Source                     *string   `json:"source"`
	Birthplace                 *string   `json:"birthplace"`
	Email                      string    `json:"email"`
	PhoneNumber                *string   `json:"phone_number"`
	Organisation               *string   `json:"organisation"`
	Departement                *string   `json:"departement"`
	Gender                     string    `json:"gender"`
	Filenumber                 string    `json:"filenumber"`
	ProfilePicture             *string   `json:"profile_picture"`
	Infix                      *string   `json:"infix"`
	Created                    time.Time `json:"created"`
	SenderID                   *int64    `json:"sender_id"`
	LocationID                 *int64    `json:"location_id"`
	DepartureReason            *string   `json:"departure_reason"`
	DepartureReport            *string   `json:"departure_report"`
	Addresses                  []Address `json:"addresses"`
	LegalMeasure               *string   `json:"legal_measure"`
	EducationCurrentlyEnrolled bool      `json:"education_currently_enrolled"`
	EducationInstitution       *string   `json:"education_institution"`
	EducationMentorName        *string   `json:"education_mentor_name"`
	EducationMentorEmail       *string   `json:"education_mentor_email"`
	EducationMentorPhone       *string   `json:"education_mentor_phone"`
	EducationAdditionalNotes   *string   `json:"education_additional_notes"`
	EducationLevel             *string   `json:"education_level"`
	WorkCurrentlyEmployed      bool      `json:"work_currently_employed"`
	WorkCurrentEmployer        *string   `json:"work_current_employer"`
	WorkCurrentEmployerPhone   *string   `json:"work_employer_phone"`
	WorkCurrentEmployerEmail   *string   `json:"work_employer_email"`
	WorkCurrentPosition        *string   `json:"work_current_position"`
	WorkStartDate              time.Time `json:"work_start_date"`
	WorkAdditionalNotes        *string   `json:"work_additional_notes"`
	LivingSituation            *string   `json:"living_situation"`
	LivingSituationNotes       *string   `json:"living_situation_notes"`
}

// ListClientsApiParams represents a request to list clients
type ListClientsApiParams struct {
	pagination.Request
	Status     *string `form:"status"`
	LocationID *int64  `form:"location_id"`
	Search     *string `form:"search"`
}

// ListClientsApiResponse represents a response to a list clients request
type ListClientsApiResponse struct {
	ID                    int64     `json:"id"`
	FirstName             string    `json:"first_name"`
	LastName              string    `json:"last_name"`
	DateOfBirth           time.Time `json:"date_of_birth"`
	Identity              bool      `json:"identity"`
	Status                *string   `json:"status"`
	Bsn                   *string   `json:"bsn"`
	Source                *string   `json:"source"`
	Birthplace            *string   `json:"birthplace"`
	Email                 string    `json:"email"`
	PhoneNumber           *string   `json:"phone_number"`
	Organisation          *string   `json:"organisation"`
	Departement           *string   `json:"departement"`
	Gender                string    `json:"gender"`
	Filenumber            string    `json:"filenumber"`
	ProfilePicture        *string   `json:"profile_picture"`
	Infix                 *string   `json:"infix"`
	CreatedAt             time.Time `json:"created_at"`
	SenderID              *int64    `json:"sender_id"`
	LocationID            *int64    `json:"location_id"`
	DepartureReason       *string   `json:"departure_reason"`
	DepartureReport       *string   `json:"departure_report"`
	Addresses             []Address `json:"addresses"`
	LegalMeasure          *string   `json:"legal_measure"`
	HasUntakenMedications bool      `json:"has_untaken_medications"`
}

// GetClientsCountApi gets the count of clients
type GetClientsCountResponse struct {
	TotalClients         int64 `json:"total_clients"`
	ClientsInCare        int64 `json:"clients_in_care"`
	ClientsOnWaitingList int64 `json:"clients_on_waiting_list"`
	ClientsOutOfCare     int64 `json:"clients_out_of_care"`
}

// GetClientApiResponse represents a response to a get client request
type GetClientApiResponse struct {
	ID                         int64     `json:"id"`
	FirstName                  string    `json:"first_name"`
	LastName                   string    `json:"last_name"`
	DateOfBirth                time.Time `json:"date_of_birth"`
	Identity                   bool      `json:"identity"`
	Status                     *string   `json:"status"`
	Bsn                        *string   `json:"bsn"`
	BsnVerifiedBy              *int64    `json:"bsn_verified_by"`
	BsnVerifiedByFirstName     *string   `json:"bsn_verified_by_first_name"`
	BsnVerifiedByLastName      *string   `json:"bsn_verified_by_last_name"`
	Source                     *string   `json:"source"`
	Birthplace                 *string   `json:"birthplace"`
	Email                      string    `json:"email"`
	PhoneNumber                *string   `json:"phone_number"`
	Organisation               *string   `json:"organisation"`
	Departement                *string   `json:"departement"`
	Gender                     string    `json:"gender"`
	Filenumber                 string    `json:"filenumber"`
	ProfilePicture             *string   `json:"profile_picture"`
	Infix                      *string   `json:"infix"`
	CreatedAt                  time.Time `json:"created_at"`
	SenderID                   *int64    `json:"sender_id"`
	LocationID                 *int64    `json:"location_id"`
	DepartureReason            *string   `json:"departure_reason"`
	DepartureReport            *string   `json:"departure_report"`
	LegalMeasure               *string   `json:"legal_measure"`
	HasUntakenMedications      bool      `json:"has_untaken_medications"`
	EducationCurrentlyEnrolled bool      `json:"education_currently_enrolled"`
	EducationInstitution       *string   `json:"education_institution"`
	EducationMentorName        *string   `json:"education_mentor_name"`
	EducationMentorEmail       *string   `json:"education_mentor_email"`
	EducationMentorPhone       *string   `json:"education_mentor_phone"`
	EducationAdditionalNotes   *string   `json:"education_additional_notes"`
	EducationLevel             *string   `json:"education_level"`
	WorkCurrentlyEmployed      bool      `json:"work_currently_employed"`
	WorkCurrentEmployer        *string   `json:"work_current_employer"`
	WorkCurrentEmployerPhone   *string   `json:"work_employer_phone"`
	WorkCurrentEmployerEmail   *string   `json:"work_employer_email"`
	WorkCurrentPosition        *string   `json:"work_current_position"`
	WorkStartDate              time.Time `json:"work_start_date"`
	WorkAdditionalNotes        *string   `json:"work_additional_notes"`
	LivingSituation            *string   `json:"living_situation"`
	LivingSituationNotes       *string   `json:"living_situation_notes"`
}

// GetClientAddressesApiResponse represents a response to a get client addresses request
type GetClientAddressesApiResponse struct {
	Addresses []Address `json:"addresses"`
}

// UpdateClientDetailsRequest represents a request to update a client
type UpdateClientDetailsRequest struct {
	FirstName                  *string   `json:"first_name"`
	LastName                   *string   `json:"last_name"`
	DateOfBirth                time.Time `json:"date_of_birth"`
	Identity                   *bool     `json:"identity"`
	Bsn                        *string   `json:"bsn"`
	BsnVerifiedBy              *int64    `json:"bsn_verified_by"`
	Source                     *string   `json:"source"`
	Birthplace                 *string   `json:"birthplace"`
	Email                      *string   `json:"email"`
	PhoneNumber                *string   `json:"phone_number"`
	Organisation               *string   `json:"organisation"`
	Departement                *string   `json:"departement"`
	Gender                     *string   `json:"gender"`
	Filenumber                 *string   `json:"filenumber"`
	ProfilePicture             *string   `json:"profile_picture"`
	Infix                      *string   `json:"infix"`
	SenderID                   *int64    `json:"sender_id"`
	LocationID                 *int64    `json:"location_id"`
	DepartureReason            *string   `json:"departure_reason"`
	DepartureReport            *string   `json:"departure_report"`
	LegalMeasure               *string   `json:"legal_measure"`
	EducationCurrentlyEnrolled *bool     `json:"education_currently_enrolled"`
	EducationInstitution       *string   `json:"education_institution"`
	EducationMentorName        *string   `json:"education_mentor_name"`
	EducationMentorPhone       *string   `json:"education_mentor_phone"`
	EducationMentorEmail       *string   `json:"education_mentor_email"`
	EducationAdditionalNotes   *string   `json:"education_additional_notes"`
	EducationLevel             *string   `json:"education_level"`
	WorkCurrentlyEmployed      *bool     `json:"work_currently_employed"`
	WorkCurrentEmployer        *string   `json:"work_current_employer"`
	WorkCurrentEmployerPhone   *string   `json:"work_employer_phone"`
	WorkCurrentEmployerEmail   *string   `json:"work_employer_email"`
	WorkCurrentPosition        *string   `json:"work_current_position"`
	WorkStartDate              time.Time `json:"work_start_date"`
	WorkAdditionalNotes        *string   `json:"work_additional_notes"`
	LivingSituation            *string   `json:"living_situation"`
	LivingSituationNotes       *string   `json:"living_situation_notes"`
}

// UpdateClientDetailsResponse represents a response to an update client request
type UpdateClientDetailsResponse struct {
	ID                    int64     `json:"id"`
	FirstName             string    `json:"first_name"`
	LastName              string    `json:"last_name"`
	DateOfBirth           time.Time `json:"date_of_birth"`
	Identity              bool      `json:"identity"`
	Status                *string   `json:"status"`
	Bsn                   *string   `json:"bsn"`
	BsnVerifiedBy         *int64    `json:"bsn_verified_by"`
	Source                *string   `json:"source"`
	Birthplace            *string   `json:"birthplace"`
	Email                 string    `json:"email"`
	PhoneNumber           *string   `json:"phone_number"`
	Organisation          *string   `json:"organisation"`
	Departement           *string   `json:"departement"`
	Gender                string    `json:"gender"`
	Filenumber            string    `json:"filenumber"`
	ProfilePicture        *string   `json:"profile_picture"`
	Infix                 *string   `json:"infix"`
	Created               time.Time `json:"created"`
	SenderID              *int64    `json:"sender_id"`
	LocationID            *int64    `json:"location_id"`
	DepartureReason       *string   `json:"departure_reason"`
	DepartureReport       *string   `json:"departure_report"`
	Addresses             []Address `json:"addresses"`
	LegalMeasure          *string   `json:"legal_measure"`
	HasUntakenMedications bool      `json:"has_untaken_medications"`
}
