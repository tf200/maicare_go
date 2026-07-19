package conv

import "time"

const DateLayout = "2006-01-02"

func DateString(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format(DateLayout)
}

func DateStringPtr(value *time.Time) *string {
	if value == nil || value.IsZero() {
		return nil
	}
	formatted := value.Format(DateLayout)
	return &formatted
}

func ParseDate(value string) (time.Time, error) {
	return time.Parse(DateLayout, value)
}

func ParseDatePtr(value *string) (*time.Time, error) {
	if value == nil || *value == "" {
		return nil, nil
	}
	parsed, err := ParseDate(*value)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}
