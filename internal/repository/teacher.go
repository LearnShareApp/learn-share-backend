package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
	internalErrs "github.com/LearnShareApp/learn-share-backend/internal/errors"
)

func (r *Repository) IsTeacherExistsByUserId(ctx context.Context, id int) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM public.teachers WHERE teacher_id = $1)`

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, id)

	if err != nil {
		return false, fmt.Errorf("failed to check teacher existence: %w", err)
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
	const query = `SELECT teacher_id, user_id FROM public.teachers WHERE user_id = $1`

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

func (r *Repository) GetAllTeachersData(ctx context.Context) ([]entities.User, error) {
	const query = `
    SELECT
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
	FROM public.users u
    INNER JOIN public.teachers t ON u.user_id = t.user_id
    INNER JOIN public.skills s ON t.teacher_id = s.teacher_id AND s.is_active
    INNER JOIN public.categories c ON s.category_id = c.category_id`

	// Временная структура для хранения результатов запроса
	type result struct {
		entities.User
		entities.Teacher
		entities.Skill
	}

	var rows []result
	err := r.db.SelectContext(ctx, &rows, query)
	if err != nil {
		return nil, err
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

		skillExists := false
		for _, skill := range user.TeacherData.Skills {
			if skill.Id == row.Skill.Id {
				skillExists = true
				break
			}
		}

		if !skillExists {
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
