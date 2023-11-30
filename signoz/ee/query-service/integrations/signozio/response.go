package signozio

type status string

const (
	statusSuccess status = "success"
	statusError   status = "error"
)

type ActivationResult struct {
	Status    status              `json:"status"`
	Data      *ActivationResponse `json:"data,omitempty"`
	ErrorType string              `json:"errorType,omitempty"`
	Error     string              `json:"error,omitempty"`
}

type ActivationResponse struct {
	ActivationId string `json:"ActivationId"`
	PlanDetails  string `json:"PlanDetails"`
}
