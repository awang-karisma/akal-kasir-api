package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
	"time"
)

type ReportService struct {
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) GetReports() (models.Report, error) {
	return s.repo.GetReports("", "")
}
func (s *ReportService) GetReportsRange(from string, to string) (models.Report, error) {
	return s.repo.GetReports(from, to)
}

func (s *ReportService) GetReportToday() (models.Report, error) {
	today := time.Now().Format("2006-01-02")
	return s.repo.GetReports(today, today)
}
