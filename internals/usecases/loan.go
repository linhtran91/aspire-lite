package usecases

type Loan struct {
	Amount     float64 `json:"amount"`
	Term       int     `json:"term"`
	Date       string  `json:"date"`
	CustomerID int64   `json:"customer_id"`
}

type Repayment struct {
	Amount float64 `json:"amount"`
	Date   string  `json:"date"`
}

type SubmittedRepayment struct {
	Amount float64 `json:"amount"`
}
