package endpoints

import (
	"handworks-api/handlers"

	"github.com/gin-gonic/gin"
)
func AccountEndpoint(r* gin.RouterGroup, h * handlers.AccountHandler){
	customer := r.Group("/customer")
	{
		customer.POST("/signup", h.SignUpCustomer)
		customer.GET("/:id", h.GetCustomer)
		customer.PUT("/:id", h.UpdateCustomer)
		customer.DELETE("/:id", h.DeleteCustomer)
	}

	employee := r.Group("/employee")
	{
		employee.POST("/signup", h.SignUpEmployee)
		employee.GET("/:id", h.GetEmployee)
		employee.PUT("/:id", h.UpdateEmployee)
		employee.PUT("/:id/performance", h.UpdateEmployeePerformanceScore)
		employee.PUT("/:id/status", h.UpdateEmployeeStatus)
		employee.DELETE("/:id", h.DeleteEmployee)
	}
}
func InventoryEndpoint(r* gin.RouterGroup, h * handlers.InventoryHandler){
	r.POST("/", h.CreateItem)
	r.GET("/:id", h.GetItem)
	r.GET("/", h.ListAllItems)
	r.GET("/type/:type", h.ListItemsByType)
	r.GET("/status/:status", h.ListItemsByStatus)
	r.GET("/category/:category", h.ListItemsByCategory)
	r.PUT("/:id", h.UpdateItem)
	r.DELETE("/:id", h.DeleteItem)
}
func BookingEndpoint(r* gin.RouterGroup, h * handlers.BookingHandler){
	r.POST("/", h.CreateBooking)
	r.GET("/id/:id", h.GetBookingById)
	r.GET("/uid/:uid", h.GetBookingByUId)
	r.PUT("/:id", h.UpdateBooking)
	r.DELETE("/:id", h.DeleteBooking)
}
func PaymentEndpoint(r* gin.RouterGroup, h * handlers.PaymentHandler){
	r.POST("/quote", h.MakeQuotation)
	r.GET("/quotes/:customerId", h.GetAllQuotesFromCustomer)
}