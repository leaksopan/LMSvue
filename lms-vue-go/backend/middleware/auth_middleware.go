package middleware

import (
	"net/http"
	"strings"

	"lms-vue-go/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Kunci rahasia untuk JWT (harus sama dengan yang di auth_handler.go)
var jwtSecret = []byte("lms-secret-key-change-in-production")

// JWTClaims adalah struktur untuk JWT claims (sama dengan di auth_handler.go)
type JWTClaims struct {
	UserID uint        `json:"user_id"`
	Role   models.Role `json:"role"`
	jwt.RegisteredClaims
}

// AuthMiddleware adalah middleware untuk memeriksa token JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil token dari header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ditemukan"})
			c.Abort()
			return
		}

		// Format token: "Bearer <token>"
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// Parse token
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
			c.Abort()
			return
		}

		// Validasi token
		if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
			// Set user ID dan role ke context
			c.Set("userID", claims.UserID)
			c.Set("userRole", claims.Role)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
			c.Abort()
		}
	}
}

// ValidateToken validates a JWT token and returns the user ID
func ValidateToken(tokenString string) (uint, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return 0, err
	}

	// Validate token
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims.UserID, nil
	}

	return 0, jwt.ErrSignatureInvalid
}

// RoleMiddleware adalah middleware untuk memeriksa peran pengguna
func RoleMiddleware(roles ...models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil role dari context (yang sudah diset oleh AuthMiddleware)
		userRole, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Tidak terautentikasi"})
			c.Abort()
			return
		}

		// Periksa apakah role pengguna ada dalam daftar role yang diizinkan
		role := userRole.(models.Role)
		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}

		// Jika role tidak diizinkan
		c.JSON(http.StatusForbidden, gin.H{"error": "Tidak memiliki izin"})
		c.Abort()
	}
}

// SecurityHeaders adalah middleware untuk menambahkan security headers
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data:; font-src 'self'; connect-src 'self'")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	}
}
