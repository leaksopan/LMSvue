package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lms-vue-go/backend/models"
)

// Untuk sementara menggunakan data dummy
var students = []models.Student{
	{ID: 1, Name: "Budi Santoso", Class: "10A", Email: "budi@example.com"},
	{ID: 2, Name: "Ani Wijaya", Class: "10A", Email: "ani@example.com"},
	{ID: 3, Name: "Dian Pratama", Class: "10B", Email: "dian@example.com"},
	{ID: 4, Name: "Rini Susanti", Class: "10B", Email: "rini@example.com"},
	{ID: 5, Name: "Ahmad Rizki", Class: "10C", Email: "ahmad@example.com"},
}

// GetAllStudents mengembalikan daftar semua siswa
func GetAllStudents(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": students,
	})
}

// GetStudentByID mengembalikan data siswa berdasarkan ID
func GetStudentByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	for _, student := range students {
		if int(student.ID) == id {
			c.JSON(http.StatusOK, gin.H{"data": student})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Siswa tidak ditemukan"})
} 