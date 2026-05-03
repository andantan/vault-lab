package v2

type transferRequest struct {
	From  string `json:"from"` // sender address
	To    string `json:"to"`
	Value string `json:"value"` // wei as decimal string
}

type transferResponse struct {
	TxHash string `json:"tx_hash"`
}
