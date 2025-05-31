package service

type ServiceHolder struct {
	PprComputeService *EfficiencyComputeService
}

var (
	Holder ServiceHolder
)
