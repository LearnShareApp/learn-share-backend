package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"time"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	internalErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"

	"github.com/lib/pq"
)

func (r *Repository) BookLesson(ctx context.Context,
	scheduleTimeID,
	studentID,
	teacherID,
	categoryID int) error {

	stateMachine, err := r.getStateMachineByName(ctx, entities.LessonStateMachineName)
	if err != nil {
		return fmt.Errorf("failed to get lesson's statemachine: %w", err)
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}
	defer tx.Rollback()

	// create stateMachineItem
	itemID, err := r.insertStateMachineItem(ctx, tx, *stateMachine)
	if err != nil {
		return fmt.Errorf("failed to create state machine item: %w", err)
	}

	// book time
	if err = r.bookActiveScheduleTimeByID(ctx, tx, scheduleTimeID); err != nil {
		return fmt.Errorf("failed to book schedule time: %w", err)
	}

	// insert lesson
	query, args, err := r.sqlBuilder.
		Insert("lessons").
		Columns(
			"student_id",
			"teacher_id",
			"category_id",
			"schedule_time_id",
			"state_machine_item_id").
		Values(
			studentID,
			teacherID,
			categoryID,
			scheduleTimeID,
			itemID).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build insert query: %w", err)
	}

	// insert lesson
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			// error code 23505 mean unique_violation
			if pqErr.Code == "23505" {
				return internalErrs.ErrorNonUniqueData
			}
		}

		return fmt.Errorf("failed to insert lesson: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

// old
//func (r *Repository) CreateUnconfirmedLesson(ctx context.Context, lesson *entities.Lesson) error {
//	tx, err := r.db.BeginTxx(ctx, nil)
//	if err != nil {
//		return fmt.Errorf("error beginning transaction: %w", err)
//	}
//	defer tx.Rollback()
//
//	query, args, err := r.sqlBuilder.
//		Insert("lessons").
//		Columns("student_id", "teacher_id", "category_id", "schedule_time_id").
//		Values(lesson.StudentID, lesson.TeacherID, lesson.CategoryID, lesson.ScheduleTimeID).
//		ToSql()
//	if err != nil {
//		return fmt.Errorf("failed to build insert query: %w", err)
//	}
//
//	// insert lesson
//	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
//		if pqErr, ok := err.(*pq.Error); ok {
//			// error code 23505 mean unique_violation
//			if pqErr.Code == "23505" {
//				return internalErrs.ErrorNonUniqueData
//			}
//		}
//
//		return fmt.Errorf("failed to insert lesson: %w", err)
//	}
//
//	// book time
//	if err = bookScheduleTime(ctx, tx, lesson.ScheduleTimeID); err != nil {
//		return fmt.Errorf("failed to book schedule time: %w", err)
//	}
//
//	if err = tx.Commit(); err != nil {
//		return fmt.Errorf("error committing transaction: %w", err)
//	}
//
//	return nil
//}

func (r *Repository) IsLessonExistsById(ctx context.Context, id int) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM lessons WHERE lesson_id = $1)`

	var exists bool

	err := r.db.GetContext(ctx, &exists, query, id)
	if err != nil {
		return false, fmt.Errorf("failed to check lesson existence by lesson id: %w", err)
	}

	return exists, nil
}

func (r *Repository) GetLessonByID(ctx context.Context, id int) (*entities.Lesson, error) {
	query, args, err := r.sqlBuilder.
		Select(
			"l.lesson_id",
			"l.student_id",
			"l.teacher_id",
			"l.category_id",
			"l.schedule_time_id",
			"l.status_id",
			"l.price",
			"l.state_machine_item_id",
			"c.name as category_name",
			"statuses.name as status_name",
			"schedule_times.datetime as schedule_time_datetime",
		).
		From("lessons l").
		InnerJoin("categories c ON l.category_id = c.category_id").
		InnerJoin("statuses ON l.status_id = statuses.status_id").
		InnerJoin("schedule_times ON l.schedule_time_id = schedule_times.schedule_time_id").
		Where(squirrel.Eq{"l.lesson_id": id}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var lesson entities.Lesson
	err = r.db.GetContext(ctx, &lesson, query, args...)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internalErrs.ErrorSelectEmpty
		}

		return nil, fmt.Errorf("failed to extract lesson by id: %w", err)
	}

	return &lesson, nil
}

func (r *Repository) GetLessonsByTeacherID(ctx context.Context, teacherID int) ([]*entities.Lesson, error) {
	query, args, err := r.sqlBuilder.
		Select(
			"l.lesson_id",
			"l.student_id",
			"l.teacher_id",
			"l.category_id",
			"l.schedule_time_id",
			"l.status_id",
			"l.price",
			"l.state_machine_item_id",
			"c.name as category_name",
			"statuses.name as status_name",
			"schedule_times.datetime as schedule_time_datetime",
		).
		From("lessons l").
		InnerJoin("categories c ON l.category_id = c.category_id").
		InnerJoin("statuses ON l.status_id = statuses.status_id").
		InnerJoin("schedule_times ON l.schedule_time_id = schedule_times.schedule_time_id").
		Where(squirrel.Eq{"l.teacher_id": teacherID}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var lessons []*entities.Lesson
	err = r.db.SelectContext(ctx, &lessons, query, args...)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internalErrs.ErrorSelectEmpty
		}

		return nil, fmt.Errorf("failed to extract lessons by teacherID: %w", err)
	}

	return lessons, nil
}

func (r *Repository) GetLessonsByStudentID(ctx context.Context, studentID int) ([]*entities.Lesson, error) {
	query, args, err := r.sqlBuilder.
		Select(
			"l.lesson_id",
			"l.student_id",
			"l.teacher_id",
			"l.category_id",
			"l.schedule_time_id",
			"l.status_id",
			"l.price",
			"l.state_machine_item_id",
			"c.name as category_name",
			"statuses.name as status_name",
			"schedule_times.datetime as schedule_time_datetime",
		).
		From("lessons l").
		InnerJoin("categories c ON l.category_id = c.category_id").
		InnerJoin("statuses ON l.status_id = statuses.status_id").
		InnerJoin("schedule_times ON l.schedule_time_id = schedule_times.schedule_time_id").
		Where(squirrel.Eq{"l.student_id": studentID}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var lessons []*entities.Lesson
	err = r.db.SelectContext(ctx, &lessons, query, args...)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internalErrs.ErrorSelectEmpty
		}

		return nil, fmt.Errorf("failed to extract lessons by studentID: %w", err)
	}

	return lessons, nil
}

func (r *Repository) IsFinishedLessonExistsByTeacherIdAndStudentIdAndCategoryId(ctx context.Context, teacherId int, studentId int, categoryId int) (bool, error) {
	const query = `
	SELECT EXISTS(
		SELECT 1 FROM lessons l
		LEFT JOIN statuses st ON l.status_id = st.status_id
		WHERE l.teacher_id = $1 AND l.student_id = $2 AND l.category_id = $3 AND st.name = $4
	)
	`

	var exists bool

	err := r.db.GetContext(ctx, &exists, query, teacherId, studentId, categoryId, entities.FinishedStatusName)
	if err != nil {
		return false, fmt.Errorf("failed to check finished lesson existence by teacher id, student id and category id: %w", err)
	}

	return exists, nil
}

func (r *Repository) ChangeLessonStatus(ctx context.Context, lessonId int, statusId int) error {
	const query = `
	UPDATE lessons SET status_id = $2 WHERE lesson_id = $1
	`

	if _, err := r.db.ExecContext(ctx, query, lessonId, statusId); err != nil {
		return fmt.Errorf("failed to update lesson status: %w", err)
	}

	return nil
}

func (r *Repository) GetTeacherLessonsByTeacherID(ctx context.Context, id int) ([]*entities.Lesson, error) {
	const query = `
    SELECT
		lessons.lesson_id,
		lessons.student_id,
		lessons.teacher_id,
		lessons.category_id,
		lessons.schedule_time_id,
		lessons.status_id,
		lessons.price,
		lessons.state_machine_item_id,
		users.user_id,
		users.email,
		users.name,
		users.surname,
		users.avatar,
		categories.name as category_name,
		statuses.name as status_name,
    	schedule_times.datetime as schedule_time_datetime
	FROM lessons
    INNER JOIN users ON lessons.student_id = users.user_id
    INNER JOIN categories ON lessons.category_id = categories.category_id
    INNER JOIN statuses ON lessons.status_id = statuses.status_id
	INNER JOIN schedule_times ON lessons.schedule_time_id = schedule_times.schedule_time_id
    WHERE lessons.teacher_id = $1`

	type result struct {
		CategoryName         string    `db:"category_name"`
		StatusName           string    `db:"status_name"`
		ScheduleTimeDatetime time.Time `db:"schedule_time_datetime"`
		entities.Lesson
		entities.User
	}

	var rows []result

	err := r.db.SelectContext(ctx, &rows, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internalErrs.ErrorSelectEmpty
		}

		return nil, err
	}

	// Мапа для группировки результатов
	lessonsMap := make(map[int]*entities.Lesson)

	// Обработка результатов
	for _, row := range rows {
		lesson, exists := lessonsMap[row.Lesson.ID]
		if !exists {
			lesson = &row.Lesson
			lesson.StatusName = row.StatusName
			lesson.CategoryName = row.CategoryName
			lesson.ScheduleTimeDatetime = row.ScheduleTimeDatetime

			lessonsMap[row.Lesson.ID] = lesson
		}

		if lesson.StudentUserData == nil {
			lesson.StudentUserData = &row.User
		}
	}

	lessons := make([]*entities.Lesson, 0, len(lessonsMap))
	for _, lesson := range lessonsMap {
		lessons = append(lessons, lesson)
	}

	return lessons, nil
}

func (r *Repository) GetStudentLessonsByUserID(ctx context.Context, id int) ([]*entities.Lesson, error) {
	const query = `
	   SELECT
		lessons.lesson_id,
		lessons.student_id,
		lessons.teacher_id,
		lessons.category_id,
		lessons.schedule_time_id,
		lessons.status_id,
		lessons.price,
		lessons.state_machine_item_id,
		users.user_id,
		users.email,
		users.name,
		users.surname,
		users.avatar,
		categories.name as category_name,
		statuses.name as status_name,
		schedule_times.datetime as schedule_time_datetime
		FROM lessons
		   INNER JOIN teachers ON lessons.teacher_id = teachers.teacher_id
		   INNER JOIN users ON teachers.user_id = users.user_id
		   INNER JOIN categories ON lessons.category_id = categories.category_id
		   INNER JOIN statuses ON lessons.status_id = statuses.status_id
		   INNER JOIN schedule_times ON lessons.schedule_time_id = schedule_times.schedule_time_id
	   WHERE lessons.student_id = $1`

	type result struct {
		CategoryName         string    `db:"category_name"`
		StatusName           string    `db:"status_name"`
		ScheduleTimeDatetime time.Time `db:"schedule_time_datetime"`
		entities.Lesson
		entities.User
	}

	var rows []result

	err := r.db.SelectContext(ctx, &rows, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internalErrs.ErrorSelectEmpty
		}

		return nil, err
	}

	lessonsMap := make(map[int]*entities.Lesson)

	// Обработка результатов
	for _, row := range rows {
		lesson, exists := lessonsMap[row.Lesson.ID]
		if !exists {
			lesson = &row.Lesson
			lesson.StatusName = row.StatusName
			lesson.CategoryName = row.CategoryName
			lesson.ScheduleTimeDatetime = row.ScheduleTimeDatetime

			lessonsMap[row.Lesson.ID] = lesson
		}

		if lesson.TeacherUserData == nil {
			lesson.TeacherUserData = &row.User
		}
	}

	lessons := make([]*entities.Lesson, 0, len(lessonsMap))
	for _, lesson := range lessonsMap {
		lessons = append(lessons, lesson)
	}

	return lessons, nil
}
