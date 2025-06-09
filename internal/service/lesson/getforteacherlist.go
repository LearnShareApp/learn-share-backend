package lesson

import (
	"context"
	"errors"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/pkg/workerpool"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
)

// GetTeacherLessonList returns all lessons for this teacher.
func (s *LessonService) GetTeacherLessonList(ctx context.Context, userID int) ([]*entities.Lesson, error) {
	// is user exists
	exists, err := s.repo.IsUserExistsByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existstance by id: %w", err)
	}

	if !exists {
		return nil, serviceErrs.ErrorUserNotFound
	}

	// get teacher by userID
	teacher, err := s.repo.GetTeacherByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return nil, serviceErrs.ErrorUserIsNotTeacher
		}

		return nil, fmt.Errorf("failed to get teacher by userID: %w", err)
	}

	// get lessons
	lessons, err := s.repo.GetLessonsByTeacherID(ctx, teacher.ID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get teacher's lessons: %w", err)
	}

	// all unique users
	users := make(map[int]*entities.User)
	for _, c := range lessons {
		users[c.StudentID] = nil
	}

	wp := workerpool.NewWorkerPool[entities.User](10)
	err = wp.FillMap(ctx, users, s.repo.GetUserByID)

	if err != nil {
		return nil, fmt.Errorf("failed to get info about users: %w", err)
	}

	for i, l := range lessons {
		lessons[i].StudentUserData = users[l.StudentID]
	}

	// all state machine items
	stateMachineItems := make(map[int]*entities.StateMachineItem)
	for _, l := range lessons {
		stateMachineItems[l.StateMachineItemID] = nil
	}

	wpSateMachineItem := workerpool.NewWorkerPool[entities.StateMachineItem](10)
	err = wpSateMachineItem.FillMap(ctx, stateMachineItems, s.repo.GetStateMachineItemByID)
	if err != nil {
		return nil, fmt.Errorf("failed to fill state machine items: %w", err)
	}

	for i, l := range lessons {
		lessons[i].StateMachineItem = stateMachineItems[l.StateMachineItemID]
	}

	//lesson.StateMachineItem, err = s.repo.GetStateMachineItemByID(ctx, lesson.StateMachineItemID)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to get statemachine item with %d id: %w", lesson.StateMachineItemID, err)
	//}

	// get lessons
	//lessons, err = s.repo.GetTeacherLessonsByTeacherID(ctx, teacher.ID)
	//if err != nil {
	//	if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
	//		return nil, nil
	//	}
	//
	//	return nil, fmt.Errorf("failed to get teacher's lessons: %w", err)
	//}

	return lessons, nil
}
