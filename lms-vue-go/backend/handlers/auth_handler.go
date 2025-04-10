package handlers

import (
	"log"
	"net/http"
	"time"

	"lms-vue-go/backend/models"
	"lms-vue-go/backend/repository"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Kunci rahasia untuk JWT (dalam produksi, gunakan environment variable)
var jwtSecret = []byte("lms-secret-key-change-in-production")

// Durasi token JWT
const tokenDuration = 24 * time.Hour

// userRepo adalah repository untuk operasi user
// Akan diinisialisasi di setiap handler untuk memastikan koneksi DB sudah ada

// LoginRequest adalah struktur untuk request login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse adalah struktur untuk response login
type LoginResponse struct {
	Token string              `json:"token"`
	User  models.UserResponse `json:"user"`
}

// JWTClaims adalah struktur untuk JWT claims
type JWTClaims struct {
	UserID uint        `json:"user_id"`
	Role   models.Role `json:"role"`
	jwt.RegisteredClaims
}

// Login menangani proses login pengguna
func Login(c *gin.Context) {
	// Inisialisasi repository
	userRepo := repository.NewUserRepository()

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	// Cari pengguna berdasarkan username
	user, err := userRepo.FindByUsername(req.Username)
	if err != nil {
		log.Printf("Error in Login: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencari pengguna: " + err.Error()})
		return
	}

	// Jika pengguna tidak ditemukan atau password salah
	if user == nil || req.Password != user.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		return
	}

	// Buat token JWT
	token, err := generateJWT(user)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	// Kirim response
	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User:  user.ToResponse(),
	})
}

// Register menangani pendaftaran pengguna baru
func Register(c *gin.Context) {
	// Inisialisasi repository
	userRepo := repository.NewUserRepository()
	// Inisialisasi student repository
	studentRepo := repository.NewStudentRepository()

	// Struktur untuk binding request
	type RegisterRequest struct {
		Username string      `json:"username" binding:"required"`
		Password string      `json:"password" binding:"required"`
		Email    string      `json:"email" binding:"required"`
		Role     models.Role `json:"role"`
		Name     string      `json:"name"`  // Nama untuk student
		Class    string      `json:"class"` // Kelas untuk student
	}

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	// Buat user dari request
	user := models.User{
		Username: req.Username,
		Password: req.Password, // Pastikan password disimpan
		Email:    req.Email,
		Role:     req.Role,
	}

	// Cek apakah username sudah digunakan
	existingUser, err := userRepo.FindByUsername(user.Username)
	if err != nil {
		log.Printf("Error in Register: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa username: " + err.Error()})
		return
	}

	if existingUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username sudah digunakan"})
		return
	}

	// Default role adalah student jika tidak disebutkan
	if user.Role == "" {
		user.Role = models.RoleStudent
	}

	// Simpan user ke database
	err = userRepo.Create(&user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendaftarkan pengguna: " + err.Error()})
		return
	}

	// Jika user adalah student, buat record di tabel students
	if user.Role == models.RoleStudent {
		// Buat default name jika tidak disediakan
		name := req.Name
		if name == "" {
			name = req.Username // Gunakan username sebagai default name
		}

		// Buat default class jika tidak disediakan
		class := req.Class
		if class == "" {
			class = "Unassigned" // Default class
		}

		// Buat student record
		student := models.Student{
			Name:  name,
			Class: class,
			Email: user.Email,
		}

		// Simpan student ke database
		err = studentRepo.Create(&student, user.ID)
		if err != nil {
			log.Printf("Error creating student record: %v", err)
			// Tidak mengembalikan error ke client karena user sudah dibuat
			// Hanya log error untuk admin
		}
	}

	// Buat token JWT
	token, err := generateJWT(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	// Kirim response
	c.JSON(http.StatusCreated, LoginResponse{
		Token: token,
		User:  user.ToResponse(),
	})
}

// GetCurrentUser mengembalikan data pengguna yang sedang login
func GetCurrentUser(c *gin.Context) {
	// Inisialisasi repository
	userRepo := repository.NewUserRepository()

	// Ambil user dari context yang sudah diset oleh middleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tidak terautentikasi"})
		return
	}

	// Cari user berdasarkan ID
	user, err := userRepo.FindByID(userID.(uint))
	if err != nil {
		log.Printf("Error in GetCurrentUser: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pengguna: " + err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user.ToResponse()})
}

// Fungsi helper untuk hash password - tidak digunakan dalam development
// func hashPassword(password string) (string, error) {
// 	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
// 	return string(bytes), err
// }

// Fungsi helper untuk memeriksa password - tidak digunakan dalam development
// func checkPasswordHash(password, hash string) bool {
// 	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
// 	return err == nil
// }

// Fungsi helper untuk generate JWT
func generateJWT(user *models.User) (string, error) {
	// Set waktu kedaluwarsa token
	expirationTime := time.Now().Add(tokenDuration)

	// Buat claims
	claims := &JWTClaims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.Username,
		},
	}

	// Buat token dengan claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token dengan kunci rahasia
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
