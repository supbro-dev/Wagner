package service

type ServiceHolder struct {
	EfficiencyComputeService *EfficiencyComputeService
}

var (
	Holder ServiceHolder
)
