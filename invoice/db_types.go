package invoice

type PaymentMethod string

const (
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodCreditCard   PaymentMethod = "credit_card"
	PaymentMethodCheck        PaymentMethod = "check"
	PaymentMethodCash         PaymentMethod = "cash"
	PaymentMethodOther        PaymentMethod = "other"
)

type PaymentStatus string

const (
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusReversed  PaymentStatus = "reversed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

type InvoiceStatus string

const (
	InvoiceStatusOutstanding   InvoiceStatus = "outstanding"
	InvoiceStatusPartiallyPaid InvoiceStatus = "partially_paid"
	InvoiceStatusPaid          InvoiceStatus = "paid"
	InvoiceStatusExpired       InvoiceStatus = "expired"
	InvoiceStatusOverpaid      InvoiceStatus = "overpaid"
	InvoiceStatusImported      InvoiceStatus = "imported"
	InvoiceStatusConcept       InvoiceStatus = "concept"
)
