package service

type ServiceHolder struct {
	PprComputeService *PprComputeService
}

var (
	Holder ServiceHolder
)
