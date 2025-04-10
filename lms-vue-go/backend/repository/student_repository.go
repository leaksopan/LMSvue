package repository

import (
	"database/sql"
	"errors"
	"lms-vue-go/backend/config"
	"lms-vue-go/backend/models"
)

// StudentRepository handles database operations for students
type StudentRepository struct {
	DB *sql.DB
}

// NewStudentRepository creates a new student repository
func NewStudentRepository() *StudentRepository {
	return &StudentRepository{
		DB: config.DB,
	}
}

// FindAll returns all students
func (r *StudentRepository) FindAll() ([]models.Student, error) {
	query := `
		SELECT s.id, s.user_id, s.name, s.class, u.email
		FROM students s
		JOIN users u ON s.user_id = u.id
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var student models.Student
		err := rows.Scan(
			&student.ID,
			&student.UserID,
			&student.Name,
			&student.Class,
			&student.Email,
		)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}

// FindByID finds a student by ID
func (r *StudentRepository) FindByID(id uint) (*models.Student, error) {
	query := `
		SELECT s.id, s.user_id, s.name, s.class, u.email
		FROM students s
		JOIN users u ON s.user_id = u.id
		WHERE s.id = ?
	`

	var student models.Student
	err := r.DB.QueryRow(query, id).Scan(
		&student.ID,
		&student.UserID,
		&student.Name,
		&student.Class,
		&student.Email,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Student not found
		}
		return nil, err
	}

	return &student, nil
}

// FindByUserID finds a student by user ID
func (r *StudentRepository) FindByUserID(userID uint) (*models.Student, error) {
	query := `
		SELECT s.id, s.user_id, s.name, s.class, u.email
		FROM students s
		JOIN users u ON s.user_id = u.id
		WHERE s.user_id = ?
	`

	var student models.Student
	err := r.DB.QueryRow(query, userID).Scan(
		&student.ID,
		&student.UserID,
		&student.Name,
		&student.Class,
		&student.Email,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Student not found
		}
		return nil, err
	}

	return &student, nil
}

// Create creates a new student
func (r *StudentRepository) Create(student *models.Student, userID uint) error {
	// Set the UserID field in the student model
	student.UserID = userID

	// Check if student record already exists for this user
	existingStudent, err := r.FindByUserID(userID)
	if err != nil {
		return err
	}

	// If student record already exists, update it instead
	if existingStudent != nil {
		student.ID = existingStudent.ID
		return r.Update(student)
	}

	// Otherwise, create a new student record
	query := `
		INSERT INTO students (user_id, name, class)
		VALUES (?, ?, ?)
	`

	result, err := r.DB.Exec(query,
		userID,
		student.Name,
		student.Class,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	student.ID = uint(id)
	return nil
}

// Update updates an existing student
func (r *StudentRepository) Update(student *models.Student) error {
	// If student ID is not set, we can't update
	if student.ID == 0 {
		return errors.New("student ID is required for update")
	}

	query := `
		UPDATE students
		SET name = ?, class = ?, user_id = ?
		WHERE id = ?
	`

	_, err := r.DB.Exec(query,
		student.Name,
		student.Class,
		student.UserID,
		student.ID,
	)

	return err
}

// Delete deletes a student
func (r *StudentRepository) Delete(id uint) error {
	query := `DELETE FROM students WHERE id = ?`
	_, err := r.DB.Exec(query, id)
	return err
}
