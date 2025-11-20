package db

// AppointmentStatusEnum helpers
func NullAppointmentStatusFromPtr(ptr *string) NullAppointmentStatusEnum {
	if ptr != nil {
		return NullAppointmentStatusEnum{
			AppointmentStatusEnum: AppointmentStatusEnum(*ptr),
			Valid:                 true,
		}
	}
	return NullAppointmentStatusEnum{Valid: false}
}

func AppointmentStatusPtrFromEnum(enum NullAppointmentStatusEnum) *string {
	if enum.Valid {
		str := string(enum.AppointmentStatusEnum)
		return &str
	}
	return nil
}

// CarePlanInterventionFrequencyEnum helpers
func NullCarePlanInterventionFrequencyFromPtr(ptr *string) NullCarePlanInterventionFrequencyEnum {
	if ptr != nil {
		return NullCarePlanInterventionFrequencyEnum{
			CarePlanInterventionFrequencyEnum: CarePlanInterventionFrequencyEnum(*ptr),
			Valid:                             true,
		}
	}
	return NullCarePlanInterventionFrequencyEnum{Valid: false}
}

func CarePlanInterventionFrequencyPtrFromEnum(enum NullCarePlanInterventionFrequencyEnum) *string {
	if enum.Valid {
		str := string(enum.CarePlanInterventionFrequencyEnum)
		return &str
	}
	return nil
}

// CarePlanObjectiveStatusEnum helpers
func NullCarePlanObjectiveStatusFromPtr(ptr *string) NullCarePlanObjectiveStatusEnum {
	if ptr != nil {
		return NullCarePlanObjectiveStatusEnum{
			CarePlanObjectiveStatusEnum: CarePlanObjectiveStatusEnum(*ptr),
			Valid:                       true,
		}
	}
	return NullCarePlanObjectiveStatusEnum{Valid: false}
}

func CarePlanObjectiveStatusPtrFromEnum(enum NullCarePlanObjectiveStatusEnum) *string {
	if enum.Valid {
		str := string(enum.CarePlanObjectiveStatusEnum)
		return &str
	}
	return nil
}

// CarePlanReportTypeEnum helpers
func NullCarePlanReportTypeFromPtr(ptr *string) NullCarePlanReportTypeEnum {
	if ptr != nil {
		return NullCarePlanReportTypeEnum{
			CarePlanReportTypeEnum: CarePlanReportTypeEnum(*ptr),
			Valid:                  true,
		}
	}
	return NullCarePlanReportTypeEnum{Valid: false}
}

func CarePlanReportTypePtrFromEnum(enum NullCarePlanReportTypeEnum) *string {
	if enum.Valid {
		str := string(enum.CarePlanReportTypeEnum)
		return &str
	}
	return nil
}

// CarePlanRiskLevelEnum helpers
func NullCarePlanRiskLevelFromPtr(ptr *string) NullCarePlanRiskLevelEnum {
	if ptr != nil {
		return NullCarePlanRiskLevelEnum{
			CarePlanRiskLevelEnum: CarePlanRiskLevelEnum(*ptr),
			Valid:                 true,
		}
	}
	return NullCarePlanRiskLevelEnum{Valid: false}
}

func CarePlanRiskLevelPtrFromEnum(enum NullCarePlanRiskLevelEnum) *string {
	if enum.Valid {
		str := string(enum.CarePlanRiskLevelEnum)
		return &str
	}
	return nil
}

// CarePlanStatusEnum helpers
func NullCarePlanStatusFromPtr(ptr *string) NullCarePlanStatusEnum {
	if ptr != nil {
		return NullCarePlanStatusEnum{
			CarePlanStatusEnum: CarePlanStatusEnum(*ptr),
			Valid:              true,
		}
	}
	return NullCarePlanStatusEnum{Valid: false}
}

func CarePlanStatusPtrFromEnum(enum NullCarePlanStatusEnum) *string {
	if enum.Valid {
		str := string(enum.CarePlanStatusEnum)
		return &str
	}
	return nil
}

// CarePlanTimeframeEnum helpers
func NullCarePlanTimeframeFromPtr(timeframePtr *string) NullCarePlanTimeframeEnum {
	if timeframePtr != nil {
		return NullCarePlanTimeframeEnum{
			CarePlanTimeframeEnum: CarePlanTimeframeEnum(*timeframePtr),
			Valid:                 true,
		}
	}
	return NullCarePlanTimeframeEnum{Valid: false}
}

func CarePlanTimeframePtrFromEnum(timeframeEnum NullCarePlanTimeframeEnum) *string {
	if timeframeEnum.Valid {
		timeframeStr := string(timeframeEnum.CarePlanTimeframeEnum)
		return &timeframeStr
	}
	return nil
}

// CareTypeEnum helpers
func NullCareTypeFromPtr(ptr *string) NullCareTypeEnum {
	if ptr != nil {
		return NullCareTypeEnum{
			CareTypeEnum: CareTypeEnum(*ptr),
			Valid:        true,
		}
	}
	return NullCareTypeEnum{Valid: false}
}

func CareTypePtrFromEnum(enum NullCareTypeEnum) *string {
	if enum.Valid {
		str := string(enum.CareTypeEnum)
		return &str
	}
	return nil
}

// ClientDocumentLabelEnum helpers
func NullClientDocumentLabelFromPtr(ptr *string) NullClientDocumentLabelEnum {
	if ptr != nil {
		return NullClientDocumentLabelEnum{
			ClientDocumentLabelEnum: ClientDocumentLabelEnum(*ptr),
			Valid:                   true,
		}
	}
	return NullClientDocumentLabelEnum{Valid: false}
}

func ClientDocumentLabelPtrFromEnum(enum NullClientDocumentLabelEnum) *string {
	if enum.Valid {
		str := string(enum.ClientDocumentLabelEnum)
		return &str
	}
	return nil
}

// ClientEducationLevelEnum helpers
func NullClientEducationLevelFromPtr(ptr *string) NullClientEducationLevelEnum {
	if ptr != nil {
		return NullClientEducationLevelEnum{
			ClientEducationLevelEnum: ClientEducationLevelEnum(*ptr),
			Valid:                    true,
		}
	}
	return NullClientEducationLevelEnum{Valid: false}
}

func ClientEducationLevelPtrFromEnum(enum NullClientEducationLevelEnum) *string {
	if enum.Valid {
		str := string(enum.ClientEducationLevelEnum)
		return &str
	}
	return nil
}

// ClientGenderEnum helpers
func NullClientGenderFromPtr(ptr *string) NullClientGenderEnum {
	if ptr != nil {
		return NullClientGenderEnum{
			ClientGenderEnum: ClientGenderEnum(*ptr),
			Valid:            true,
		}
	}
	return NullClientGenderEnum{Valid: false}
}

func ClientGenderPtrFromEnum(enum NullClientGenderEnum) *string {
	if enum.Valid {
		str := string(enum.ClientGenderEnum)
		return &str
	}
	return nil
}

// ClientLivingSituationEnum helpers
func NullClientLivingSituationFromPtr(ptr *string) NullClientLivingSituationEnum {
	if ptr != nil {
		return NullClientLivingSituationEnum{
			ClientLivingSituationEnum: ClientLivingSituationEnum(*ptr),
			Valid:                     true,
		}
	}
	return NullClientLivingSituationEnum{Valid: false}
}

func ClientLivingSituationPtrFromEnum(enum NullClientLivingSituationEnum) *string {
	if enum.Valid {
		str := string(enum.ClientLivingSituationEnum)
		return &str
	}
	return nil
}

// ClientLocationTransferStatusEnum helpers
func NullClientLocationTransferStatusFromPtr(ptr *string) NullClientLocationTransferStatusEnum {
	if ptr != nil {
		return NullClientLocationTransferStatusEnum{
			ClientLocationTransferStatusEnum: ClientLocationTransferStatusEnum(*ptr),
			Valid:                            true,
		}
	}
	return NullClientLocationTransferStatusEnum{Valid: false}
}

func ClientLocationTransferStatusPtrFromEnum(enum NullClientLocationTransferStatusEnum) *string {
	if enum.Valid {
		str := string(enum.ClientLocationTransferStatusEnum)
		return &str
	}
	return nil
}

// ClientStatusEnum helpers
func NullClientStatusFromPtr(ptr *string) NullClientStatusEnum {
	if ptr != nil {
		return NullClientStatusEnum{
			ClientStatusEnum: ClientStatusEnum(*ptr),
			Valid:            true,
		}
	}
	return NullClientStatusEnum{Valid: false}
}

func ClientStatusPtrFromEnum(enum NullClientStatusEnum) *string {
	if enum.Valid {
		str := string(enum.ClientStatusEnum)
		return &str
	}
	return nil
}

// ContractAuditOperationEnum helpers
func NullContractAuditOperationFromPtr(ptr *string) NullContractAuditOperationEnum {
	if ptr != nil {
		return NullContractAuditOperationEnum{
			ContractAuditOperationEnum: ContractAuditOperationEnum(*ptr),
			Valid:                      true,
		}
	}
	return NullContractAuditOperationEnum{Valid: false}
}

func ContractAuditOperationPtrFromEnum(enum NullContractAuditOperationEnum) *string {
	if enum.Valid {
		str := string(enum.ContractAuditOperationEnum)
		return &str
	}
	return nil
}

// ContractReminderTypeEnum helpers
func NullContractReminderTypeFromPtr(ptr *string) NullContractReminderTypeEnum {
	if ptr != nil {
		return NullContractReminderTypeEnum{
			ContractReminderTypeEnum: ContractReminderTypeEnum(*ptr),
			Valid:                    true,
		}
	}
	return NullContractReminderTypeEnum{Valid: false}
}

func ContractReminderTypePtrFromEnum(enum NullContractReminderTypeEnum) *string {
	if enum.Valid {
		str := string(enum.ContractReminderTypeEnum)
		return &str
	}
	return nil
}

// ContractStatusEnum helpers
func NullContractStatusFromPtr(ptr *string) NullContractStatusEnum {
	if ptr != nil {
		return NullContractStatusEnum{
			ContractStatusEnum: ContractStatusEnum(*ptr),
			Valid:              true,
		}
	}
	return NullContractStatusEnum{Valid: false}
}

func ContractStatusPtrFromEnum(enum NullContractStatusEnum) *string {
	if enum.Valid {
		str := string(enum.ContractStatusEnum)
		return &str
	}
	return nil
}

// EmotionalStateEnum helpers
func NullEmotionalStateFromPtr(ptr *string) NullEmotionalStateEnum {
	if ptr != nil {
		return NullEmotionalStateEnum{
			EmotionalStateEnum: EmotionalStateEnum(*ptr),
			Valid:              true,
		}
	}
	return NullEmotionalStateEnum{Valid: false}
}

func EmotionalStatePtrFromEnum(enum NullEmotionalStateEnum) *string {
	if enum.Valid {
		str := string(enum.EmotionalStateEnum)
		return &str
	}
	return nil
}

// EmployeeContractTypeEnum helpers
func NullEmployeeContractTypeFromPtr(ptr *string) NullEmployeeContractTypeEnum {
	if ptr != nil {
		return NullEmployeeContractTypeEnum{
			EmployeeContractTypeEnum: EmployeeContractTypeEnum(*ptr),
			Valid:                    true,
		}
	}
	return NullEmployeeContractTypeEnum{Valid: false}
}

func EmployeeContractTypePtrFromEnum(enum NullEmployeeContractTypeEnum) *string {
	if enum.Valid {
		str := string(enum.EmployeeContractTypeEnum)
		return &str
	}
	return nil
}

// EmployeeGenderEnum helpers
func NullEmployeeGenderFromPtr(ptr *string) NullEmployeeGenderEnum {
	if ptr != nil {
		return NullEmployeeGenderEnum{
			EmployeeGenderEnum: EmployeeGenderEnum(*ptr),
			Valid:              true,
		}
	}
	return NullEmployeeGenderEnum{Valid: false}
}

func EmployeeGenderPtrFromEnum(enum NullEmployeeGenderEnum) *string {
	if enum.Valid {
		str := string(enum.EmployeeGenderEnum)
		return &str
	}
	return nil
}

// FinancingActEnum helpers
func NullFinancingActFromPtr(ptr *string) NullFinancingActEnum {
	if ptr != nil {
		return NullFinancingActEnum{
			FinancingActEnum: FinancingActEnum(*ptr),
			Valid:            true,
		}
	}
	return NullFinancingActEnum{Valid: false}
}

func FinancingActPtrFromEnum(enum NullFinancingActEnum) *string {
	if enum.Valid {
		str := string(enum.FinancingActEnum)
		return &str
	}
	return nil
}

// FinancingOptionEnum helpers
func NullFinancingOptionFromPtr(ptr *string) NullFinancingOptionEnum {
	if ptr != nil {
		return NullFinancingOptionEnum{
			FinancingOptionEnum: FinancingOptionEnum(*ptr),
			Valid:               true,
		}
	}
	return NullFinancingOptionEnum{Valid: false}
}

func FinancingOptionPtrFromEnum(enum NullFinancingOptionEnum) *string {
	if enum.Valid {
		str := string(enum.FinancingOptionEnum)
		return &str
	}
	return nil
}

// FormStatusEnum helpers
func NullFormStatusFromPtr(ptr *string) NullFormStatusEnum {
	if ptr != nil {
		return NullFormStatusEnum{
			FormStatusEnum: FormStatusEnum(*ptr),
			Valid:          true,
		}
	}
	return NullFormStatusEnum{Valid: false}
}

func FormStatusPtrFromEnum(enum NullFormStatusEnum) *string {
	if enum.Valid {
		str := string(enum.FormStatusEnum)
		return &str
	}
	return nil
}

// HoursTypeEnum helpers
func NullHoursTypeFromPtr(ptr *string) NullHoursTypeEnum {
	if ptr != nil {
		return NullHoursTypeEnum{
			HoursTypeEnum: HoursTypeEnum(*ptr),
			Valid:         true,
		}
	}
	return NullHoursTypeEnum{Valid: false}
}

func HoursTypePtrFromEnum(enum NullHoursTypeEnum) *string {
	if enum.Valid {
		str := string(enum.HoursTypeEnum)
		return &str
	}
	return nil
}

// IncidentReporterInvolvementEnum helpers
func NullIncidentReporterInvolvementFromPtr(ptr *string) NullIncidentReporterInvolvementEnum {
	if ptr != nil {
		return NullIncidentReporterInvolvementEnum{
			IncidentReporterInvolvementEnum: IncidentReporterInvolvementEnum(*ptr),
			Valid:                           true,
		}
	}
	return NullIncidentReporterInvolvementEnum{Valid: false}
}

func IncidentReporterInvolvementPtrFromEnum(enum NullIncidentReporterInvolvementEnum) *string {
	if enum.Valid {
		str := string(enum.IncidentReporterInvolvementEnum)
		return &str
	}
	return nil
}

// InvoiceAuditOperationEnum helpers
func NullInvoiceAuditOperationFromPtr(ptr *string) NullInvoiceAuditOperationEnum {
	if ptr != nil {
		return NullInvoiceAuditOperationEnum{
			InvoiceAuditOperationEnum: InvoiceAuditOperationEnum(*ptr),
			Valid:                     true,
		}
	}
	return NullInvoiceAuditOperationEnum{Valid: false}
}

func InvoiceAuditOperationPtrFromEnum(enum NullInvoiceAuditOperationEnum) *string {
	if enum.Valid {
		str := string(enum.InvoiceAuditOperationEnum)
		return &str
	}
	return nil
}

// InvoiceStatusEnum helpers
func NullInvoiceStatusFromPtr(ptr *string) NullInvoiceStatusEnum {
	if ptr != nil {
		return NullInvoiceStatusEnum{
			InvoiceStatusEnum: InvoiceStatusEnum(*ptr),
			Valid:             true,
		}
	}
	return NullInvoiceStatusEnum{Valid: false}
}

func InvoiceStatusPtrFromEnum(enum NullInvoiceStatusEnum) *string {
	if enum.Valid {
		str := string(enum.InvoiceStatusEnum)
		return &str
	}
	return nil
}

// InvoiceTypeEnum helpers
func NullInvoiceTypeFromPtr(ptr *string) NullInvoiceTypeEnum {
	if ptr != nil {
		return NullInvoiceTypeEnum{
			InvoiceTypeEnum: InvoiceTypeEnum(*ptr),
			Valid:           true,
		}
	}
	return NullInvoiceTypeEnum{Valid: false}
}

func InvoiceTypePtrFromEnum(enum NullInvoiceTypeEnum) *string {
	if enum.Valid {
		str := string(enum.InvoiceTypeEnum)
		return &str
	}
	return nil
}

// LocationTypeEnum helpers
func NullLocationTypeFromPtr(ptr *string) NullLocationTypeEnum {
	if ptr != nil {
		return NullLocationTypeEnum{
			LocationTypeEnum: LocationTypeEnum(*ptr),
			Valid:            true,
		}
	}
	return NullLocationTypeEnum{Valid: false}
}

func LocationTypePtrFromEnum(enum NullLocationTypeEnum) *string {
	if enum.Valid {
		str := string(enum.LocationTypeEnum)
		return &str
	}
	return nil
}

// NeededConsultationEnum helpers
func NullNeededConsultationFromPtr(ptr *string) NullNeededConsultationEnum {
	if ptr != nil {
		return NullNeededConsultationEnum{
			NeededConsultationEnum: NeededConsultationEnum(*ptr),
			Valid:                  true,
		}
	}
	return NullNeededConsultationEnum{Valid: false}
}

func NeededConsultationPtrFromEnum(enum NullNeededConsultationEnum) *string {
	if enum.Valid {
		str := string(enum.NeededConsultationEnum)
		return &str
	}
	return nil
}

// NotificationTypeEnum helpers
func NullNotificationTypeFromPtr(ptr *string) NullNotificationTypeEnum {
	if ptr != nil {
		return NullNotificationTypeEnum{
			NotificationTypeEnum: NotificationTypeEnum(*ptr),
			Valid:                true,
		}
	}
	return NullNotificationTypeEnum{Valid: false}
}

func NotificationTypePtrFromEnum(enum NullNotificationTypeEnum) *string {
	if enum.Valid {
		str := string(enum.NotificationTypeEnum)
		return &str
	}
	return nil
}

// PaymentMethodEnum helpers
func NullPaymentMethodFromPtr(ptr *string) NullPaymentMethodEnum {
	if ptr != nil {
		return NullPaymentMethodEnum{
			PaymentMethodEnum: PaymentMethodEnum(*ptr),
			Valid:             true,
		}
	}
	return NullPaymentMethodEnum{Valid: false}
}

func PaymentMethodPtrFromEnum(enum NullPaymentMethodEnum) *string {
	if enum.Valid {
		str := string(enum.PaymentMethodEnum)
		return &str
	}
	return nil
}

// PaymentStatusEnum helpers
func NullPaymentStatusFromPtr(ptr *string) NullPaymentStatusEnum {
	if ptr != nil {
		return NullPaymentStatusEnum{
			PaymentStatusEnum: PaymentStatusEnum(*ptr),
			Valid:             true,
		}
	}
	return NullPaymentStatusEnum{Valid: false}
}

func PaymentStatusPtrFromEnum(enum NullPaymentStatusEnum) *string {
	if enum.Valid {
		str := string(enum.PaymentStatusEnum)
		return &str
	}
	return nil
}

// PhysicalInjuryEnum helpers
func NullPhysicalInjuryFromPtr(ptr *string) NullPhysicalInjuryEnum {
	if ptr != nil {
		return NullPhysicalInjuryEnum{
			PhysicalInjuryEnum: PhysicalInjuryEnum(*ptr),
			Valid:              true,
		}
	}
	return NullPhysicalInjuryEnum{Valid: false}
}

func PhysicalInjuryPtrFromEnum(enum NullPhysicalInjuryEnum) *string {
	if enum.Valid {
		str := string(enum.PhysicalInjuryEnum)
		return &str
	}
	return nil
}

// PriceTimeUnitEnum helpers
func NullPriceTimeUnitFromPtr(ptr *string) NullPriceTimeUnitEnum {
	if ptr != nil {
		return NullPriceTimeUnitEnum{
			PriceTimeUnitEnum: PriceTimeUnitEnum(*ptr),
			Valid:             true,
		}
	}
	return NullPriceTimeUnitEnum{Valid: false}
}

func PriceTimeUnitPtrFromEnum(enum NullPriceTimeUnitEnum) *string {
	if enum.Valid {
		str := string(enum.PriceTimeUnitEnum)
		return &str
	}
	return nil
}

// ProgressReportTypeEnum helpers
func NullProgressReportTypeFromPtr(ptr *string) NullProgressReportTypeEnum {
	if ptr != nil {
		return NullProgressReportTypeEnum{
			ProgressReportTypeEnum: ProgressReportTypeEnum(*ptr),
			Valid:                  true,
		}
	}
	return NullProgressReportTypeEnum{Valid: false}
}

func ProgressReportTypePtrFromEnum(enum NullProgressReportTypeEnum) *string {
	if enum.Valid {
		str := string(enum.ProgressReportTypeEnum)
		return &str
	}
	return nil
}

// PsychologicalDamageEnum helpers
func NullPsychologicalDamageFromPtr(ptr *string) NullPsychologicalDamageEnum {
	if ptr != nil {
		return NullPsychologicalDamageEnum{
			PsychologicalDamageEnum: PsychologicalDamageEnum(*ptr),
			Valid:                   true,
		}
	}
	return NullPsychologicalDamageEnum{Valid: false}
}

func PsychologicalDamagePtrFromEnum(enum NullPsychologicalDamageEnum) *string {
	if enum.Valid {
		str := string(enum.PsychologicalDamageEnum)
		return &str
	}
	return nil
}

// RecurrenceRiskEnum helpers
func NullRecurrenceRiskFromPtr(ptr *string) NullRecurrenceRiskEnum {
	if ptr != nil {
		return NullRecurrenceRiskEnum{
			RecurrenceRiskEnum: RecurrenceRiskEnum(*ptr),
			Valid:              true,
		}
	}
	return NullRecurrenceRiskEnum{Valid: false}
}

func RecurrenceRiskPtrFromEnum(enum NullRecurrenceRiskEnum) *string {
	if enum.Valid {
		str := string(enum.RecurrenceRiskEnum)
		return &str
	}
	return nil
}

// RecurrenceTypeEnum helpers
func NullRecurrenceTypeFromPtr(ptr *string) NullRecurrenceTypeEnum {
	if ptr != nil {
		return NullRecurrenceTypeEnum{
			RecurrenceTypeEnum: RecurrenceTypeEnum(*ptr),
			Valid:              true,
		}
	}
	return NullRecurrenceTypeEnum{Valid: false}
}

func RecurrenceTypePtrFromEnum(enum NullRecurrenceTypeEnum) *string {
	if enum.Valid {
		str := string(enum.RecurrenceTypeEnum)
		return &str
	}
	return nil
}

// RelationStatusEnum helpers
func NullRelationStatusFromPtr(ptr *string) NullRelationStatusEnum {
	if ptr != nil {
		return NullRelationStatusEnum{
			RelationStatusEnum: RelationStatusEnum(*ptr),
			Valid:              true,
		}
	}
	return NullRelationStatusEnum{Valid: false}
}

func RelationStatusPtrFromEnum(enum NullRelationStatusEnum) *string {
	if enum.Valid {
		str := string(enum.RelationStatusEnum)
		return &str
	}
	return nil
}

// SenderTypesEnum helpers
func NullSenderTypesFromPtr(ptr *string) NullSenderTypesEnum {
	if ptr != nil {
		return NullSenderTypesEnum{
			SenderTypesEnum: SenderTypesEnum(*ptr),
			Valid:           true,
		}
	}
	return NullSenderTypesEnum{Valid: false}
}

func SenderTypesPtrFromEnum(enum NullSenderTypesEnum) *string {
	if enum.Valid {
		str := string(enum.SenderTypesEnum)
		return &str
	}
	return nil
}

// SeverityOfIncidentEnum helpers
func NullSeverityOfIncidentFromPtr(ptr *string) NullSeverityOfIncidentEnum {
	if ptr != nil {
		return NullSeverityOfIncidentEnum{
			SeverityOfIncidentEnum: SeverityOfIncidentEnum(*ptr),
			Valid:                  true,
		}
	}
	return NullSeverityOfIncidentEnum{Valid: false}
}

func SeverityOfIncidentPtrFromEnum(enum NullSeverityOfIncidentEnum) *string {
	if enum.Valid {
		str := string(enum.SeverityOfIncidentEnum)
		return &str
	}
	return nil
}
