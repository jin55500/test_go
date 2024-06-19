package auth

import (
	"net/http"
	"os"
	"time"

	"test/go/database"
	Model "test/go/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type LoginForm struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

var secretKey []byte

// @Description สมัครสมาชิก
// @Description ข้อมูล object ที่ต้องการ ส่งเป็น raw json
// @Description {
// @Description  	username : ชื่อผู้ใช้ | ห้ามซ้ำ
// @Description  	password : รหัสผ่าน
// @Description  	name : ชื่อ
// @Description  	surname : นามกสุล
// @Description  	bankNumber : เลขที่บัญชี | 10 หลัก | ห้ามซ้ำ
// @Description }
// @Tags user
// @Accept json
// @Produce json
// @param User body Model.UserCreateForm true "register success!"
// @Success 200
// @Router /user/register [post]
func RegisterUser(c *gin.Context) {

	var json Model.UserCreateForm
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//check username exist
	var checkUserName Model.User
	if err := database.Db.Where("username = ?", json.Username).First(&checkUserName).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Username already exists",
			"data":    nil,
		})
		return
	}

	//check banknumber exist
	var checkBankNumber Model.User
	if err := database.Db.Where("bank_number = ?", json.BankNumber).First(&checkBankNumber).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "BankNumber already exists",
			"data":    nil,
		})
		return
	}

	//hash password
	passWordHash, _ := bcrypt.GenerateFromPassword([]byte(json.Password), 10)

	//save data to db
	user := Model.User{
		Username:   json.Username,
		Password:   string(passWordHash),
		Name:       json.Name,
		Surname:    json.Surname,
		BankNumber: json.BankNumber,
	}
	if err := database.Db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err,
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    nil,
		"message": "register success!",
	})
	return
}

// @Description เข้าสู่ระบบ
// @Description ข้อมูล object ที่ต้องการ ส่งเป็น raw json
// @Description {
// @Description  	username : ชื่อผู้ใช้
// @Description  	password : รหัสผ่าน
// @Description }
// @Tags user
// @Accept json
// @Produce json
// @param User body LoginForm true "login success!"
// @Success 200
// @Router /user/login [post]
func Login(c *gin.Context) {

	var json LoginForm
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//check user exist and password
	var user Model.User
	database.Db.Where("username = ?", json.Username).First(&user)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(json.Password))
	if user.ID == 0 || err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "username or password invalid",
			"data":    nil,
		})
		return
	}

	//create token
	secretKey = []byte(os.Getenv("my_secret_key"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Minute * 30).Unix(), //set time out
	})
	tokenString, err := token.SignedString(secretKey)

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"accessToken": tokenString,
		"message":     "login success!",
	})
	return
}
