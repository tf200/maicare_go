package domain

import (
	"strings"
)

// PermissionKey represents a system permission string key (e.g. "CLIENT.VIEW").
type PermissionKey string

func (p PermissionKey) String() string {
	return string(p)
}

// Appointment Permissions
const (
	PermAppointmentCreate             PermissionKey = "APPOINTMENT.CREATE"
	PermAppointmentDelete             PermissionKey = "APPOINTMENT.DELETE"
	PermAppointmentUpdate             PermissionKey = "APPOINTMENT.UPDATE"
	PermAppointmentView               PermissionKey = "APPOINTMENT.VIEW"
	PermAppointmentViewAll            PermissionKey = "APPOINTMENT.VIEW_ALL"
	PermAppointmentWorkApprovalUpdate PermissionKey = "APPOINTMENT.WORK_APPROVAL.UPDATE"
)

// Appointment Card Permissions
const (
	PermAppointmentCardDelete           PermissionKey = "APPOINTMENT_CARD.DELETE"
	PermAppointmentCardUpdate           PermissionKey = "APPOINTMENT_CARD.UPDATE"
	PermAppointmentCardView             PermissionKey = "APPOINTMENT_CARD.VIEW"
	PermAppointmentCardGenerateDocument PermissionKey = "APPOINTMENT_CARD.GENERATE_DOCUMENT"
)

// Care Coordination Permissions
const (
	PermCareCoordinationView PermissionKey = "CARE_COORDINATION.VIEW"
)

// Client Permissions
const (
	PermClientCreate       PermissionKey = "CLIENT.CREATE"
	PermClientDelete       PermissionKey = "CLIENT.DELETE"
	PermClientUpdate       PermissionKey = "CLIENT.UPDATE"
	PermClientView         PermissionKey = "CLIENT.VIEW"
	PermClientStatusUpdate PermissionKey = "CLIENT.STATUS.UPDATE"

	PermClientCarePlanCreate PermissionKey = "CLIENT.CARE_PLAN.CREATE"
	PermClientCarePlanDelete PermissionKey = "CLIENT.CARE_PLAN.DELETE"
	PermClientCarePlanUpdate PermissionKey = "CLIENT.CARE_PLAN.UPDATE"
	PermClientCarePlanView   PermissionKey = "CLIENT.CARE_PLAN.VIEW"

	PermClientDocumentsView   PermissionKey = "CLIENT.DOCUMENTS.VIEW"
	PermClientDocumentsUpload PermissionKey = "CLIENT.DOCUMENTS.UPLOAD"
	PermClientDocumentsDelete PermissionKey = "CLIENT.DOCUMENTS.DELETE"

	PermClientIncidentCreate  PermissionKey = "CLIENT.INCIDENT.CREATE"
	PermClientIncidentDelete  PermissionKey = "CLIENT.INCIDENT.DELETE"
	PermClientIncidentUpdate  PermissionKey = "CLIENT.INCIDENT.UPDATE"
	PermClientIncidentView    PermissionKey = "CLIENT.INCIDENT.VIEW"
	PermClientIncidentConfirm PermissionKey = "CLIENT.INCIDENT.CONFIRM"

	PermClientDiagnosisCreate PermissionKey = "CLIENT.DIAGNOSIS.CREATE"
	PermClientDiagnosisDelete PermissionKey = "CLIENT.DIAGNOSIS.DELETE"
	PermClientDiagnosisUpdate PermissionKey = "CLIENT.DIAGNOSIS.UPDATE"
	PermClientDiagnosisView   PermissionKey = "CLIENT.DIAGNOSIS.VIEW"

	PermClientMedicationCreate PermissionKey = "CLIENT.MEDICATION.CREATE"
	PermClientMedicationDelete PermissionKey = "CLIENT.MEDICATION.DELETE"
	PermClientMedicationUpdate PermissionKey = "CLIENT.MEDICATION.UPDATE"
	PermClientMedicationView   PermissionKey = "CLIENT.MEDICATION.VIEW"

	PermClientEmergencyContactCreate PermissionKey = "CLIENT.EMERGENCY_CONTACT.CREATE"
	PermClientEmergencyContactDelete PermissionKey = "CLIENT.EMERGENCY_CONTACT.DELETE"
	PermClientEmergencyContactUpdate PermissionKey = "CLIENT.EMERGENCY_CONTACT.UPDATE"
	PermClientEmergencyContactView   PermissionKey = "CLIENT.EMERGENCY_CONTACT.VIEW"

	PermClientInvolvedEmployeeCreate PermissionKey = "CLIENT.INVOLVED_EMPLOYEE.CREATE"
	PermClientInvolvedEmployeeDelete PermissionKey = "CLIENT.INVOLVED_EMPLOYEE.DELETE"
	PermClientInvolvedEmployeeUpdate PermissionKey = "CLIENT.INVOLVED_EMPLOYEE.UPDATE"
	PermClientInvolvedEmployeeView   PermissionKey = "CLIENT.INVOLVED_EMPLOYEE.VIEW"

	PermClientProgressReportCreate PermissionKey = "CLIENT.PROGRESS_REPORT.CREATE"
	PermClientProgressReportDelete PermissionKey = "CLIENT.PROGRESS_REPORT.DELETE"
	PermClientProgressReportUpdate PermissionKey = "CLIENT.PROGRESS_REPORT.UPDATE"
	PermClientProgressReportView   PermissionKey = "CLIENT.PROGRESS_REPORT.VIEW"

	PermClientAIProgressReportGenerate PermissionKey = "CLIENT.AI_PROGRESS_REPORT.GENERATE"
	PermClientAIProgressReportConfirm  PermissionKey = "CLIENT.AI_PROGRESS_REPORT.CONFIRM"
	PermClientAIProgressReportView     PermissionKey = "CLIENT.AI_PROGRESS_REPORT.VIEW"
)

// Contract Permissions
const (
	PermContractCreate PermissionKey = "CONTRACT.CREATE"
	PermContractDelete PermissionKey = "CONTRACT.DELETE"
	PermContractUpdate PermissionKey = "CONTRACT.UPDATE"
	PermContractView   PermissionKey = "CONTRACT.VIEW"

	PermContractTypeCreate PermissionKey = "CONTRACT_TYPE.CREATE"
	PermContractTypeDelete PermissionKey = "CONTRACT_TYPE.DELETE"
	PermContractTypeView   PermissionKey = "CONTRACT_TYPE.VIEW"
)

// Dashboard Permissions
const (
	PermDashboardView PermissionKey = "DASHBOARD.VIEW"
)

// Employee Permissions
const (
	PermEmployeeCreate PermissionKey = "EMPLOYEE.CREATE"
	PermEmployeeDelete PermissionKey = "EMPLOYEE.DELETE"
	PermEmployeeUpdate PermissionKey = "EMPLOYEE.UPDATE"
	PermEmployeeView   PermissionKey = "EMPLOYEE.VIEW"

	PermEmployeeWorkingHoursView PermissionKey = "EMPLOYEE.WORKING_HOURS.VIEW"
	PermEmployeeContractView     PermissionKey = "EMPLOYEE.CONTRACT.VIEW"
	PermEmployeeContractUpdate   PermissionKey = "EMPLOYEE.CONTRACT.UPDATE"
)

// Evaluation Permissions
const (
	PermEvaluationCreate PermissionKey = "EVALUATION.CREATE"
	PermEvaluationDelete PermissionKey = "EVALUATION.DELETE"
	PermEvaluationView   PermissionKey = "EVALUATION.VIEW"
)

// Finance & Invoice Permissions
const (
	PermFinanceView PermissionKey = "FINANCE.VIEW"

	PermInvoiceCreate PermissionKey = "INVOICE.CREATE"
	PermInvoiceDelete PermissionKey = "INVOICE.DELETE"
	PermInvoiceUpdate PermissionKey = "INVOICE.UPDATE"
	PermInvoiceView   PermissionKey = "INVOICE.VIEW"

	PermInvoicePaymentCreate PermissionKey = "INVOICE.PAYMENT.CREATE"
	PermInvoicePaymentDelete PermissionKey = "INVOICE.PAYMENT.DELETE"
	PermInvoicePaymentUpdate PermissionKey = "INVOICE.PAYMENT.UPDATE"
	PermInvoicePaymentView   PermissionKey = "INVOICE.PAYMENT.VIEW"
)

// Handbook Permissions
const (
	PermHandbookSelfView   PermissionKey = "HANDBOOK.SELF.VIEW"
	PermHandbookSelfUpdate PermissionKey = "HANDBOOK.SELF.UPDATE"

	PermHandbookDepartmentView   PermissionKey = "HANDBOOK.DEPARTMENT.VIEW"
	PermHandbookDepartmentCreate PermissionKey = "HANDBOOK.DEPARTMENT.CREATE"

	PermHandbookTemplateView    PermissionKey = "HANDBOOK.TEMPLATE.VIEW"
	PermHandbookTemplateCreate  PermissionKey = "HANDBOOK.TEMPLATE.CREATE"
	PermHandbookTemplateUpdate  PermissionKey = "HANDBOOK.TEMPLATE.UPDATE"
	PermHandbookTemplatePublish PermissionKey = "HANDBOOK.TEMPLATE.PUBLISH"

	PermHandbookStepView   PermissionKey = "HANDBOOK.STEP.VIEW"
	PermHandbookStepCreate PermissionKey = "HANDBOOK.STEP.CREATE"
	PermHandbookStepUpdate PermissionKey = "HANDBOOK.STEP.UPDATE"
	PermHandbookStepDelete PermissionKey = "HANDBOOK.STEP.DELETE"

	PermHandbookAssign            PermissionKey = "HANDBOOK.ASSIGN"
	PermHandbookEligibleEmployees PermissionKey = "HANDBOOK.ELIGIBLE_EMPLOYEES.VIEW_ALL"
)

// Incident Permissions
const (
	PermIncidentView PermissionKey = "INCIDENT.VIEW"
)

// Late Arrival Permissions
const (
	PermLateArrivalCreate    PermissionKey = "LATE_ARRIVAL.CREATE"
	PermLateArrivalCreateAll PermissionKey = "LATE_ARRIVAL.CREATE_ALL"
	PermLateArrivalView      PermissionKey = "LATE_ARRIVAL.VIEW"
	PermLateArrivalViewAll   PermissionKey = "LATE_ARRIVAL.VIEW_ALL"
)

// Leave Permissions
const (
	PermLeaveRequestCreate    PermissionKey = "LEAVE.REQUEST.CREATE"
	PermLeaveRequestUpdate    PermissionKey = "LEAVE.REQUEST.UPDATE"
	PermLeaveRequestUpdateAll PermissionKey = "LEAVE.REQUEST.UPDATE_ALL"
	PermLeaveRequestDecide    PermissionKey = "LEAVE.REQUEST.DECIDE"
	PermLeaveRequestView      PermissionKey = "LEAVE.REQUEST.VIEW"
	PermLeaveRequestViewAll   PermissionKey = "LEAVE.REQUEST.VIEW_ALL"

	PermLeaveBalanceView    PermissionKey = "LEAVE.BALANCE.VIEW"
	PermLeaveBalanceViewAll PermissionKey = "LEAVE.BALANCE.VIEW_ALL"
	PermLeaveBalanceAdjust  PermissionKey = "LEAVE.BALANCE.ADJUST"
)

// Location Permissions
const (
	PermLocationCreate PermissionKey = "LOCATION.CREATE"
	PermLocationDelete PermissionKey = "LOCATION.DELETE"
	PermLocationUpdate PermissionKey = "LOCATION.UPDATE"
	PermLocationView   PermissionKey = "LOCATION.VIEW"
)

// Organisation Permissions
const (
	PermOrganisationCreate PermissionKey = "ORGANISATION.CREATE"
	PermOrganisationDelete PermissionKey = "ORGANISATION.DELETE"
	PermOrganisationUpdate PermissionKey = "ORGANISATION.UPDATE"
	PermOrganisationView   PermissionKey = "ORGANISATION.VIEW"
)

// Permissions Management
const (
	PermPermissionsCreate PermissionKey = "PERMISSIONS.CREATE"
	PermPermissionsDelete PermissionKey = "PERMISSIONS.DELETE"
	PermPermissionsUpdate PermissionKey = "PERMISSIONS.UPDATE"
	PermPermissionsView   PermissionKey = "PERMISSIONS.VIEW"
	PermPermissionsGrant  PermissionKey = "PERMISSIONS.GRANT"
)

// Profile Permissions
const (
	PermProfileView PermissionKey = "PROFILE.VIEW"
)

// Registration Form Permissions
const (
	PermRegistrationFormDelete PermissionKey = "REGISTRATION_FORM.DELETE"
	PermRegistrationFormUpdate PermissionKey = "REGISTRATION_FORM.UPDATE"
	PermRegistrationFormView   PermissionKey = "REGISTRATION_FORM.VIEW"
)

// Intake form Permissions
const (
	PermIntakeFormDelete PermissionKey = "INTAKE_FORM.DELETE"
	PermIntakeFormUpdate PermissionKey = "INTAKE_FORM.UPDATE"
	PermIntakeFormView   PermissionKey = "INTAKE_FORM.VIEW"
	PermIntakeFormCreate PermissionKey = "INTAKE_FORM.CREATE"
)

// Roles Permissions
const (
	PermRolesCreate PermissionKey = "ROLES.CREATE"
	PermRolesDelete PermissionKey = "ROLES.DELETE"
	PermRolesUpdate PermissionKey = "ROLES.UPDATE"
	PermRolesView   PermissionKey = "ROLES.VIEW"
	PermRolesAssign PermissionKey = "ROLES.ASSIGN"
)

// Schedule & Shift Permissions
const (
	PermScheduleCreate PermissionKey = "SCHEDULE.CREATE"
	PermScheduleDelete PermissionKey = "SCHEDULE.DELETE"
	PermScheduleUpdate PermissionKey = "SCHEDULE.UPDATE"
	PermScheduleView   PermissionKey = "SCHEDULE.VIEW"

	PermScheduleSwapRequest PermissionKey = "SCHEDULE_SWAP.REQUEST"
	PermScheduleSwapRespond PermissionKey = "SCHEDULE_SWAP.RESPOND"
	PermScheduleSwapApprove PermissionKey = "SCHEDULE_SWAP.APPROVE"
	PermScheduleSwapView    PermissionKey = "SCHEDULE_SWAP.VIEW"

	PermShiftCreate PermissionKey = "SHIFT.CREATE"
	PermShiftDelete PermissionKey = "SHIFT.DELETE"
	PermShiftUpdate PermissionKey = "SHIFT.UPDATE"
	PermShiftView   PermissionKey = "SHIFT.VIEW"
)

// Sender Permissions
const (
	PermSenderCreate PermissionKey = "SENDER.CREATE"
	PermSenderDelete PermissionKey = "SENDER.DELETE"
	PermSenderUpdate PermissionKey = "SENDER.UPDATE"
	PermSenderView   PermissionKey = "SENDER.VIEW"
)

// Settings Permissions
const (
	PermSettingsView PermissionKey = "SETTINGS.VIEW"

	PermSettingsDepartmentView   PermissionKey = "SETTINGS.DEPARTMENT.VIEW"
	PermSettingsDepartmentCreate PermissionKey = "SETTINGS.DEPARTMENT.CREATE"
	PermSettingsDepartmentUpdate PermissionKey = "SETTINGS.DEPARTMENT.UPDATE"

	PermSettingsOrgProfileView   PermissionKey = "SETTINGS.ORGANIZATION_PROFILE.VIEW"
	PermSettingsOrgProfileUpdate PermissionKey = "SETTINGS.ORGANIZATION_PROFILE.UPDATE"
)

// Audit & Reports Permissions
const (
	PermAuditLogView PermissionKey = "AUDIT.LOG.VIEW"
	PermReportsView  PermissionKey = "REPORTS.VIEW"
)

// PermissionDefinition defines metadata for a system permission
type PermissionDefinition struct {
	Key         PermissionKey
	GroupKey    string
	SectionKey  string
	DisplayName string
	Description string
}

// RoleSeedDefinition defines a initial role seed with its permissions
type RoleSeedDefinition struct {
	Name        string
	Description string
	Permissions []PermissionKey
}

// AllPermissionKeys is the complete registry of system permission keys
var AllPermissionKeys = []PermissionKey{
	PermAppointmentCreate,
	PermAppointmentDelete,
	PermAppointmentUpdate,
	PermAppointmentView,
	PermAppointmentViewAll,
	PermAppointmentWorkApprovalUpdate,

	PermAppointmentCardDelete,
	PermAppointmentCardUpdate,
	PermAppointmentCardView,
	PermAppointmentCardGenerateDocument,

	PermCareCoordinationView,

	PermClientCreate,
	PermClientDelete,
	PermClientUpdate,
	PermClientView,
	PermClientStatusUpdate,

	PermClientCarePlanCreate,
	PermClientCarePlanDelete,
	PermClientCarePlanUpdate,
	PermClientCarePlanView,

	PermClientDocumentsView,
	PermClientDocumentsUpload,
	PermClientDocumentsDelete,

	PermClientIncidentCreate,
	PermClientIncidentDelete,
	PermClientIncidentUpdate,
	PermClientIncidentView,
	PermClientIncidentConfirm,

	PermClientDiagnosisCreate,
	PermClientDiagnosisDelete,
	PermClientDiagnosisUpdate,
	PermClientDiagnosisView,

	PermClientMedicationCreate,
	PermClientMedicationDelete,
	PermClientMedicationUpdate,
	PermClientMedicationView,

	PermClientEmergencyContactCreate,
	PermClientEmergencyContactDelete,
	PermClientEmergencyContactUpdate,
	PermClientEmergencyContactView,

	PermClientInvolvedEmployeeCreate,
	PermClientInvolvedEmployeeDelete,
	PermClientInvolvedEmployeeUpdate,
	PermClientInvolvedEmployeeView,

	PermClientProgressReportCreate,
	PermClientProgressReportDelete,
	PermClientProgressReportUpdate,
	PermClientProgressReportView,

	PermClientAIProgressReportGenerate,
	PermClientAIProgressReportConfirm,
	PermClientAIProgressReportView,

	PermContractCreate,
	PermContractDelete,
	PermContractUpdate,
	PermContractView,

	PermContractTypeCreate,
	PermContractTypeDelete,
	PermContractTypeView,

	PermDashboardView,

	PermEmployeeCreate,
	PermEmployeeDelete,
	PermEmployeeUpdate,
	PermEmployeeView,

	PermEmployeeWorkingHoursView,
	PermEmployeeContractView,
	PermEmployeeContractUpdate,

	PermEvaluationCreate,
	PermEvaluationDelete,
	PermEvaluationView,

	PermFinanceView,

	PermInvoiceCreate,
	PermInvoiceDelete,
	PermInvoiceUpdate,
	PermInvoiceView,

	PermInvoicePaymentCreate,
	PermInvoicePaymentDelete,
	PermInvoicePaymentUpdate,
	PermInvoicePaymentView,

	PermHandbookSelfView,
	PermHandbookSelfUpdate,

	PermHandbookDepartmentView,
	PermHandbookDepartmentCreate,

	PermHandbookTemplateView,
	PermHandbookTemplateCreate,
	PermHandbookTemplateUpdate,
	PermHandbookTemplatePublish,

	PermHandbookStepView,
	PermHandbookStepCreate,
	PermHandbookStepUpdate,
	PermHandbookStepDelete,

	PermHandbookAssign,
	PermHandbookEligibleEmployees,

	PermIncidentView,

	PermLateArrivalCreate,
	PermLateArrivalCreateAll,
	PermLateArrivalView,
	PermLateArrivalViewAll,

	PermLeaveRequestCreate,
	PermLeaveRequestUpdate,
	PermLeaveRequestUpdateAll,
	PermLeaveRequestDecide,
	PermLeaveRequestView,
	PermLeaveRequestViewAll,

	PermLeaveBalanceView,
	PermLeaveBalanceViewAll,
	PermLeaveBalanceAdjust,

	PermLocationCreate,
	PermLocationDelete,
	PermLocationUpdate,
	PermLocationView,

	PermOrganisationCreate,
	PermOrganisationDelete,
	PermOrganisationUpdate,
	PermOrganisationView,

	PermPermissionsCreate,
	PermPermissionsDelete,
	PermPermissionsUpdate,
	PermPermissionsView,
	PermPermissionsGrant,

	PermProfileView,

	PermRegistrationFormDelete,
	PermRegistrationFormUpdate,
	PermRegistrationFormView,

	PermIntakeFormCreate,
	PermIntakeFormDelete,
	PermIntakeFormUpdate,
	PermIntakeFormView,

	PermRolesCreate,
	PermRolesDelete,
	PermRolesUpdate,
	PermRolesView,
	PermRolesAssign,

	PermScheduleCreate,
	PermScheduleDelete,
	PermScheduleUpdate,
	PermScheduleView,

	PermScheduleSwapRequest,
	PermScheduleSwapRespond,
	PermScheduleSwapApprove,
	PermScheduleSwapView,

	PermShiftCreate,
	PermShiftDelete,
	PermShiftUpdate,
	PermShiftView,

	PermSenderCreate,
	PermSenderDelete,
	PermSenderUpdate,
	PermSenderView,

	PermSettingsView,

	PermSettingsDepartmentView,
	PermSettingsDepartmentCreate,
	PermSettingsDepartmentUpdate,

	PermSettingsOrgProfileView,
	PermSettingsOrgProfileUpdate,

	PermAuditLogView,
	PermReportsView,
}

// AllPermissionDefinitions returns all permission definitions with derived metadata
func AllPermissionDefinitions() []PermissionDefinition {
	defs := make([]PermissionDefinition, len(AllPermissionKeys))
	for i, key := range AllPermissionKeys {
		defs[i] = GetPermissionDefinition(key)
	}
	return defs
}

// GetPermissionDefinition derives metadata for a permission key
func GetPermissionDefinition(key PermissionKey) PermissionDefinition {
	s := string(key)
	parts := splitKeyParts(s)

	groupKey := "general"
	if len(parts) > 0 {
		groupKey = strings.ToLower(parts[0])
	}

	sectionKey := "general"
	if len(parts) > 2 {
		sectionKey = strings.ToLower(strings.Join(parts[1:len(parts)-1], "_"))
	}

	displayName := deriveDisplayName(parts)

	return PermissionDefinition{
		Key:         key,
		GroupKey:    groupKey,
		SectionKey:  sectionKey,
		DisplayName: displayName,
	}
}

func splitKeyParts(name string) []string {
	rawParts := strings.Split(strings.TrimSpace(name), ".")
	parts := make([]string, 0, len(rawParts))
	for _, part := range rawParts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	return parts
}

func deriveDisplayName(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	if len(parts) == 1 {
		return humanizeWord(parts[0])
	}
	action := humanizeWord(parts[len(parts)-1])
	contextParts := parts[:len(parts)-1]
	if len(contextParts) > 1 {
		contextParts = contextParts[1:]
	}
	ctxStr := humanizeWord(strings.Join(contextParts, " "))
	if ctxStr == "" {
		return action
	}
	return action + " " + ctxStr
}

func humanizeWord(val string) string {
	val = strings.ReplaceAll(val, "_", " ")
	words := strings.Fields(strings.ToLower(val))
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}

// DefaultRoleSeeds returns the default system roles to seed
func DefaultRoleSeeds() []RoleSeedDefinition {
	return []RoleSeedDefinition{
		{
			Name:        "admin",
			Description: "Full administrative access",
			Permissions: AllPermissionKeys,
		},
		{
			Name:        "coordinator",
			Description: "Coordinator role for operational management",
			Permissions: []PermissionKey{
				PermAppointmentCreate, PermAppointmentDelete, PermAppointmentUpdate, PermAppointmentView,
				PermAppointmentCardDelete, PermAppointmentCardUpdate, PermAppointmentCardView, PermAppointmentCardGenerateDocument,
				PermClientView, PermClientUpdate, PermClientStatusUpdate,
				PermClientCarePlanCreate, PermClientCarePlanDelete, PermClientCarePlanUpdate, PermClientCarePlanView,
				PermClientIncidentCreate, PermClientIncidentDelete, PermClientIncidentUpdate, PermClientIncidentView,
				PermClientDiagnosisCreate, PermClientDiagnosisDelete, PermClientDiagnosisUpdate, PermClientDiagnosisView,
				PermClientMedicationCreate, PermClientMedicationDelete, PermClientMedicationUpdate, PermClientMedicationView,
				PermClientEmergencyContactCreate, PermClientEmergencyContactDelete, PermClientEmergencyContactUpdate, PermClientEmergencyContactView,
				PermClientInvolvedEmployeeView,
				PermClientProgressReportCreate, PermClientProgressReportDelete, PermClientProgressReportUpdate, PermClientProgressReportView,
				PermClientAIProgressReportGenerate, PermClientAIProgressReportConfirm, PermClientAIProgressReportView,
				PermClientDocumentsView, PermClientDocumentsUpload, PermClientDocumentsDelete,
				PermHandbookSelfView, PermHandbookSelfUpdate, PermHandbookDepartmentView, PermHandbookDepartmentCreate,
				PermHandbookTemplateView, PermHandbookTemplateCreate, PermHandbookTemplateUpdate, PermHandbookTemplatePublish,
				PermHandbookStepView, PermHandbookStepCreate, PermHandbookStepUpdate, PermHandbookStepDelete, PermHandbookAssign,
				PermScheduleSwapRequest, PermScheduleSwapRespond, PermScheduleSwapView, PermScheduleSwapApprove,
				PermLeaveRequestCreate, PermLeaveRequestUpdate, PermLeaveRequestUpdateAll, PermLeaveRequestDecide, PermLeaveRequestView, PermLeaveRequestViewAll,
				PermLeaveBalanceView, PermLeaveBalanceViewAll, PermLeaveBalanceAdjust,
				PermLateArrivalCreate, PermLateArrivalCreateAll, PermLateArrivalView, PermLateArrivalViewAll,
				PermSenderCreate, PermShiftView, PermCareCoordinationView,
				PermEvaluationCreate, PermEvaluationDelete, PermEvaluationView,
			},
		},
	}
}
