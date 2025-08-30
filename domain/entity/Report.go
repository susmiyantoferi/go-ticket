package entity

type ReportsSales struct {
	EventID          uint    `json:"event_id"`
	EventName        string  `json:"event_name"`
	EventDescription string  `json:"event_description"`
	Month            string  `json:"month"`
	TotalQty         int     `json:"total_qty"`
	TotalSales       float64 `json:"total_sales"`
}
