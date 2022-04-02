package service_mock

import (
	"github.com/stretchr/testify/mock"

	"github.com/kujilabo/cocotola-api/pkg_app/service"
)

type ProblemQuotaProcessorMock struct {
	mock.Mock
}

func (m *ProblemQuotaProcessorMock) GetUnitForSizeQuota() service.QuotaUnit {
	args := m.Called()
	return args.Get(0).(service.QuotaUnit)
}
func (m *ProblemQuotaProcessorMock) GetLimitForSizeQuota() int {
	args := m.Called()
	return args.Int(0)
}
func (m *ProblemQuotaProcessorMock) GetUnitForUpdateQuota() service.QuotaUnit {
	args := m.Called()
	return args.Get(0).(service.QuotaUnit)
}
func (m *ProblemQuotaProcessorMock) GetLimitForUpdateQuota() int {
	args := m.Called()
	return args.Int(0)
}
