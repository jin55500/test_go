package user

import (
	"net/http"
	"test/go/database"
	Model "test/go/model"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// @Description ดึงข้อมูลโปรไฟล์ โดย auth
// @Tags user
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} Model.UserSwagger
// @Router /user/me [get]
func GetUser(c *gin.Context) {

	var user Model.UserSwagger
	database.Db.Raw("SELECT * FROM users WHERE id = ?", c.MustGet("id")).Scan(&user)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "get data success!",
		"data":    user,
	})
	return
}

// @Description แก้ไขข้อมูลส่วนตัวโดย auth
// @Description ข้อมูล object ที่ต้องการ ส่งเป็น raw json
// @Description {
// @Description  	username : ชื่อผู้ใช้ | ห้ามซ้ำ
// @Description  	password : รหัสผ่าน
// @Description  	name : ชื่อ
// @Description  	surname : นามกสุล
// @Description  	bankNumber : เลขที่บัญชี | 10 หลัก | ห้ามซ้ำ
// @Description }
// @Tags user
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @param User body Model.UserUpdateForm true "update data success!"
// @Success 200
// @Router /user/me [patch]
func UpdateUser(c *gin.Context) {

	var json Model.UserUpdateForm
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user Model.User
	database.Db.First(&user, c.MustGet("id").(float64))

	//check username
	var checkUserName Model.User
	if err := database.Db.Where("username = ? AND id != ?", json.Username, c.MustGet("id")).First(&checkUserName).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Username already exists",
			"data":    nil,
		})
		return
	}

	//check banknumber
	var checkBankNumber Model.User
	if err := database.Db.Where("bank_number = ? AND id != ?", json.BankNumber, c.MustGet("id")).First(&checkBankNumber).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "BankNumber already exists",
			"data":    nil,
		})
		return
	}

	// Update user data
	user.Username = json.Username
	user.Name = json.Name
	user.Surname = json.Surname
	user.BankNumber = json.BankNumber
	if json.Password != "" {
		passWordHash, _ := bcrypt.GenerateFromPassword([]byte(json.Password), 10)
		user.Password = string(passWordHash)
	}

	if err := database.Db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "update data success!",
		"data":    nil,
	})
	return

}
