package domain

import (
	"context"
	"errors"
)

type QuotaUnit string

var (
	ErrQuotaExceeded           = errors.New("quota exceeded")
	UnitPersitance   QuotaUnit = "persitance"
	UnitMonth        QuotaUnit = "month"
	UnitDay          QuotaUnit = "day"
)

type UserQuotaRepository interface {
	IsExceeded(ctx context.Context, operator Student, name string, unit QuotaUnit, limit int) (bool, error)
	Increment(ctx context.Context, operator Student, name string, unit QuotaUnit, limit int, count int) (bool, error)
}
