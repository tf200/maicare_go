package api

import (
	"bytes"
	"maicare_go/util"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

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
