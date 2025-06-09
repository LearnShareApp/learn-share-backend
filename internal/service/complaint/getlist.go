package complaint

import (
	"context"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	"github.com/LearnShareApp/learn-share-backend/pkg/workerpool"
)

func (s *ComplaintService) GetComplaintList(ctx context.Context) ([]*entities.Complaint, error) {
	complaints, err := s.repo.GetAllComplaints(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get complaints: %w", err)
	}

	// all unique users
	users := make(map[int]*entities.User)
	for _, c := range complaints {
		users[c.ComplainerID] = nil
		users[c.ReportedID] = nil
	}

	wp := workerpool.NewWorkerPool[entities.User](10)
	err = wp.FillMap(ctx, users, s.repo.GetUserByID)

	if err != nil {
		return nil, fmt.Errorf("failed to get info about users: %w", err)
	}

	for i, c := range complaints {
		complaints[i].Complainer = users[c.ComplainerID]
		complaints[i].Reported = users[c.ReportedID]
	}

	return complaints, nil
}
