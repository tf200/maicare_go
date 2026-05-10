package main

import (
	"fmt"
	"strings"

	db "maicare_go/db/sqlc"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

func oneOf[T any](values []T) T {
	return values[gofakeit.Number(0, len(values)-1)]
}

func boolValue(v *bool) bool {
	return v != nil && *v
}

func deriveCareType(form db.GetRegistrationFormRow) db.IntakeCareTypeEnum {
	switch {
	case boolValue(form.CareProtectedLiving):
		return db.IntakeCareTypeEnumProtectedLiving
	case boolValue(form.CareRoomTrainingCenter):
		return db.IntakeCareTypeEnumTrainingCenter
	case boolValue(form.CareAssistedIndependentLiving):
		return db.IntakeCareTypeEnumSupportedIndependentLiving
	case boolValue(form.CareAmbulatoryGuidance):
		return db.IntakeCareTypeEnumAmbulatorySupport
	default:
		if form.WorkCurrentlyEmployed || form.EducationCurrentlyEnrolled {
			return db.IntakeCareTypeEnumAmbulatorySupport
		}
		return db.IntakeCareTypeEnumProtectedLiving
	}
}

func deriveAdmissionType(form db.GetRegistrationFormRow) db.AdmissionTypeEnum {
	highUrgency := boolValue(form.RiskSuicidalSelfharm) || boolValue(form.RiskAggressiveBehavior) || boolValue(form.RiskWeaponPossession)
	if highUrgency || chance(0.18) {
		return db.AdmissionTypeEnumCrisisAdmission
	}
	return db.AdmissionTypeEnumRegularPlacement
}

func deriveSelfSufficiency(form db.GetRegistrationFormRow) int32 {
	riskCount := 0
	for _, risk := range []*bool{
		form.RiskAggressiveBehavior,
		form.RiskSuicidalSelfharm,
		form.RiskSubstanceAbuse,
		form.RiskPsychiatricIssues,
		form.RiskCriminalHistory,
		form.RiskFlightBehavior,
		form.RiskWeaponPossession,
		form.RiskSexualBehavior,
		form.RiskDayNightRhythm,
		form.RiskOther,
	} {
		if boolValue(risk) {
			riskCount++
		}
	}

	base := int32(4)
	switch {
	case riskCount >= 6:
		base = 1
	case riskCount >= 4:
		base = 2
	case riskCount >= 2:
		base = 3
	}

	if form.WorkCurrentlyEmployed || form.EducationCurrentlyEnrolled {
		base++
	}
	return boundedLevel(base)
}

func boundedLevel(level int32) int32 {
	if level < 1 {
		return 1
	}
	if level > 5 {
		return 5
	}
	return level
}

func pickUniqueTopics(topics []db.Topic, n int) []db.Topic {
	if n >= len(topics) {
		return topics
	}

	chosen := make([]db.Topic, 0, n)
	used := make(map[uuid.UUID]struct{}, n)
	for len(chosen) < n {
		topic := oneOf(topics)
		if _, exists := used[topic.ID]; exists {
			continue
		}
		used[topic.ID] = struct{}{}
		chosen = append(chosen, topic)
	}
	return chosen
}

func buildProposedGoals(topicName string) []map[string]string {
	return []map[string]string{
		{
			"title":       fmt.Sprintf("Improve %s stability", strings.ToLower(topicName)),
			"description": gofakeit.Sentence(10),
			"priority":    oneOf([]string{"medium", "high"}),
		},
		{
			"title":       fmt.Sprintf("Reach next level in %s", strings.ToLower(topicName)),
			"description": gofakeit.Sentence(8),
			"priority":    "medium",
		},
	}
}

func contractSettingsFromIntakeCareType(careType *db.IntakeCareTypeEnum) (db.CareTypeEnum, db.PriceTimeUnitEnum, *float64, *db.HoursTypeEnum) {
	if careType != nil && *careType == db.IntakeCareTypeEnumAmbulatorySupport {
		hours := float64(gofakeit.Number(4, 24))
		hoursType := db.HoursTypeEnumWeekly
		return db.CareTypeEnumAmbulante, db.PriceTimeUnitEnumHourly, &hours, &hoursType
	}

	return db.CareTypeEnumAccommodation, db.PriceTimeUnitEnumWeekly, nil, nil
}

func randomGoals() []string {
	pool := []string{
		"Improve school attendance",
		"Build stable daily routine",
		"Improve communication with guardians",
		"Reduce stress and anxiety",
		"Increase self-reliance in daily tasks",
		"Find suitable education pathway",
		"Develop healthy sleep schedule",
		"Improve social network support",
	}

	goalCount := gofakeit.Number(1, 3)
	goals := make([]string, 0, goalCount)
	used := map[int]struct{}{}

	for len(goals) < goalCount {
		idx := gofakeit.Number(0, len(pool)-1)
		if _, exists := used[idx]; exists {
			continue
		}
		used[idx] = struct{}{}
		goals = append(goals, pool[idx])
	}

	return goals
}
