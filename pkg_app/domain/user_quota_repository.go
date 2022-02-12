package domain

import (
	"context"
	"errors"
)

type QuotaUnit string
type QuotaName string

var (
	ErrQuotaExceeded              = errors.New("quota exceeded")
	QuotaUnitPersitance QuotaUnit = "persitance"
	QuotaUnitMonth      QuotaUnit = "month"
	QuotaUnitDay        QuotaUnit = "day"
	QuotaNameSize       QuotaName = "Size"
	QuotaNameUpdate     QuotaName = "Update"
)

type UserQuotaRepository interface {
	IsExceeded(ctx context.Context, operator Student, name string, unit QuotaUnit, limit int) (bool, error)
	Increment(ctx context.Context, operator Student, name string, unit QuotaUnit, limit int, count int) (bool, error)
}
