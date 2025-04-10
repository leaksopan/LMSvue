package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"lms-vue-go/backend/config"
	"lms-vue-go/backend/models"
	"log"
)

// QuestionRepository handles database operations for questions
type QuestionRepository struct {
	DB *sql.DB
}

// NewQuestionRepository creates a new question repository
func NewQuestionRepository() *QuestionRepository {
	// Check if DB is initialized
	if config.DB == nil {
		log.Println("WARNING: Database connection is nil in QuestionRepository")
	}
	return &QuestionRepository{
		DB: config.DB,
	}
}

// FindAll returns all questions
func (r *QuestionRepository) FindAll() ([]models.Question, error) {
	// Check if DB is nil
	if r.DB == nil {
		log.Println("ERROR: Database connection is nil in FindAll")
		return nil, errors.New("database connection not initialized")
	}

	query := `
		SELECT id, type, question, options, answer, image_url, score
		FROM questions
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []models.Question
	for rows.Next() {
		var question models.Question
		var optionsJSON sql.NullString
		var imageURL sql.NullString
		var answer sql.NullString

		err := rows.Scan(
			&question.ID,
			&question.Type,
			&question.Question,
			&optionsJSON,
			&answer,
			&imageURL,
			&question.Score,
		)
		if err != nil {
			return nil, err
		}

		// Parse options JSON if present
		if optionsJSON.Valid && optionsJSON.String != "" {
			err = json.Unmarshal([]byte(optionsJSON.String), &question.Options)
			if err != nil {
				return nil, err
			}
		}

		// Set answer if present
		if answer.Valid {
			question.Answer = answer.String
		}

		// Set image URL if present
		if imageURL.Valid {
			question.ImageURL = imageURL.String
		}

		questions = append(questions, question)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return questions, nil
}

// FindByID finds a question by ID
func (r *QuestionRepository) FindByID(id uint) (*models.Question, error) {
	// Check if DB is nil
	if r.DB == nil {
		log.Println("ERROR: Database connection is nil in FindByID")
		return nil, errors.New("database connection not initialized")
	}

	query := `
		SELECT id, type, question, options, answer, image_url, score
		FROM questions
		WHERE id = ?
	`

	var question models.Question
	var optionsJSON sql.NullString
	var imageURL sql.NullString
	var answer sql.NullString

	err := r.DB.QueryRow(query, id).Scan(
		&question.ID,
		&question.Type,
		&question.Question,
		&optionsJSON,
		&answer,
		&imageURL,
		&question.Score,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Question not found
		}
		return nil, err
	}

	// Parse options JSON if present
	if optionsJSON.Valid && optionsJSON.String != "" {
		err = json.Unmarshal([]byte(optionsJSON.String), &question.Options)
		if err != nil {
			return nil, err
		}
	}

	// Set answer if present
	if answer.Valid {
		question.Answer = answer.String
	}

	// Set image URL if present
	if imageURL.Valid {
		question.ImageURL = imageURL.String
	}

	return &question, nil
}

// Create creates a new question
func (r *QuestionRepository) Create(question *models.Question) error {
	// Check if DB is nil
	if r.DB == nil {
		log.Println("ERROR: Database connection is nil in Create")
		return errors.New("database connection not initialized")
	}

	query := `
		INSERT INTO questions (type, question, options, answer, image_url, score)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	// Convert options to JSON
	var optionsJSON []byte
	var err error
	if len(question.Options) > 0 {
		optionsJSON, err = json.Marshal(question.Options)
		if err != nil {
			return err
		}
	}

	result, err := r.DB.Exec(query,
		question.Type,
		question.Question,
		optionsJSON,
		sql.NullString{String: question.Answer, Valid: question.Answer != ""},
		sql.NullString{String: question.ImageURL, Valid: question.ImageURL != ""},
		question.Score,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	question.ID = uint(id)
	return nil
}

// Update updates an existing question
func (r *QuestionRepository) Update(question *models.Question) error {
	// Check if DB is nil
	if r.DB == nil {
		log.Println("ERROR: Database connection is nil in Update")
		return errors.New("database connection not initialized")
	}

	query := `
		UPDATE questions
		SET type = ?, question = ?, options = ?, answer = ?, image_url = ?, score = ?
		WHERE id = ?
	`

	// Convert options to JSON
	var optionsJSON []byte
	var err error
	if len(question.Options) > 0 {
		optionsJSON, err = json.Marshal(question.Options)
		if err != nil {
			return err
		}
	}

	_, err = r.DB.Exec(query,
		question.Type,
		question.Question,
		optionsJSON,
		sql.NullString{String: question.Answer, Valid: question.Answer != ""},
		sql.NullString{String: question.ImageURL, Valid: question.ImageURL != ""},
		question.Score,
		question.ID,
	)

	return err
}

// Delete deletes a question
func (r *QuestionRepository) Delete(id uint) error {
	// Check if DB is nil
	if r.DB == nil {
		log.Println("ERROR: Database connection is nil in Delete")
		return errors.New("database connection not initialized")
	}

	query := `DELETE FROM questions WHERE id = ?`
	_, err := r.DB.Exec(query, id)
	return err
}
