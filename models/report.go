package models

type TotalSalesReport struct {
	TotalSales float64 `json:"total_sales"`
}

type PopularItemReport struct {
	ProductID string `json:"product_id"`
	Name      string `json:"name"`
	Count     int    `json:"count"`
}
