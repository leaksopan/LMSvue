package repository

import (
	"database/sql"
	"errors"
	"lms-vue-go/backend/config"
	"lms-vue-go/backend/models"
	"log"
)

// StudentAnswerRepository handles database operations for student answers
type StudentAnswerRepository struct {
	DB *sql.DB
}

// NewStudentAnswerRepository creates a new student answer repository
func NewStudentAnswerRepository() *StudentAnswerRepository {
	// Check if DB is initialized
	if config.DB == nil {
		log.Println("WARNING: Database connection is nil in StudentAnswerRepository")
	}
	return &StudentAnswerRepository{
		DB: config.DB,
	}
}

// FindByStudentAndQuestion finds an answer by student ID and question ID
func (r *StudentAnswerRepository) FindByStudentAndQuestion(studentID, questionID uint) (*models.StudentAnswer, error) {
	query := `
		SELECT id, student_id, question_id, answer, score
		FROM student_answers
		WHERE student_id = ? AND question_id = ?
	`

	var answer models.StudentAnswer
	var score sql.NullInt32

	err := r.DB.QueryRow(query, studentID, questionID).Scan(
		&answer.ID,
		&answer.StudentID,
		&answer.QuestionID,
		&answer.Answer,
		&score,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Answer not found
		}
		return nil, err
	}

	// Set score if present
	if score.Valid {
		scoreInt := int(score.Int32)
		answer.Score = &scoreInt
	}

	return &answer, nil
}

// FindByStudent finds all answers for a student
func (r *StudentAnswerRepository) FindByStudent(studentID uint) ([]models.StudentAnswer, error) {
	query := `
		SELECT id, student_id, question_id, answer, score
		FROM student_answers
		WHERE student_id = ?
	`

	rows, err := r.DB.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers []models.StudentAnswer
	for rows.Next() {
		var answer models.StudentAnswer
		var score sql.NullInt32

		err := rows.Scan(
			&answer.ID,
			&answer.StudentID,
			&answer.QuestionID,
			&answer.Answer,
			&score,
		)
		if err != nil {
			return nil, err
		}

		// Set score if present
		if score.Valid {
			scoreInt := int(score.Int32)
			answer.Score = &scoreInt
		}

		answers = append(answers, answer)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return answers, nil
}

// Create creates a new student answer
func (r *StudentAnswerRepository) Create(answer *models.StudentAnswer) error {
	query := `
		INSERT INTO student_answers (student_id, question_id, answer, score)
		VALUES (?, ?, ?, ?)
	`

	var scoreSQL sql.NullInt32
	if answer.Score != nil {
		scoreSQL = sql.NullInt32{Int32: int32(*answer.Score), Valid: true}
	}

	result, err := r.DB.Exec(query,
		answer.StudentID,
		answer.QuestionID,
		answer.Answer,
		scoreSQL,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	answer.ID = uint(id)
	return nil
}

// Update updates an existing student answer
func (r *StudentAnswerRepository) Update(answer *models.StudentAnswer) error {
	query := `
		UPDATE student_answers
		SET answer = ?, score = ?
		WHERE id = ?
	`

	var scoreSQL sql.NullInt32
	if answer.Score != nil {
		scoreSQL = sql.NullInt32{Int32: int32(*answer.Score), Valid: true}
	}

	_, err := r.DB.Exec(query,
		answer.Answer,
		scoreSQL,
		answer.ID,
	)

	return err
}

// Delete deletes a student answer
func (r *StudentAnswerRepository) Delete(id uint) error {
	query := `DELETE FROM student_answers WHERE id = ?`
	_, err := r.DB.Exec(query, id)
	return err
}

// FindAll returns all student answers with student and question details
func (r *StudentAnswerRepository) FindAll() ([]models.StudentAnswerWithDetails, error) {
	query := `
		SELECT sa.id, sa.student_id, sa.question_id, sa.answer, sa.score,
		       s.name as student_name, s.class as student_class, s.user_id,
		       q.question, q.type, q.score as question_score
		FROM student_answers sa
		JOIN students s ON sa.student_id = s.id
		JOIN questions q ON sa.question_id = q.id
		ORDER BY sa.id DESC
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers []models.StudentAnswerWithDetails
	for rows.Next() {
		var answer models.StudentAnswerWithDetails
		var score sql.NullInt32

		err := rows.Scan(
			&answer.ID,
			&answer.StudentID,
			&answer.QuestionID,
			&answer.Answer,
			&score,
			&answer.StudentName,
			&answer.StudentClass,
			&answer.UserID,
			&answer.QuestionText,
			&answer.QuestionType,
			&answer.QuestionScore,
		)
		if err != nil {
			return nil, err
		}

		// Set score if present
		if score.Valid {
			scoreInt := int(score.Int32)
			answer.Score = &scoreInt
		}

		answers = append(answers, answer)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return answers, nil
}

// FindByID finds a student answer by ID
func (r *StudentAnswerRepository) FindByID(id uint) (*models.StudentAnswer, error) {
	query := `
		SELECT id, student_id, question_id, answer, score
		FROM student_answers
		WHERE id = ?
	`

	var answer models.StudentAnswer
	var score sql.NullInt32

	err := r.DB.QueryRow(query, id).Scan(
		&answer.ID,
		&answer.StudentID,
		&answer.QuestionID,
		&answer.Answer,
		&score,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Answer not found
		}
		return nil, err
	}

	// Set score if present
	if score.Valid {
		scoreInt := int(score.Int32)
		answer.Score = &scoreInt
	}

	return &answer, nil
}
