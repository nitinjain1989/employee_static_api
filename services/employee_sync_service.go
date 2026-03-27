package services

import (
	"context"
	"static-api/dto"
	"static-api/repositories"
)

type SyncService struct {
	repo *repositories.SyncRepository
}

func NewSyncService(r *repositories.SyncRepository) *SyncService {
	return &SyncService{repo: r}
}

func (s *SyncService) Sync(ctx context.Context, req dto.SyncRequest) (dto.SyncResponse, error) {

	data, cursor, hasMore, err := s.repo.Sync(ctx, req)
	if err != nil {
		return dto.SyncResponse{}, err
	}

	return dto.SyncResponse{
		Employees:  data,
		NextCursor: cursor,
		HasMore:    hasMore,
	}, nil
}
