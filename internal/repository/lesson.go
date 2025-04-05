package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	internalErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"

	"github.com/lib/pq"
)

func (r *Repository) CreateUnconfirmedLesson(ctx context.Context, lesson *entities.Lesson) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}
	defer tx.Rollback()

	query, args, err := r.sqlBuilder.
		Insert("lessons").
		Columns("student_id", "teacher_id", "category_id", "schedule_time_id").
		Values(lesson.StudentID, lesson.TeacherID, lesson.CategoryID, lesson.ScheduleTimeID).
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

	// book time
	if err = bookScheduleTime(ctx, tx, lesson.ScheduleTimeID); err != nil {
		return fmt.Errorf("failed to book schedule time: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

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
	const query = `
	SELECT
		lessons.lesson_id,
		lessons.student_id,
		lessons.teacher_id,
		lessons.category_id,
		lessons.schedule_time_id,
		lessons.status_id,
		lessons.price,
		categories.name as category_name,
		statuses.name as status_name,
    	schedule_times.datetime as schedule_time_datetime
	FROM lessons
		INNER JOIN categories ON lessons.category_id = categories.category_id
    	INNER JOIN statuses ON lessons.status_id = statuses.status_id
		INNER JOIN schedule_times ON lessons.schedule_time_id = schedule_times.schedule_time_id
	WHERE lesson_id = $1
	`

	type result struct {
		CategoryName         string    `db:"category_name"`
		StatusName           string    `db:"status_name"`
		ScheduleTimeDatetime time.Time `db:"schedule_time_datetime"`
		entities.Lesson
	}

	var resp result

	err := r.db.GetContext(ctx, &resp, query, id)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, internalErrs.ErrorSelectEmpty
	}

	if err != nil {
		return nil, fmt.Errorf("failed to extract lesson by id: %w", err)
	}

	lesson := resp.Lesson
	lesson.CategoryName = resp.CategoryName
	lesson.StatusName = resp.StatusName
	lesson.ScheduleTimeDatetime = resp.ScheduleTimeDatetime

	return &lesson, nil
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

	// Временная структура для хранения результатов запроса
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
		lesson, exists := lessonsMap[row.Lesson.Id]
		if !exists {
			lesson = &row.Lesson
			lesson.StatusName = row.StatusName
			lesson.CategoryName = row.CategoryName
			lesson.ScheduleTimeDatetime = row.ScheduleTimeDatetime

			lessonsMap[row.Lesson.Id] = lesson
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

	// Временная структура для хранения результатов запроса
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
		lesson, exists := lessonsMap[row.Lesson.Id]
		if !exists {
			lesson = &row.Lesson
			lesson.StatusName = row.StatusName
			lesson.CategoryName = row.CategoryName
			lesson.ScheduleTimeDatetime = row.ScheduleTimeDatetime

			lessonsMap[row.Lesson.Id] = lesson
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
