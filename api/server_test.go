package api

import (
	"context"
	grpclient "maicare_go/grpclient/proto"
)

type MockGrpcClient struct {
	CarePlanResponse *grpclient.PersonalizedCarePlanResponse
	SpellingResponse *grpclient.CorrectSpellingResponse
	Err              error
}

func (m *MockGrpcClient) GenerateCarePlan(ctx context.Context, req *grpclient.PersonalizedCarePlanRequest) (*grpclient.PersonalizedCarePlanResponse, error) {
	return m.CarePlanResponse, m.Err
}

func (m *MockGrpcClient) CorrectSpelling(ctx context.Context, req *grpclient.CorrectSpellingRequest) (*grpclient.CorrectSpellingResponse, error) {
	return m.SpellingResponse, m.Err
}

func (m *MockGrpcClient) Close() error {
	return nil
}

func CreateMockGrpcClient() *MockGrpcClient {
	return &MockGrpcClient{
		CarePlanResponse: &grpclient.PersonalizedCarePlanResponse{
			ResourcesRequired: []string{
				"Certified youth financial advisor",
				"Financial literacy resources (digital tools, workbooks)",
				"Supportive foster parents or guardians for daily financial oversight",
			},
			SuccessMetrics: []*grpclient.SuccessMetric{
				{
					Metric:            "New debts accumulated",
					Target:            "Zero new debts monthly",
					MeasurementMethod: "Analysis of financial records and client reports",
				},
				{
					Metric:            "Participation in financial counseling",
					Target:            "Regular weekly or bi-weekly counselor sessions",
					MeasurementMethod: "Counseling service attendance logs",
				},
				{
					Metric:            "Budget compliance",
					Target:            "Adhere to budget for 3 out of 4 weeks",
					MeasurementMethod: "Evaluation of income and expense tracking",
				},
			},
			RiskFactors: []*grpclient.RiskFactor{
				{
					Risk:       "Spending driven by peer influence",
					Mitigation: "Teach strategies to resist peer pressure, engage client in goal-setting, and promote delayed gratification.",
				},
				{
					Risk:       "Difficulty grasping complex financial concepts",
					Mitigation: "Use clear language, visual aids, and break concepts into smaller parts with practical examples.",
				},
				{
					Risk:       "Restricted access to personal income or resources",
					Mitigation: "Identify part-time job opportunities, allowances, or youth savings programs, and seek agency support.",
				},
			},
			SupportNetwork: []*grpclient.SupportRole{
				{
					Role:           "Foster Parents/Guardians",
					Responsibility: "Offer daily oversight, guide spending decisions, and create a safe space for financial learning.",
				},
				{
					Role:           "Case Manager/Youth Advocate",
					Responsibility: "Arrange access to counseling, track progress, and provide emotional support.",
				},
				{
					Role:           "Financial Counselor/Educator",
					Responsibility: "Deliver expert guidance, create debt management plans, and teach financial skills.",
				},
			},
			EmergencyProtocols: []string{},
			ClientProfile: &grpclient.ClientProfile{
				Age:              16,
				LivingSituation:  "Foster care",
				EducationLevel:   "VMBO",
				AssessmentDomain: "Financial Management",
				CurrentLevel:     1,
				LevelDescription: "Increasingly complex debt issues",
			},
			AssessmentSummary: "The client, a 16-year-old male in foster care, is enrolled in VMBO education. He faces significant financial challenges, with growing complex debts classified as Level 1 in financial management.",
			CarePlanObjectives: &grpclient.CarePlanObjectives{
				ShortTermGoals: []*grpclient.Goal{
					{
						GoalTitle:   "Address immediate financial strain",
						Description: "Halt debt growth and reduce urgent financial pressures.",
						Timeframe:   "1-3 months",
						SpecificActions: []string{
							"Compile a list of all creditors and debts owed.",
							"Stop all non-essential spending and adopt a needs-based budget.",
							"Discuss available support for basic needs with foster parents or guardians.",
						},
					},
				},
				LongTermGoals: []*grpclient.Goal{
					{
						GoalTitle:   "Build sustainable financial practices",
						Description: "Develop skills for independent financial management, including saving and credit awareness.",
						Timeframe:   "6-12 months",
						SpecificActions: []string{
							"Follow a tailored personal budget.",
							"Learn saving principles and save toward a goal (e.g., education or driving lessons).",
							"Understand credit implications and responsible borrowing practices.",
						},
					},
				},
				MediumTermGoals: []*grpclient.Goal{
					{
						GoalTitle:   "Begin debt restructuring and financial education",
						Description: "Start formal debt management and learn basics of income and expenses.",
						Timeframe:   "3-6 months",
						SpecificActions: []string{
							"Enroll in youth-focused debt counseling services.",
							"Maintain a basic income and expense tracker.",
							"Learn about income sources (e.g., allowances, part-time work) and typical expenses (e.g., transport, personal items).",
						},
					},
				},
			},
			Interventions: &grpclient.Interventions{
				DailyActivities: []string{
					"Review daily spending with a foster parent or guardian, if applicable.",
					"Log all income and expenses.",
					"Practice setting small financial boundaries, such as avoiding impulse purchases.",
				},
				WeeklyActivities: []string{
					"Attend a financial literacy session with a mentor or case manager.",
					"Analyze income and expense tracker for patterns.",
					"Address financial concerns with a trusted adult.",
				},
				MonthlyActivities: []string{
					"Meet with a financial counselor to assess progress and plan next steps.",
					"Set a measurable financial goal for the upcoming month.",
					"Join a youth-oriented financial literacy workshop.",
				},
			},
			TransitionCriteria: nil,
		},
		SpellingResponse: &grpclient.CorrectSpellingResponse{
			CorrectedText: "This is a sample corrected text.",
		},
		Err: nil,
	}
}
