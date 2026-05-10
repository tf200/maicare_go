package db

func NullCareTypeFromPtr(ptr *string) *CareTypeEnum {
	if ptr == nil {
		return nil
	}
	value := CareTypeEnum(*ptr)
	return &value
}

func CareTypePtrFromEnum(value *CareTypeEnum) *string {
	if value == nil {
		return nil
	}
	str := string(*value)
	return &str
}

func NullClientDocumentLabelFromPtr(ptr *string) *ClientDocumentLabelEnum {
	if ptr == nil {
		return nil
	}
	value := ClientDocumentLabelEnum(*ptr)
	return &value
}

func ClientDocumentLabelPtrFromEnum(value *ClientDocumentLabelEnum) *string {
	if value == nil {
		return nil
	}
	str := string(*value)
	return &str
}

func NullClientEducationLevelFromPtr(ptr *string) *EducationLevelEnum {
	if ptr == nil {
		return nil
	}
	value := EducationLevelEnum(*ptr)
	return &value
}

func ClientEducationLevelPtrFromEnum(value *EducationLevelEnum) *string {
	if value == nil {
		return nil
	}
	str := string(*value)
	return &str
}

func NullClientLivingSituationFromPtr(ptr *string) *ClientLivingSituationEnum {
	if ptr == nil {
		return nil
	}
	value := ClientLivingSituationEnum(*ptr)
	return &value
}

func ClientLivingSituationPtrFromEnum(value *ClientLivingSituationEnum) *string {
	if value == nil {
		return nil
	}
	str := string(*value)
	return &str
}

func NullClientLocationTransferStatusFromPtr(ptr *string) *ClientLocationTransferStatusEnum {
	if ptr == nil {
		return nil
	}
	value := ClientLocationTransferStatusEnum(*ptr)
	return &value
}

func ClientLocationTransferStatusPtrFromEnum(value *ClientLocationTransferStatusEnum) *string {
	if value == nil {
		return nil
	}
	str := string(*value)
	return &str
}

func NullClientStatusFromPtr(ptr *string) *ClientStatusEnum {
	if ptr == nil {
		return nil
	}
	value := ClientStatusEnum(*ptr)
	return &value
}

func ClientStatusPtrFromEnum(value *ClientStatusEnum) *string {
	if value == nil {
		return nil
	}
	str := string(*value)
	return &str
}

func NullIntakeCareTypeFromPtr(ptr *string) *IntakeCareTypeEnum {
	if ptr == nil {
		return nil
	}
	value := IntakeCareTypeEnum(*ptr)
	return &value
}

func IntakeCareTypePtrFromEnum(value *IntakeCareTypeEnum) *string {
	if value == nil {
		return nil
	}
	str := string(*value)
	return &str
}

func NullIntakeConclusionFromPtr(ptr *string) *IntakeConclusionEnum {
	if ptr == nil {
		return nil
	}
	value := IntakeConclusionEnum(*ptr)
	return &value
}

func IntakeConclusionPtrFromEnum(value *IntakeConclusionEnum) *string {
	if value == nil {
		return nil
	}
	str := string(*value)
	return &str
}

func NullContractAuditOperationFromPtr(ptr *string) *ContractAuditOperationEnum {
	if ptr == nil {
		return nil
	}
	value := ContractAuditOperationEnum(*ptr)
	return &value
}

func ContractAuditOperationPtrFromEnum(value *ContractAuditOperationEnum) *string {
	if value == nil {
		return nil
	}
	str := string(*value)
	return &str
}

func NullContractReminderTypeFromPtr(ptr *string) *ContractReminderTypeEnum {
	if ptr == nil {
		return nil
	}
	value := ContractReminderTypeEnum(*ptr)
	return &value
}

func ContractReminderTypePtrFromEnum(value *ContractReminderTypeEnum) *string {
	if value == nil {
		return nil
	}
	str := string(*value)
	return &str
}

func NullContractStatusFromPtr(ptr *string) *ContractStatusEnum {
	if ptr == nil {
		return nil
	}
	value := ContractStatusEnum(*ptr)
	return &value
}

func ContractStatusPtrFromEnum(value *ContractStatusEnum) *string {
	if value == nil {
		return nil
	}
	str := string(*value)
	return &str
}

func NullEmotionalStateFromPtr(ptr *string) *EmotionalStateEnum {
	if ptr == nil {
		return nil
	}
	value := EmotionalStateEnum(*ptr)
	return &value
}

func EmotionalStatePtrFromEnum(value *EmotionalStateEnum) *string {
	if value == nil {
		return nil
	}
	str := string(*value)
	return &str
}

func NullEmployeeContractTypeFromPtr(ptr *string) *EmployeeContractTypeEnum {
	if ptr == nil {
		return nil
	}
	value := EmployeeContractTypeEnum(*ptr)
	return &value
}

func EmployeeContractTypePtrFromEnum(value *EmployeeContractTypeEnum) *string {
	if value == nil {
		return nil
	}
	str := string(*value)
	return &str
}

func NullGenderFromPtr(ptr *string) *GenderEnum {
	if ptr == nil {
		return nil
	}
	value := GenderEnum(*ptr)
	return &value
}

func EmployeeGenderPtrFromEnum(value *GenderEnum) *string {
	if value == nil {
		return nil
	}
	str := string(*value)
	return &str
}

func NullFinancingActFromPtr(ptr *string) *FinancingActEnum {
	if ptr == nil {
		return nil
	}
	value := FinancingActEnum(*ptr)
	return &value
}

func FinancingActPtrFromEnum(value *FinancingActEnum) *string {
	if value == nil {
		return nil
	}
	str := string(*value)
	return &str
}

func NullFinancingOptionFromPtr(ptr *string) *FinancingOptionEnum {
	if ptr == nil {
		return nil
	}
	value := FinancingOptionEnum(*ptr)
	return &value
}

func FinancingOptionPtrFromEnum(value *FinancingOptionEnum) *string {
	if value == nil {
		return nil
	}
	str := string(*value)
	return &str
}

func NullFormStatusFromPtr(ptr *string) *FormStatusEnum {
	if ptr == nil {
		return nil
	}
	value := FormStatusEnum(*ptr)
	return &value
}

func FormStatusPtrFromEnum(value *FormStatusEnum) *string {
	if value == nil {
		return nil
	}
	str := string(*value)
	return &str
}

func NullHoursTypeFromPtr(ptr *string) *HoursTypeEnum {
	if ptr == nil {
		return nil
	}
	value := HoursTypeEnum(*ptr)
	return &value
}

func HoursTypePtrFromEnum(value *HoursTypeEnum) *string {
	if value == nil {
		return nil
	}
	str := string(*value)
	return &str
}

func NullIncidentReporterInvolvementFromPtr(ptr *string) *IncidentReporterInvolvementEnum {
	if ptr == nil {
		return nil
	}
	value := IncidentReporterInvolvementEnum(*ptr)
	return &value
}

func IncidentReporterInvolvementPtrFromEnum(value *IncidentReporterInvolvementEnum) *string {
	if value == nil {
		return nil
	}
	str := string(*value)
	return &str
}

func NullInvoiceAuditOperationFromPtr(ptr *string) *InvoiceAuditOperationEnum {
	if ptr == nil {
		return nil
	}
	value := InvoiceAuditOperationEnum(*ptr)
	return &value
}

func InvoiceAuditOperationPtrFromEnum(value *InvoiceAuditOperationEnum) *string {
	if value == nil {
		return nil
	}
	str := string(*value)
	return &str
}

func NullInvoiceStatusFromPtr(ptr *string) *InvoiceStatusEnum {
	if ptr == nil {
		return nil
	}
	value := InvoiceStatusEnum(*ptr)
	return &value
}

func InvoiceStatusPtrFromEnum(value *InvoiceStatusEnum) *string {
	if value == nil {
		return nil
	}
	str := string(*value)
	return &str
}

func NullInvoiceTypeFromPtr(ptr *string) *InvoiceTypeEnum {
	if ptr == nil {
		return nil
	}
	value := InvoiceTypeEnum(*ptr)
	return &value
}

func InvoiceTypePtrFromEnum(value *InvoiceTypeEnum) *string {
	if value == nil {
		return nil
	}
	str := string(*value)
	return &str
}

func NullLocationTypeFromPtr(ptr *string) *LocationTypeEnum {
	if ptr == nil {
		return nil
	}
	value := LocationTypeEnum(*ptr)
	return &value
}

func LocationTypePtrFromEnum(value *LocationTypeEnum) *string {
	if value == nil {
		return nil
	}
	str := string(*value)
	return &str
}

func NullNeededConsultationFromPtr(ptr *string) *NeededConsultationEnum {
	if ptr == nil {
		return nil
	}
	value := NeededConsultationEnum(*ptr)
	return &value
}

func NeededConsultationPtrFromEnum(value *NeededConsultationEnum) *string {
	if value == nil {
		return nil
	}
	str := string(*value)
	return &str
}

func NullNotificationTypeFromPtr(ptr *string) *NotificationTypeEnum {
	if ptr == nil {
		return nil
	}
	value := NotificationTypeEnum(*ptr)
	return &value
}

func NotificationTypePtrFromEnum(value *NotificationTypeEnum) *string {
	if value == nil {
		return nil
	}
	str := string(*value)
	return &str
}

func NullPaymentMethodFromPtr(ptr *string) *PaymentMethodEnum {
	if ptr == nil {
		return nil
	}
	value := PaymentMethodEnum(*ptr)
	return &value
}

func PaymentMethodPtrFromEnum(value *PaymentMethodEnum) *string {
	if value == nil {
		return nil
	}
	str := string(*value)
	return &str
}

func NullPaymentStatusFromPtr(ptr *string) *PaymentStatusEnum {
	if ptr == nil {
		return nil
	}
	value := PaymentStatusEnum(*ptr)
	return &value
}

func PaymentStatusPtrFromEnum(value *PaymentStatusEnum) *string {
	if value == nil {
		return nil
	}
	str := string(*value)
	return &str
}

func NullPhysicalInjuryFromPtr(ptr *string) *PhysicalInjuryEnum {
	if ptr == nil {
		return nil
	}
	value := PhysicalInjuryEnum(*ptr)
	return &value
}

func PhysicalInjuryPtrFromEnum(value *PhysicalInjuryEnum) *string {
	if value == nil {
		return nil
	}
	str := string(*value)
	return &str
}

func NullPriceTimeUnitFromPtr(ptr *string) *PriceTimeUnitEnum {
	if ptr == nil {
		return nil
	}
	value := PriceTimeUnitEnum(*ptr)
	return &value
}

func PriceTimeUnitPtrFromEnum(value *PriceTimeUnitEnum) *string {
	if value == nil {
		return nil
	}
	str := string(*value)
	return &str
}

func NullProgressReportTypeFromPtr(ptr *string) *ProgressReportTypeEnum {
	if ptr == nil {
		return nil
	}
	value := ProgressReportTypeEnum(*ptr)
	return &value
}

func ProgressReportTypePtrFromEnum(value *ProgressReportTypeEnum) *string {
	if value == nil {
		return nil
	}
	str := string(*value)
	return &str
}

func NullPsychologicalDamageFromPtr(ptr *string) *PsychologicalDamageEnum {
	if ptr == nil {
		return nil
	}
	value := PsychologicalDamageEnum(*ptr)
	return &value
}

func PsychologicalDamagePtrFromEnum(value *PsychologicalDamageEnum) *string {
	if value == nil {
		return nil
	}
	str := string(*value)
	return &str
}

func NullRecurrenceRiskFromPtr(ptr *string) *RecurrenceRiskEnum {
	if ptr == nil {
		return nil
	}
	value := RecurrenceRiskEnum(*ptr)
	return &value
}

func RecurrenceRiskPtrFromEnum(value *RecurrenceRiskEnum) *string {
	if value == nil {
		return nil
	}
	str := string(*value)
	return &str
}

func NullRelationStatusFromPtr(ptr *string) *RelationStatusEnum {
	if ptr == nil {
		return nil
	}
	value := RelationStatusEnum(*ptr)
	return &value
}

func RelationStatusPtrFromEnum(value *RelationStatusEnum) *string {
	if value == nil {
		return nil
	}
	str := string(*value)
	return &str
}

func NullSenderTypesFromPtr(ptr *string) *SenderTypesEnum {
	if ptr == nil {
		return nil
	}
	value := SenderTypesEnum(*ptr)
	return &value
}

func SenderTypesPtrFromEnum(value *SenderTypesEnum) *string {
	if value == nil {
		return nil
	}
	str := string(*value)
	return &str
}

func NullSeverityOfIncidentFromPtr(ptr *string) *SeverityOfIncidentEnum {
	if ptr == nil {
		return nil
	}
	value := SeverityOfIncidentEnum(*ptr)
	return &value
}

func SeverityOfIncidentPtrFromEnum(value *SeverityOfIncidentEnum) *string {
	if value == nil {
		return nil
	}
	str := string(*value)
	return &str
}
