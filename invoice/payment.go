package invoice

import "fmt"

const (
	PAYMENT_TOLERANCE float64 = 50
)

func DetermineInvoiceStatus(invoiceTotal, totalPaid float64) (InvoiceStatus, error) {
	diffrence := totalPaid - invoiceTotal

	if totalPaid <= PAYMENT_TOLERANCE {
		return InvoiceStatusOutstanding, nil
	}

	if diffrence < -PAYMENT_TOLERANCE {
		return InvoiceStatusPartiallyPaid, nil
	}

	if diffrence >= -PAYMENT_TOLERANCE && diffrence <= PAYMENT_TOLERANCE {
		return InvoiceStatusPaid, nil
	}
	if diffrence > PAYMENT_TOLERANCE {
		return InvoiceStatusOverpaid, nil
	}

	return "", fmt.Errorf("could not determine invoice status for totalPaid: %f, invoiceTotal: %f", totalPaid, invoiceTotal)

}
