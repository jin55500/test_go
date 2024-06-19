package accounting

import (
	"net/http"
	"test/go/database"
	Model "test/go/model"

	"github.com/gin-gonic/gin"
)

// @Description โอนเครดิตให้ user อื่น
// @Description ข้อมูล object ที่ต้องการ ส่งเป็น raw json
// @Description {
// @Description  	bankNumber : เลขที่บัญชี
// @Description  	credit : จำนวนเครดิตที่ต้องการโอน
// @Description }
// @Tags accounting
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @param accounting body Model.TransferCreateForm true "tranfer credit success!"
// @Success 200
// @Router /accounting/transfer [post]
func Transfer(c *gin.Context) {
	var json Model.TransferCreateForm
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//get user transfer
	var UserTransfer Model.User
	database.Db.First(&UserTransfer, c.MustGet("id").(float64))
	if UserTransfer.Credit < json.Credit {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "credit not enough",
			"data":    nil,
		})
		return
	}

	//get user receiver
	var UserReceiver Model.User
	if err := database.Db.Where("bank_number = ?", json.BankNumber).First(&UserReceiver).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "bankNumber not found",
			"data":    nil,
		})
		return
	}

	if UserTransfer.ID == UserReceiver.ID {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "cant transfer to your self",
			"data":    nil,
		})
		return
	}

	//update credit
	UserTransfer.Credit -= json.Credit
	database.Db.Save(&UserTransfer)
	UserReceiver.Credit += json.Credit
	database.Db.Save(&UserReceiver)

	//create accounting
	accounting := Model.Accounting{
		UserTransfer: UserTransfer.ID,
		Credit:       json.Credit,
		UserReceiver: UserReceiver.ID,
		BankNumber:   UserReceiver.BankNumber,
	}
	if err := database.Db.Create(&accounting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err,
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "tranfer credit success!",
		"data":    nil,
	})
	return
}

// @Description ประวัติการโอนเครดิต โดย auth
// @Description ข้อมูล params ที่ต้องการ
// @Description  	date_start : วันที่เริ่มต้น format ("2024-06-19")
// @Description  	date_end : วันที่สิ้นสุด format ("2024-06-20")
// @Description  	===============================================================
// @Description  	response ตัวแปร Type = ชนิดของ transection
// @Description  	receive = ได้รับ credit
// @Description  	transfer = โอน credit
// @Tags accounting
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param date_start query string true "Start date (format: '2024-06-19')"
// @Param date_end query string true "End date (format: '2024-06-20')"
// @Success 200
// @Router /accounting/transfer-list [get]
func GetTransferList(c *gin.Context) {
	var accounting []Model.AccountingResult
	result := database.Db.Raw(
		`SELECT 
			*,
			(
				CASE
					WHEN user_transfer = ? THEN 
						'transfer'
					WHEN user_receiver = ? THEN 
						'receive'
				END
			) as type
		FROM accountings
		WHERE ( user_receiver = ? OR user_transfer = ? ) 
		AND ( created_at BETWEEN ? AND ? )`,
		c.MustGet("id"),
		c.MustGet("id"),
		c.MustGet("id"),
		c.MustGet("id"),
		c.Query("date_start"),
		c.Query("date_end")).Scan(&accounting)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": result.Error.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "get data success!",
		"data":    accounting,
	})
	return
}
