package models

type ReportBestSeller struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    int64  `json:"quantity"`
	TotalAmount int64  `json:"total_amount"`
}

type Report struct {
	TotalRevenue      int64            `json:"total_revenue"`
	TotalTransactions int64            `json:"total_transactions"`
	BestSeller        ReportBestSeller `json:"best_seller"`
}
