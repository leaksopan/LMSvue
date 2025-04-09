package routes

import (
	// "net/http" // Dihapus karena tidak digunakan
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"lms-vue-go/backend/handlers"
)

// SetupRouter mengatur semua endpoint API
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Konfigurasi CORS untuk development
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		MaxAge:           12 * time.Hour,
	}))

	// Endpoint untuk health check
	r.GET("/api/status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "online",
			"message": "LMS API berjalan dengan baik",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// Grup untuk API
	api := r.Group("/api")
	{
		// Routes untuk siswa
		students := api.Group("/students")
		{
			students.GET("/", handlers.GetAllStudents)
			students.GET("/:id", handlers.GetStudentByID)
			// Tambahkan endpoint lain seperti POST, PUT, DELETE nanti
		}

		// Routes untuk soal
		questions := api.Group("/questions")
		{
			questions.GET("/", handlers.GetAllQuestions)
			questions.GET("/:id", handlers.GetQuestionByID)
			// Tambahkan endpoint lain nanti
		}
	}

	return r
} 