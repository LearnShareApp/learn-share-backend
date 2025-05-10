package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	internalErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/Masterminds/squirrel"

	"github.com/jmoiron/sqlx"
)

func (r *Repository) IsTeacherExistsByUserID(ctx context.Context, id int) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM teachers WHERE user_id = $1)`

	var exists bool

	err := r.db.GetContext(ctx, &exists, query, id)
	if err != nil {
		return false, fmt.Errorf("failed to check teacher existence by user id: %w", err)
	}

	return exists, nil
}

func (r *Repository) IsTeacherExistsById(ctx context.Context, id int) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM teachers WHERE teacher_id = $1)`

	var exists bool

	err := r.db.GetContext(ctx, &exists, query, id)
	if err != nil {
		return false, fmt.Errorf("failed to check teacher existence by teacher id: %w", err)
	}

	return exists, nil
}

func (r *Repository) CreateTeacher(ctx context.Context, userId int) error {
	const query = `
	INSERT INTO teachers (user_id) 
	VALUES ($1)
	`

	if _, err := r.db.ExecContext(ctx, query, userId); err != nil {
		return fmt.Errorf("failed to insert teacher: %w", err)
	}

	return nil
}

// CreateTeacherIfNotExists check is teacher exists in db, create if not and return teacher_id.
func (r *Repository) CreateTeacherIfNotExists(ctx context.Context, userId int) (int, error) {
	const (
		selectQuery = `
		SELECT teacher_id FROM teachers WHERE user_id = $1
		`

		insertQuery = `
		INSERT INTO teachers (user_id) 
		VALUES ($1)
		RETURNING teacher_id
		`
	)

	var teacherId int

	err := r.db.GetContext(ctx, &teacherId, selectQuery, userId)
	if err == nil {
		return teacherId, nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("failed to select teacher: %w", err)
	}

	if err := r.db.QueryRowContext(ctx, insertQuery, userId).Scan(&teacherId); err != nil {
		return 0, fmt.Errorf("failed to insert teacher: %w", err)
	}

	return teacherId, nil
}

func (r *Repository) GetTeacherIdByUserId(ctx context.Context, id int) (int, error) {
	const query = `
		SELECT 
		    teacher_id
		FROM teachers 
		WHERE user_id = $1`

	var teacherId int
	err := r.db.GetContext(ctx, &teacherId, query, id)

	if errors.Is(err, sql.ErrNoRows) {
		return 0, internalErrs.ErrorSelectEmpty
	}

	if err != nil {
		return 0, fmt.Errorf("failed to find teacherId by user id: %w", err)
	}

	return teacherId, nil
}

func (r *Repository) GetTeacherByUserID(ctx context.Context, id int) (*entities.Teacher, error) {
	const query = `
		SELECT 
		    teacher_id, 
		    user_id,
		    rate,
		    reviews_count
		FROM teachers 
		WHERE user_id = $1`

	var teacher entities.Teacher
	err := r.db.GetContext(ctx, &teacher, query, id)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, internalErrs.ErrorSelectEmpty
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find teacher by user id: %w", err)
	}

	return &teacher, nil
}

func (r *Repository) GetTeacherByID(ctx context.Context, id int) (*entities.Teacher, error) {
	const query = `
		SELECT 
    		teacher_id, 
    		user_id,
    		rate,
		    reviews_count
		FROM teachers WHERE teacher_id = $1`

	var teacher entities.Teacher
	err := r.db.GetContext(ctx, &teacher, query, id)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, internalErrs.ErrorSelectEmpty
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find teacher by id: %w", err)
	}

	return &teacher, nil
}

func (r *Repository) GetShortStatTeacherByID(ctx context.Context, teacherId int) (*entities.TeacherStatistic, error) {
	const query = `
    SELECT 
        COUNT(DISTINCT l.lesson_id) FILTER (WHERE st.name = $1) as count_of_finished_lesson,
		COUNT(DISTINCT l.student_id) FILTER (WHERE st.name = $1) as count_of_students
    FROM lessons l
    INNER JOIN statuses st ON st.status_id = l.status_id
    WHERE l.teacher_id = $2
    `

	var stat entities.TeacherStatistic

	err := r.db.GetContext(ctx, &stat, query,
		entities.FinishedStatusName, // $1
		teacherId,                   // $2
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("teacher statistic not found: %w", err)
		}

		return nil, fmt.Errorf("failed to calculate teacher's statistic by teacherId: %w", err)
	}

	return &stat, nil
}

func (r *Repository) GetShortTeacherDatasByIDs(ctx context.Context, teacherIDs map[int]bool) ([]entities.User, error) {
	var ids []int
	for id := range teacherIDs {
		ids = append(ids, id)
	}

	query, args, err := r.sqlBuilder.
		Select(
			"t.teacher_id",
			"t.user_id",
			"u.name",
			"u.surname",
			"u.avatar",
		).
		From("teachers t").
		InnerJoin("users u ON t.user_id = u.user_id").
		Where(squirrel.Eq{"t.teacher_id": ids}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	type shortTeacherData struct {
		TeacherID int    `db:"teacher_id"`
		UserID    int    `db:"user_id"`
		Name      string `db:"name"`
		Surname   string `db:"surname"`
		Avatar    string `db:"avatar"`
	}

	var rows []shortTeacherData
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf("failed to execute teacher by ids: %w", err)
	}

	users := make([]entities.User, 0, len(rows))
	for _, row := range rows {
		users = append(users, entities.User{
			ID:        row.UserID,
			Name:      row.Name,
			Surname:   row.Surname,
			Avatar:    row.Avatar,
			IsTeacher: true,
			TeacherData: &entities.Teacher{
				ID: row.TeacherID,
			},
		})
	}

	return users, nil
}

func (r *Repository) GetAllTeachersDataFiltered(ctx context.Context, userId int, isUsersTeachers bool, category string, isFilteredByCategory bool) ([]entities.User, error) {
	// Base query
	baseQuery := `
	WITH teacher_stats AS (
		SELECT
			l.teacher_id,
			COUNT(DISTINCT l.lesson_id) FILTER (WHERE st.name = :finished_status_name) as count_of_finished_lesson,
			COUNT(DISTINCT l.student_id) FILTER (WHERE st.name = :finished_status_name) as count_of_students
		FROM lessons l
		LEFT JOIN statuses st ON l.status_id = st.status_id
		GROUP BY l.teacher_id
	)

	SELECT
		u.user_id,
		u.email,
		u.name,
		u.surname,
		u.registration_date,
		u.birthdate,
		u.avatar,
		t.teacher_id,
		t.rate,
		t.reviews_count,
		s.skill_id,
		s.category_id,
		s.video_card_link,
		s.about,
		s.rate,
		s.total_rate_score,
		s.reviews_count,
		s.is_active,
		c.name as category_name,
		COALESCE(ts.count_of_finished_lesson, 0) as count_of_finished_lesson,
		COALESCE(ts.count_of_students, 0) as count_of_students
	FROM users u
	INNER JOIN teachers t ON u.user_id = t.user_id
	INNER JOIN skills s ON t.teacher_id = s.teacher_id
	INNER JOIN categories c ON s.category_id = c.category_id
	LEFT JOIN teacher_stats ts ON t.teacher_id = ts.teacher_id
	LEFT JOIN lessons l ON t.teacher_id = l.teacher_id
	LEFT JOIN statuses st ON l.status_id = st.status_id
	WHERE s.is_active
	`

	// named params for query
	namedParams := make(map[string]interface{})
	namedParams["finished_status_name"] = entities.FinishedStatusName

	var conditions []string

	if isUsersTeachers {
		conditions = append(conditions, "st.name = :finished_status_name")
		conditions = append(conditions, "l.student_id = :user_id")
		namedParams["user_id"] = userId
	}

	if isFilteredByCategory {
		conditions = append(conditions, "c.name = :category")
		namedParams["category"] = category
	}

	query := baseQuery
	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	// prepare named query
	namedQuery, args, err := sqlx.Named(query, namedParams)
	if err != nil {
		return nil, fmt.Errorf("failed to build named query: %w", err)
	}

	// converting into $1, $2, ... PostgreSQL format
	query = r.db.Rebind(namedQuery)

	type result struct {
		entities.User
		entities.Teacher
		entities.Skill
		entities.TeacherStatistic
	}

	var rows []result

	err = r.db.SelectContext(ctx, &rows, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []entities.User{}, nil
		}

		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	// Мапа для группировки результатов
	usersMap := make(map[int]*entities.User)

	// Обработка результатов
	for _, row := range rows {
		user, exists := usersMap[row.User.ID]
		if !exists {
			user = &row.User
			usersMap[row.User.ID] = user
		}

		if user.TeacherData == nil {
			user.IsTeacher = true
			user.TeacherData = &entities.Teacher{
				ID:          row.Teacher.ID,
				UserID:      row.Teacher.UserID,
				Skills:      make([]*entities.Skill, 0),
				TeacherStat: row.TeacherStatistic,
			}
		}

		if !hasSkill(user.TeacherData.Skills, row.Skill.ID) {
			user.TeacherData.Skills = append(user.TeacherData.Skills, &row.Skill)
		}
	}

	users := make([]entities.User, 0, len(usersMap))
	for _, user := range usersMap {
		users = append(users, *user)
	}

	return users, nil
}

func hasSkill(skills []*entities.Skill, skillId int) bool {
	for _, skill := range skills {
		if skill.ID == skillId {
			return true
		}
	}

	return false
}
