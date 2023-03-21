package models

import "time"

type TransactionReport struct {
	Invoice      string    `json:"invoice"`
	Member       string    `json:"transacting_member"`
	Date         time.Time `json:"date"`
	PaidDate     time.Time `json:"paid_date"`
	AddCharge    int       `json:"charge_amount"`
	Discount     int       `json:"discount"`
	Tax          int       `json:"tax"`
	TransactedAt string    `json:"transacted_at"`
	PaidStatus   string    `json:"paid_status"`
	Cashier      string    `json:"cashier"`
	LaundryType  string    `json:"laundry_type"`
	UnitPrice    int       `json:"unit_price"`
	Qty          int       `json:"qty"`
}
