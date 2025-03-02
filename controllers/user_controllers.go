package controllers

import (
	"net/http"

	"server-crud/database"
	"server-crud/middleware"
	"server-crud/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

// Fungsi untuk mengubah error validator ke format yang lebih jelas
func formatValidationError(err error) map[string]string {
	errors := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := e.Field()
			switch e.Tag() {
			case "required":
				errors[field] = field + " harus diisi"
			case "email":
				errors[field] = "Format email tidak valid"
			case "min":
				errors[field] = field + " minimal " + e.Param() + " karakter"
			case "max":
				errors[field] = field + " maksimal " + e.Param() + " karakter"
			default:
				errors[field] = "Format tidak valid"
			}
		}
	}
	return errors
}

// func isValidEmail(email string) bool {
// 	_, err := mail.ParseAddress(email)
// 	return err == nil
// }

func Home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ini halaman home"})
}

// func Register(c *gin.Context) {
// 	var user models.User
// 	if err := c.ShouldBindJSON(&user); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
//
// 	// Validasi panjang password minimal 6 karakter
// 	if len(user.Password) < 6 {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Password harus minimal 6 karakter"})
// 		return
// 	}
//
// 	// ✅ Validasi format email
// 	if !isValidEmail(user.Email) {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Format email tidak valid"})
// 		return
// 	}
//
// 	// 🔍 Cek apakah email sudah digunakan
// 	var existingUser models.User
// 	if err := database.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah digunakan"})
// 		return
// 	}
//
// 	// 🔍 Cek apakah username sudah digunakan
// 	if err := database.DB.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Username sudah digunakan"})
// 		return
// 	}
//
// 	// Hash password
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
// 		return
// 	}
// 	user.Password = string(hashedPassword)
//
// 	// Save user to database
// 	if err := database.DB.Create(&user).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
// 		return
// 	}
//
// 	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
// }

func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 🔥 Validasi otomatis dengan library validator
	if err := validate.Struct(user); err != nil {
		formattedErrors := formatValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": formattedErrors})
		return
	}

	// 🔥 Cek apakah username atau email sudah digunakan
	var existingUser models.User
	if err := database.DB.Where("username = ?", user.Username).Or("email = ?", user.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username atau Email sudah digunakan"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}
	user.Password = string(hashedPassword)

	// Save user to database
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}
}

func Login(c *gin.Context) {
	var inputUser models.User
	if err := c.ShouldBindJSON(&inputUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dbUser models.User
	if err := database.DB.Where("username = ?", inputUser.Username).First(&dbUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credential"})
		return
	}

	//compare password
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(inputUser.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credential"})
		return
	}

	// Gunakan fungsi GenerateToken dari middleware
	tokenString, err := middleware.GenerateToken(dbUser.Id, dbUser.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func Logout(c *gin.Context) {
	//ambil token dari header Authorization
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization token required"})
		return
	}

	//tambahkan token ke blacklist
	middleware.AddToBlacklist(tokenString)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
