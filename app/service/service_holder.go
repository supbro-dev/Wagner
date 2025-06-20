package service

import (
	"wagner/app/service/sink"
)

type ServiceHolder struct {
	EfficiencyComputeService  *EfficiencyComputeService
	EfficiencyService         *EfficiencyService
	SummarySinkService        *sink.SummarySinkService
	EmployeeStatusSinkService *sink.EmployeeStatusSinkService
}

var (
	Holder ServiceHolder
)
