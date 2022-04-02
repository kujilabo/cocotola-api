package service

import (
	"context"
	"errors"

	"github.com/kujilabo/cocotola-api/pkg_app/domain"
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
	IsExceeded(ctx context.Context, operator domain.StudentModel, name string, unit QuotaUnit, limit int) (bool, error)
	Increment(ctx context.Context, operator domain.StudentModel, name string, unit QuotaUnit, limit int, count int) (bool, error)
}
