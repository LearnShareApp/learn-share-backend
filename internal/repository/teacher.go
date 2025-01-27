package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	internalErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
	"github.com/jmoiron/sqlx"
	"strings"
)

func (r *Repository) IsTeacherExistsByUserId(ctx context.Context, id int) (bool, error) {
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
		// Если ошибка не sql.ErrNoRows, возвращаем её
		return 0, fmt.Errorf("failed to select teacher: %w", err)
	}

	if err := r.db.QueryRowContext(ctx, insertQuery, userId).Scan(&teacherId); err != nil {
		return 0, fmt.Errorf("failed to insert teacher: %w", err)
	}
	return teacherId, nil
}

func (r *Repository) GetTeacherByUserId(ctx context.Context, id int) (*entities.Teacher, error) {
	const query = `SELECT teacher_id, user_id FROM teachers WHERE user_id = $1`

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

func (r *Repository) GetTeacherById(ctx context.Context, id int) (*entities.Teacher, error) {
	const query = `SELECT teacher_id, user_id FROM teachers WHERE teacher_id = $1`

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

func (r *Repository) GetAllTeachersDataFiltered(ctx context.Context, userId int, isUsersTeachers bool, category string, isFilteredByCategory bool) ([]entities.User, error) {
	// Base query
	baseQuery := `
    SELECT DISTINCT
        u.user_id,
        u.email,
        u.name,
        u.surname,
        u.registration_date,
        u.birthdate,
        u.avatar,
        t.teacher_id,
        s.skill_id,
        s.category_id,
        s.video_card_link,
        s.about,
        s.rate,
        s.is_active,
        c.name as category_name
    FROM users u
    INNER JOIN teachers t ON u.user_id = t.user_id
    INNER JOIN skills s ON t.teacher_id = s.teacher_id AND s.is_active
    INNER JOIN categories c ON s.category_id = c.category_id
	`

	// Создаем именованные параметры для sqlx
	namedParams := make(map[string]interface{})
	var conditions []string

	if isUsersTeachers {
		baseQuery += `
        INNER JOIN lessons l ON t.teacher_id = l.teacher_id 
        INNER JOIN statuses st ON l.status_id = st.status_id AND st.name = :status_name
    `
		conditions = append(conditions, "l.student_id = :user_id")
		namedParams["user_id"] = userId
		namedParams["status_name"] = "finished"
	}

	if isFilteredByCategory {
		conditions = append(conditions, "c.name = :category")
		namedParams["category"] = category
	}

	query := baseQuery
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Подготавливаем именованный запрос
	namedQuery, args, err := sqlx.Named(query, namedParams)
	if err != nil {
		return nil, fmt.Errorf("failed to build named query: %w", err)
	}

	// Конвертируем в $1, $2 формат для PostgreSQL
	query, args, err = sqlx.In(namedQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to convert named query: %w", err)
	}
	query = r.db.Rebind(query)

	// Временная структура для результатов
	type result struct {
		entities.User
		entities.Teacher
		entities.Skill
	}

	var rows []result
	err = r.db.SelectContext(ctx, &rows, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	// Мапа для группировки результатов
	usersMap := make(map[int]*entities.User)

	// Обработка результатов
	for _, row := range rows {
		user, exists := usersMap[row.User.Id]
		if !exists {
			user = &entities.User{
				Id:               row.User.Id,
				Email:            row.User.Email,
				Name:             row.User.Name,
				Surname:          row.User.Surname,
				Password:         row.User.Password,
				RegistrationDate: row.User.RegistrationDate,
				Birthdate:        row.User.Birthdate,
				IsTeacher:        false,
				TeacherData:      nil,
			}
			usersMap[row.User.Id] = user
		}

		if user.TeacherData == nil {
			user.IsTeacher = true
			user.TeacherData = &entities.Teacher{
				Id:     row.Teacher.Id,
				UserId: row.Teacher.UserId,
				Skills: make([]*entities.Skill, 0),
			}
		}

		if !hasSkill(user.TeacherData.Skills, row.Skill.Id) {
			user.TeacherData.Skills = append(user.TeacherData.Skills, &entities.Skill{
				Id:            row.Skill.Id,
				TeacherId:     row.Skill.TeacherId,
				CategoryId:    row.Skill.CategoryId,
				CategoryName:  row.Skill.CategoryName,
				VideoCardLink: row.Skill.VideoCardLink,
				About:         row.Skill.About,
				Rate:          row.Skill.Rate,
				IsActive:      row.Skill.IsActive,
			})
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
		if skill.Id == skillId {
			return true
		}
	}
	return false
}
