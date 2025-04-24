package complaint

import (
	"context"
	"fmt"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
)

func (s *ComplaintService) CreateComplaint(ctx context.Context, complaint *entities.Complaint) error {
	// is complainer exists
	exists, err := s.repo.IsUserExistsByID(ctx, complaint.ComplainerID)
	if err != nil {
		return fmt.Errorf("failed to check user existstance by id: %w", err)
	}

	if !exists {
		return serviceErrs.ErrorUserNotFound
	}

	// is reported exists
	exists, err = s.repo.IsUserExistsByID(ctx, complaint.ReportedID)
	if err != nil {
		return fmt.Errorf("failed to check user existstance by id: %w", err)
	}

	if !exists {
		return serviceErrs.ErrorReportedUserNotFound
	}

	// is complainer != reported
	if complaint.ComplainerID == complaint.ReportedID {
		return serviceErrs.ErrorComplainerAndReportedSame
	}

	if err := s.repo.CreateComplaint(ctx, complaint); err != nil {
		return fmt.Errorf("failed to create complaint: %w", err)
	}

	return nil
}
