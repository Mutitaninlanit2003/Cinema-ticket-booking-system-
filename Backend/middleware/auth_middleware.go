package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware verifies JWT token.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// -----------------------------------------
		// อ่าน JWT_SECRET จาก environment
		// -----------------------------------------

		jwtSecret := os.Getenv("JWT_SECRET")

		if jwtSecret == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "JWT_SECRET is not configured",
			})
			c.Abort()
			return
		}

		// -----------------------------------------
		// อ่าน Authorization header
		// -----------------------------------------

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header is required",
			})
			c.Abort()
			return
		}

		// -----------------------------------------
		// ตรวจสอบรูปแบบ
		//
		// Authorization: Bearer <token>
		// -----------------------------------------

		parts := strings.Fields(authHeader)

		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// -----------------------------------------
		// Parse JWT
		// -----------------------------------------

		token, err := jwt.Parse(
			tokenString,
			func(token *jwt.Token) (interface{}, error) {

				// ตรวจสอบว่าใช้ HS256 เท่านั้น
				if token.Method != jwt.SigningMethodHS256 {
					return nil, jwt.ErrSignatureInvalid
				}

				return []byte(jwtSecret), nil
			},
		)

		// -----------------------------------------
		// ตรวจสอบ Token
		// -----------------------------------------

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			c.Abort()
			return
		}

		// -----------------------------------------
		// ดึง claims
		// -----------------------------------------

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token claims",
			})
			c.Abort()
			return
		}

		// -----------------------------------------
		// ตรวจสอบ user_id
		// -----------------------------------------

		userID, ok := claims["user_id"].(string)

		if !ok || userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "user_id not found in token",
			})
			c.Abort()
			return
		}

		// -----------------------------------------
		// ตรวจสอบ role
		// -----------------------------------------

		role, ok := claims["role"].(string)

		if !ok || role == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "role not found in token",
			})
			c.Abort()
			return
		}

		// -----------------------------------------
		// ส่งข้อมูลให้ Handler ถัดไป
		// -----------------------------------------

		c.Set("user_id", userID)
		c.Set("role", role)

		c.Next()
	}
}

// RequireRole allows only users with the specified role.
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {

		// -----------------------------------------
		// ดึง role จาก AuthMiddleware
		// -----------------------------------------

		roleValue, exists := c.Get("role")

		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authentication required",
			})
			c.Abort()
			return
		}

		userRole, ok := roleValue.(string)

		if !ok || userRole == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid user role",
			})
			c.Abort()
			return
		}

		// -----------------------------------------
		// ตรวจสอบสิทธิ์
		// -----------------------------------------

		if userRole != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "forbidden: insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
