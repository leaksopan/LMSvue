package handlers

import (
	"net/http"
	"strconv"

	"lms-vue-go/backend/models"
	"lms-vue-go/backend/repository"

	"github.com/gin-gonic/gin"
)

// studentRepo adalah repository untuk operasi student
// Akan diinisialisasi di setiap handler untuk memastikan koneksi DB sudah ada

// GetAllStudents mengembalikan daftar semua siswa
func GetAllStudents(c *gin.Context) {
	// Inisialisasi repository
	studentRepo := repository.NewStudentRepository()

	students, err := studentRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data siswa"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": students,
	})
}

// GetStudentByID mengembalikan data siswa berdasarkan ID
func GetStudentByID(c *gin.Context) {
	// Inisialisasi repository
	studentRepo := repository.NewStudentRepository()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	// Validasi ID
	if id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID harus positif"})
		return
	}

	// Cek apakah user memiliki akses ke data siswa
	userRole, exists := c.Get("userRole")
	userID, userExists := c.Get("userID")

	// Siswa hanya bisa melihat data dirinya sendiri
	if exists && userExists && userRole == models.RoleStudent {
		// Cari student berdasarkan userID
		studentData, err := studentRepo.FindByUserID(userID.(uint))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa data siswa"})
			return
		}

		// Jika student ditemukan dan ID tidak sama
		if studentData != nil && studentData.ID != uint(id) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Tidak memiliki izin untuk melihat data siswa lain"})
			return
		}
	}

	// Cari student berdasarkan ID
	student, err := studentRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data siswa"})
		return
	}

	if student == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Siswa tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": student})
}

// CreateStudent menambahkan data siswa baru
func CreateStudent(c *gin.Context) {
	// Inisialisasi repository
	studentRepo := repository.NewStudentRepository()
	var student models.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	// Ambil user ID dari context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tidak terautentikasi"})
		return
	}

	// Simpan student ke database
	err := studentRepo.Create(&student, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan data siswa"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": student})
}

// UpdateStudent mengupdate data siswa
func UpdateStudent(c *gin.Context) {
	// Inisialisasi repository
	studentRepo := repository.NewStudentRepository()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	// Cari student berdasarkan ID
	student, err := studentRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data siswa"})
		return
	}

	if student == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Siswa tidak ditemukan"})
		return
	}

	// Bind data baru
	var updatedStudent models.Student
	if err := c.ShouldBindJSON(&updatedStudent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	// Update data
	updatedStudent.ID = uint(id)
	err = studentRepo.Update(&updatedStudent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate data siswa"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": updatedStudent})
}

// DeleteStudent menghapus data siswa
func DeleteStudent(c *gin.Context) {
	// Inisialisasi repository
	studentRepo := repository.NewStudentRepository()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	// Cari student berdasarkan ID
	student, err := studentRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data siswa"})
		return
	}

	if student == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Siswa tidak ditemukan"})
		return
	}

	// Hapus data
	err = studentRepo.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data siswa"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Siswa berhasil dihapus"})
}

// GetCurrentStudentProfile mengembalikan profil siswa untuk user yang sedang login
func GetCurrentStudentProfile(c *gin.Context) {
	// Inisialisasi repository
	studentRepo := repository.NewStudentRepository()

	// Ambil user ID dari context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tidak terautentikasi"})
		return
	}

	// Ambil role dari context
	userRole, roleExists := c.Get("userRole")
	if !roleExists || userRole != models.RoleStudent {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya siswa yang dapat mengakses profil siswa"})
		return
	}

	// Cari student berdasarkan user ID
	student, err := studentRepo.FindByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data siswa"})
		return
	}

	if student == nil {
		// Jika tidak ada profil siswa, buat profil default
		// Ambil data user untuk mendapatkan email
		userRepo := repository.NewUserRepository()
		user, err := userRepo.FindByID(userID.(uint))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pengguna"})
			return
		}

		if user == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
			return
		}

		// Buat profil siswa baru
		newStudent := models.Student{
			Name:  user.Username, // Gunakan username sebagai nama default
			Class: "Unassigned",  // Default class
			Email: user.Email,
		}

		// Simpan student ke database
		err = studentRepo.Create(&newStudent, userID.(uint))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat profil siswa"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": newStudent})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": student})
}
