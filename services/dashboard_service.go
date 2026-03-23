package services

import (
	"context"
	"fs-backend/models"
	"fs-backend/repository"
	"time"
)

type DashboardResponse struct {
	TotalDocuments int64           `json:"total_documents"`
	TotalHBLs      int64           `json:"total_hbls"`
	RecentActivity int64           `json:"recent_activity"`
	Data           []models.HBLDoc `json:"data"`
}

type DashboardService interface {
	GetDashboardDetails(ctx context.Context) (*DashboardResponse, error)
	DeleteDocument(ctx context.Context, id string) error
}

type dashboardService struct {
	hblDocRepo repository.HBLDocRepository
	hblRepo    repository.HBLRepository
}

func NewDashboardService(hblDocRepo repository.HBLDocRepository, hblRepo repository.HBLRepository) DashboardService {
	return &dashboardService{
		hblDocRepo: hblDocRepo,
		hblRepo:    hblRepo,
	}
}

func (s *dashboardService) GetDashboardDetails(ctx context.Context) (*DashboardResponse, error) {
	totalDocs, err := s.hblDocRepo.CountTotal(ctx)
	if err != nil {
		return nil, err
	}

	totalHBLs, err := s.hblRepo.CountTotal(ctx)
	if err != nil {
		return nil, err
	}

	allDocs, err := s.hblDocRepo.GetRecent(ctx, 0) // fetch all documents
	if err != nil {
		return nil, err
	}
	
	now := time.Now()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	
	var recentCount int64 = 0
	for _, doc := range allDocs {
		if !doc.CreatedAt.Before(startOfToday) {
			recentCount++
		}
	}

	// Ensure an empty slice is returned instead of null for JSON serialization
	if allDocs == nil {
		allDocs = []models.HBLDoc{}
	}

	return &DashboardResponse{
		TotalDocuments: totalDocs,
		TotalHBLs:      totalHBLs,
		RecentActivity: recentCount,
		Data:           allDocs,
	}, nil
}

func (s *dashboardService) DeleteDocument(ctx context.Context, id string) error {
	return s.hblDocRepo.DeleteByID(ctx, id)
}
