package complaint

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
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

	var mu sync.Mutex
	eg, ctx := errgroup.WithContext(ctx)

	// limited count of parallel requests
	sem := make(chan struct{}, 10)

	for userID := range users {
		userIDCopy := userID // local copy for gorutine (shadow variable)
		eg.Go(func() error {
			sem <- struct{}{}        // book slot
			defer func() { <-sem }() // unbook slot

			user, err := s.repo.GetUserByID(ctx, userIDCopy)
			if err != nil {
				return err
			}

			mu.Lock()
			defer mu.Unlock()
			users[userIDCopy] = user

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("failed to get info about users: %w", err)
	}

	for i, c := range complaints {
		complaints[i].Complainer = users[c.ComplainerID]
		complaints[i].Reported = users[c.ReportedID]
	}

	return complaints, nil
}
