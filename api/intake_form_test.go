package api

import (
	"bytes"
	"context"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
	"maicare_go/token"
	"maicare_go/util"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomIntakeForm(t *testing.T) db.IntakeForm {
	// Helper function to get random enum value

	idTypes := []string{"passport", "id_card", "residence_permit"}
	signedByOptions := []string{"Referrer", "Parent/Guardian", "Client"}
	lawTypes := []string{"Youth Act", "WLZ", "WMO", "Other"}
	registrationTypes := []string{"Protected Living", "Supervised Independent Living", "Outpatient Guidance"}
	livingSituations := []string{"Home", "Foster care", "Youth care institution", "Other"}

	arg := db.CreateIntakeFormParams{
		FirstName: faker.FirstName(),
		LastName:  faker.LastName(),
		DateOfBirth: pgtype.Date{
			Time:  time.Now().AddDate(-20, 0, 0), // 20 years ago
			Valid: true,
		},
		Nationality: faker.Word(),
		Bsn:         faker.CCNumber(),
		Address:     faker.GetRealAddress().Address,
		City:        faker.GetRealAddress().City,
		PostalCode:  faker.GetRealAddress().PostalCode,
		PhoneNumber: faker.Phonenumber(),
		Gender:      faker.Gender(),
		Email:       faker.Email(),
		IDType:      util.RandomEnum(idTypes),
		IDNumber:    faker.CCNumber(),

		ReferrerName:         util.StringPtr(faker.Name()),
		ReferrerOrganization: util.StringPtr(faker.Name()),
		ReferrerFunction:     util.StringPtr(faker.NAME),
		ReferrerPhone:        util.StringPtr(faker.Phonenumber()),
		ReferrerEmail:        util.StringPtr(faker.Email()),
		SignedBy:             util.StringPtr(util.RandomEnum(signedByOptions)),

		HasValidIndication:  true,
		LawType:             util.StringPtr(util.RandomEnum(lawTypes)),
		MainProviderName:    util.StringPtr(faker.Name()),
		MainProviderContact: util.StringPtr(faker.Phonenumber()),
		IndicationStartDate: pgtype.Date{
			Time:  time.Now(),
			Valid: true,
		},
		IndicationEndDate: pgtype.Date{
			Time:  time.Now().AddDate(1, 0, 0),
			Valid: true,
		},
		RegistrationReason: util.StringPtr(faker.Sentence()),
		GuidanceGoals:      util.StringPtr(faker.Sentence()),
		RegistrationType:   util.StringPtr(util.RandomEnum(registrationTypes)),

		LivingSituation:   util.StringPtr(util.RandomEnum(livingSituations)),
		ParentalAuthority: false,
		CurrentSchool:     util.StringPtr(faker.Word()),
		MentorName:        util.StringPtr(faker.Name()),
		MentorPhone:       util.StringPtr(faker.Phonenumber()),
		MentorEmail:       util.StringPtr(faker.Email()),
		PreviousCare:      util.StringPtr(faker.Sentence()),

		GuardianDetails: []byte(`[{
			"first_name": "` + faker.FirstName() + `",
			"last_name": "` + faker.LastName() + `",
			"phone_number": "` + faker.Phonenumber() + `",
			"email": "` + faker.Email() + `",
			"address": "` + faker.GetRealAddress().City + `"
		}]`),

		UsesMedication:      util.RandomBool(),
		AddictionIssues:     util.RandomBool(),
		JudicialInvolvement: util.RandomBool(),

		RiskAggression:       util.RandomBool(),
		RiskSuicidality:      util.RandomBool(),
		RiskRunningAway:      util.RandomBool(),
		RiskSelfHarm:         util.RandomBool(),
		RiskWeaponPossession: util.RandomBool(),
		RiskDrugDealing:      util.RandomBool(),
		OtherRisks:           util.StringPtr(faker.Sentence()),

		SharingPermission: util.RandomBool(),
		TruthDeclaration:  util.RandomBool(),
		ClientSignature:   util.RandomBool(),
		GuardianSignature: util.BoolPtr(util.RandomBool()),
		ReferrerSignature: util.BoolPtr(util.RandomBool()),
		SignatureDate: pgtype.Date{
			Time:  time.Now(),
			Valid: true,
		},
		UrgencyScore: util.RandomEnum([]string{"low", "medium", "high"}),
	}

	form, err := testStore.CreateIntakeForm(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, form)
	require.NotEmpty(t, form.ID)
	require.Equal(t, arg.FirstName, form.FirstName)
	require.Equal(t, arg.LastName, form.LastName)
	require.Equal(t, arg.Email, form.Email)
	return form

}

func TestIntakeFormUploadHandlerApi(t *testing.T) {
	filename, fileContent := createRandomFile(t)

	testCases := []struct {
		name          string
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildRequest: func() (*http.Request, error) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)

				part, err := writer.CreateFormFile("file", filename)
				if err != nil {
					return nil, err
				}

				_, err = part.Write(fileContent)
				if err != nil {
					return nil, err
				}

				err = writer.Close()
				if err != nil {
					return nil, err
				}
				url := "/intake_form/upload"
				req, err := http.NewRequest(http.MethodPost, url, body)
				if err != nil {
					return nil, err
				}
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})

	}
}

func TestCreateIntakeFormApi(t *testing.T) {
	getRandomEnum := func(values []string) string {
		r, _ := faker.RandomInt(0, len(values)-1)
		return values[r[0]]
	}

	idTypes := []string{"passport", "id_card", "residence_permit"}
	signedByOptions := []string{"Referrer", "Parent/Guardian", "Client"}
	lawTypes := []string{"Youth Act", "WLZ", "WMO", "Other"}
	registrationTypes := []string{"Protected Living", "Supervised Independent Living", "Outpatient Guidance"}
	livingSituations := []string{"Home", "Foster care", "Youth care institution", "Other"}
	urgencyScore := []string{"low", "medium", "high"}

	testCases := []struct {
		name          string
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildRequest: func() (*http.Request, error) {
				reqBody := CreateIntakeFormRequest{
					FirstName:   faker.FirstName(),
					LastName:    faker.LastName(),
					DateOfBirth: time.Now().AddDate(-20, 0, 0), // 20 years ago
					Nationality: faker.Word(),
					Bsn:         faker.CCNumber(),
					Address:     faker.GetRealAddress().Address,
					City:        faker.GetRealAddress().City,
					PostalCode:  faker.GetRealAddress().PostalCode,
					PhoneNumber: faker.Phonenumber(),
					Gender:      faker.Gender(),
					Email:       faker.Email(),
					IDType:      getRandomEnum(idTypes),
					IDNumber:    faker.CCNumber(),

					ReferrerName:         util.StringPtr(faker.Name()),
					ReferrerOrganization: util.StringPtr(faker.Name()),
					ReferrerFunction:     util.StringPtr(faker.NAME),
					ReferrerPhone:        util.StringPtr(faker.Phonenumber()),
					ReferrerEmail:        util.StringPtr(faker.Email()),
					SignedBy:             util.StringPtr(getRandomEnum(signedByOptions)),

					HasValidIndication:  true,
					LawType:             util.StringPtr(getRandomEnum(lawTypes)),
					MainProviderName:    util.StringPtr(faker.Name()),
					MainProviderContact: util.StringPtr(faker.Phonenumber()),
					IndicationStartDate: time.Now(),
					IndicationEndDate:   time.Now().AddDate(1, 0, 0),

					RegistrationReason: util.StringPtr(faker.Sentence()),
					GuidanceGoals:      util.StringPtr(faker.Sentence()),
					RegistrationType:   util.StringPtr(getRandomEnum(registrationTypes)),

					LivingSituation:   util.StringPtr(getRandomEnum(livingSituations)),
					ParentalAuthority: false,
					CurrentSchool:     util.StringPtr(faker.Word()),
					MentorName:        util.StringPtr(faker.Name()),
					MentorPhone:       util.StringPtr(faker.Phonenumber()),
					MentorEmail:       util.StringPtr(faker.Email()),
					PreviousCare:      util.StringPtr(faker.Sentence()),

					GuardianDetails: []GuardionInfo{
						{
							FirstName:   faker.FirstName(),
							LastName:    faker.LastName(),
							PhoneNumber: faker.Phonenumber(),
							Email:       faker.Email(),
							Address:     faker.GetRealAddress().Address,
						},
					},

					UsesMedication:      util.RandomBool(),
					AddictionIssues:     util.RandomBool(),
					JudicialInvolvement: util.RandomBool(),

					RiskAggression:       util.RandomBool(),
					RiskSuicidality:      util.RandomBool(),
					RiskRunningAway:      util.RandomBool(),
					RiskSelfHarm:         util.RandomBool(),
					RiskWeaponPossession: util.RandomBool(),
					RiskDrugDealing:      util.RandomBool(),
					OtherRisks:           util.StringPtr(faker.Sentence()),

					SharingPermission: util.RandomBool(),
					TruthDeclaration:  util.RandomBool(),
					ClientSignature:   util.RandomBool(),
					GuardianSignature: util.BoolPtr(util.RandomBool()),
					ReferrerSignature: util.BoolPtr(util.RandomBool()),
					SignatureDate:     time.Now(),
					AttachementIds:    []uuid.UUID{},
					UrgencyScore:      getRandomEnum(urgencyScore),
				}
				data, err := json.Marshal(reqBody)
				require.NoError(t, err)
				url := "/intake_form"
				req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
				require.NoError(t, err)
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})

	}
}

func TestListIntakeFormsApi(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomIntakeForm(t)
	}

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				url := "/intake_form?page=1&page_size=5"
				req, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var res Response[pagination.Response[ListIntakeFormsResponse]]
				err := json.Unmarshal(recorder.Body.Bytes(), &res)
				require.NoError(t, err)
				require.Len(t, res.Data.Results, 5)

			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)
			tc.setupAuth(t, req, testServer.tokenMaker)

			recorder := httptest.NewRecorder()
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})

	}
}

func TestGetIntakeForm(t *testing.T) {
	form := createRandomIntakeForm(t)
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				url := fmt.Sprintf("/intake_form/%d", form.ID)
				req, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var res Response[GetIntakeFormResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &res)
				require.NoError(t, err)
				require.Equal(t, form.ID, res.Data.ID)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)
			tc.setupAuth(t, req, testServer.tokenMaker)

			recorder := httptest.NewRecorder()
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})

	}
}

func TestAddUrgencyScoreApi(t *testing.T) {
	form := createRandomIntakeForm(t)
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				reqBody := AddUrgencyScoreRequest{
					UrgencyScore: "medium",
				}
				data, err := json.Marshal(reqBody)
				require.NoError(t, err)
				url := fmt.Sprintf("/intake_form/%d/urgency_score", form.ID)
				req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
				require.NoError(t, err)
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)
			tc.setupAuth(t, req, testServer.tokenMaker)

			recorder := httptest.NewRecorder()
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})

	}
}

func TestMoveToWaitingList(t *testing.T) {
	form := createRandomIntakeForm(t)
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {

				url := fmt.Sprintf("/intake_form/%d/move_to_waiting_list", form.ID)
				req, err := http.NewRequest(http.MethodPost, url, nil)
				require.NoError(t, err)
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)
			tc.setupAuth(t, req, testServer.tokenMaker)

			recorder := httptest.NewRecorder()
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})

	}
}
