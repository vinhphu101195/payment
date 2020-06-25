package handle

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

var db *sql.DB

func GetPaymentMethod(ctx *gin.context) {
	conn, err := db.Conn(ctx)
	if err != nil {
		ctx.JSON(200, gin.H{"error": "500", "data": {"error": "Can not connect to database"}})
		return
	}

	row,err:= conn.QueryContext(ctx,"Select * from")
}
