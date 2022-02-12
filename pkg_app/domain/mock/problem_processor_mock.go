package domain_mock

import (
	"github.com/kujilabo/cocotola-api/pkg_app/domain"
	"github.com/stretchr/testify/mock"
)

type ProblemQuotaProcessorMock struct {
	mock.Mock
}

func (m *ProblemQuotaProcessorMock) GetUnitForSizeQuota() domain.QuotaUnit {
	args := m.Called()
	return args.Get(0).(domain.QuotaUnit)
}
func (m *ProblemQuotaProcessorMock) GetLimitForSizeQuota() int {
	args := m.Called()
	return args.Int(0)
}
func (m *ProblemQuotaProcessorMock) GetUnitForUpdateQuota() domain.QuotaUnit {
	args := m.Called()
	return args.Get(0).(domain.QuotaUnit)
}
func (m *ProblemQuotaProcessorMock) GetLimitForUpdateQuota() int {
	args := m.Called()
	return args.Int(0)
}
