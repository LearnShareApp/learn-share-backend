package lesson

import (
	"context"
	"errors"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/pkg/workerpool"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	serviceErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
)

// GetStudentLessonList returns all lessons for this student.
func (s *LessonService) GetStudentLessonList(ctx context.Context, userID int) ([]*entities.Lesson, error) {
	// is user exists
	exists, err := s.repo.IsUserExistsByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existstance by id: %w", err)
	}

	if !exists {
		return nil, serviceErrs.ErrorUserNotFound
	}

	lessons, err := s.repo.GetLessonsByStudentID(ctx, userID)
	if err != nil {
		if errors.Is(err, serviceErrs.ErrorSelectEmpty) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get student's lessons: %w", err)
	}

	teachersUserIDs := make(map[int]*int)                         // all unique teacher's userIDs
	stateMachineItems := make(map[int]*entities.StateMachineItem) // all state machine items
	for _, l := range lessons {
		teachersUserIDs[l.TeacherID] = nil
		stateMachineItems[l.StateMachineItemID] = nil
	}

	wp := workerpool.NewWorkerPool[int](20)
	err = wp.FillMap(ctx, teachersUserIDs, func(ctx context.Context, id int) (*int, error) {
		u, err := s.repo.GetUserIDByTeacherID(ctx, id)
		return &u, err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get teacher's userIDs: %w", err)
	}

	for i, l := range lessons {
		lessons[i].TeacherUserData = &entities.User{ID: *teachersUserIDs[l.TeacherID]}
	}

	users := make(map[int]*entities.User) // all unique users (teacher) data
	for _, l := range lessons {
		users[l.TeacherUserData.ID] = nil
	}

	wp2 := workerpool.NewWorkerPool[entities.User](10)
	err = wp2.FillMap(ctx, users, s.repo.GetUserByID)

	if err != nil {
		return nil, fmt.Errorf("failed to get info about users: %w", err)
	}

	for i, l := range lessons {
		lessons[i].TeacherUserData = users[l.TeacherUserData.ID]
	}

	wpSateMachineItem := workerpool.NewWorkerPool[entities.StateMachineItem](10)
	err = wpSateMachineItem.FillMap(ctx, stateMachineItems, s.repo.GetStateMachineItemByID)
	if err != nil {
		return nil, fmt.Errorf("failed to fill state machine items: %w", err)
	}

	for i, l := range lessons {
		lessons[i].StateMachineItem = stateMachineItems[l.StateMachineItemID]
	}

	return lessons, nil
}
