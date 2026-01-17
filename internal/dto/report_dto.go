package dto

import "time"

type ReportType string

const (
	REPORT_SALES    ReportType = "SALES"
	REPORT_PURCHASE ReportType = "PURCHASE"
)

type ReportByDate struct {
	StartDate  time.Time  `json:"start_date"`
	EndDate    time.Time  `json:"end_date"`
	ReportType ReportType `json:"report_type"`
}

type ReportAll struct {
	ReportType ReportType `json:"report_type"`
}
