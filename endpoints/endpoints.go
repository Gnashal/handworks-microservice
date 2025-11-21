package endpoints

import (
	"handworks-api/handlers"

	"github.com/gin-gonic/gin"
)
func AccountEndpoint(r* gin.RouterGroup, h * handlers.AccountHandler){}
func InventoryEndpoint(r* gin.RouterGroup, h * handlers.InventoryHandler){}
func BookingEndpoint(r* gin.RouterGroup, h * handlers.BookingHandler){}
func PaymentEndpoint(r* gin.RouterGroup, h * handlers.PaymentHandler){}