package service

import "wagner/app/service/sink"

type ServiceHolder struct {
	EfficiencyComputeService *EfficiencyComputeService
	SummarySinkService       *sink.SummarySinkService
}

var (
	Holder ServiceHolder
)
