package pdf

import (
	"context"
	"log"
	"maicare_go/bucket"
	"maicare_go/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateAppointmentcardPDF(t *testing.T) {

	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatalf("Could not load conf %v", err)
	}

	testb2Client, err := bucket.NewB2Client(config)
	if err != nil {
		log.Fatal("cannot create b2 client:", err)
	}

	testCard := AppointmentCard{
		ClientName: "Chrystal Reyes",
		Date:       "16-11-2023",
		Mentor:     "Sanae Aktitou",
		GeneralInformation: []string{
			"Chrystal is geboren op 01-09-2006 (16 jaar jong)",
			"Chrystal blijft bereikbaar en houdt de begeleiding op de hoogte van waar hij is (Tel. Chrystal: (+31 648736925)",
			"Chrystal hoort om 22:30 uur terug te zijn op de woonlocatie",
		},
		ImportantContacts: []string{
			"Jeugdbeschermer Junesca: 0652713825",
			"Moeder van Chrystal: 0684701284",
			"Schooldocent Saul: 0643239809",
			"School Trajectbegeleider Martin: 0643270037",
			"Mentor Diversitas Sanae: 0623502061",
		},
		HouseholdInfo: []string{
			"Chrystal helpt mee met de wekelijkse huisschoonmaak en voert de taken uit die de begeleiding haar oplegt",
			"Chrystal mag onder begeleiding gebruik maken van de wasmachine en droger om haar kleding te wassen",
			"Chrystal mag onder begeleiding haar ontbijt en lunch maaltijden bereiden, maar ruimt de keuken na gebruik zelfstandig weer op",
		},
		OrganizationAgreements: []string{
			"Chrystal meldt zich ziek bij de begeleiding, vanuit Diversitas wordt Chrystal ziekgemeld",
			"Chrystal maakt zich niet schuldig aan een strafbaar feit",
			"Chrystal gaat naar school of stage volgens het rooster",
			"Chrystal is voor binnenkomsttijd 22.30 uur terug op de groep",
		},
		YouthOfficerAgreements: []string{
			"Chrystal maakt zich niet schuldig aan een strafbaar feit",
			"Chrystal gaat naar school of stage volgens het rooster",
			"Chrystal is voor binnenkomsttijd 22.30 uur terug op de groep",
		},
		TreatmentAgreements: []string{
			"Medicatie NVT",
		},
		SmokingRules: []string{
			"Jeugdige mag roken in de tuin (niet in kamer)",
		},
		Work: []string{
			"Chrystal heeft op dit moment geen werk",
		},
		SchoolInternship: []string{
			"Chrystal gaat naar school: Educatief centrum te Rotterdam (Schiehaven)",
			"Chrystal heeft op dit moment geen stage",
			"Chrystal heeft les van ma t/m do van 09.00 uur tot 14.30 en vrij van 09.00 uur tot 13.00 uur",
		},
		Travel: []string{
			"Chrystal haar reisabonnement wordt gefinancierd vanuit Diversitas Zorg",
		},
		Leave: []string{
			"Verlof gaat in contact met ouders, begeleiding en jeugdbeschermer",
			"Chrystal mag van zaterdag tot zondag met verlof bij moeder",
		},
	}

	pdfBytes, err := GenerateAndUploadAppointmentCardPDF(context.Background(), testCard, testb2Client)
	require.NoError(t, err)
	require.NotEmpty(t, pdfBytes)

}
