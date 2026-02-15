package repositories

import (
	"database/sql"
	"kasir-api/models"
	"strconv"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) GetReports(from string, to string) (models.Report, error) {
	var params []string
	query := `
		with transaction_range as 
		(
			select
				product_id, p.name as product_name, quantity, subtotal, price, td.created_at as created_at
			FROM transaction_details as td
			left join products as p on p.id = td.product_id
			WHERE 1=1
	`
	if len(from) != 0 {
		params = append(params, from)
		query += "AND td.created_at >= $" + strconv.Itoa(len(params))
	}
	if len(to) != 0 {
		params = append(params, to)
		query += "AND td.created_at <= $" + strconv.Itoa(len(params))
	}
	query += `)
		select product_id, product_name, sum(quantity) as quantity, sum(subtotal) as subtotal, COUNT(*) OVER() as transaction_count
		from transaction_range
		group by product_id, product_name
		order by quantity desc
	`
	var rows *sql.Rows
	var err error
	if len(from) != 0 || len(to) != 0 {
		rows, err = r.db.Query(query, from, to)
	} else {
		rows, err = r.db.Query(query)
	}
	if err != nil {
		return models.Report{}, err
	}
	defer rows.Close()

	var report models.Report

	for rows.Next() {
		var transactionDetail models.TransactionDetail
		var transactionCount int64
		err := rows.Scan(&transactionDetail.ProductID, &transactionDetail.ProductName, &transactionDetail.Quantity, &transactionDetail.Subtotal, &transactionCount)
		if err != nil {
			return models.Report{}, err
		}
		report.TotalRevenue += transactionDetail.Subtotal
		report.TotalTransactions += transactionCount
		if transactionCount > int64(report.BestSeller.Quantity) {
			report.BestSeller.ProductID = transactionDetail.ProductID
			report.BestSeller.ProductName = transactionDetail.ProductName
			report.BestSeller.Quantity = transactionCount
			report.BestSeller.TotalAmount = transactionDetail.Subtotal
		}
	}

	return report, nil
}
