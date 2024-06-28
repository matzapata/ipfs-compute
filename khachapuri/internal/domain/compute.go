package domain

type IComputeService interface {
	Compute(cid string, payerHeader string, computeArgs string) (res *ComputeResponse, ctx *ComputeContext, err error)
	ExecuteProgram(deploymentPath string, execEnv []string, execArgs string) (*ComputeResponse, error)
}

type ComputeResponse struct {
	Data    string            `json:"data"`
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
}

type ComputeContext struct {
	EscrowTransaction string `json:"escrow_transaction"`
}