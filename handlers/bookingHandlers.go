package handlers

import (
	"net/http"

	"handworks-api/types"

	"github.com/gin-gonic/gin"
)

// CreateBooking godoc
// @Summary Create a new booking
// @Description Creates a booking record
// @Tags Booking
// @Accept json
// @Produce json
// @Param input body map[string]interface{} true "Booking info"
// @Success 200 {object} map[string]string
// @Router /booking [post]
func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var input types.CreateBookingEvent
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}
	_, err := h.Service.CreateBooking(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// GetBookingById godoc
// @Summary Get booking by ID
// @Description Retrieve booking information by its database ID
// @Tags Booking
// @Accept json
// @Produce json
// @Param id path string true "Booking ID"
// @Success 200 {object} map[string]interface{}
// @Router /booking/{id} [get]
func (h *BookingHandler) GetBookingById(c *gin.Context) {
	_ = h.Service.GetBookingById(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// GetBookingByUId godoc
// @Summary Get booking by user ID
// @Description Retrieve all bookings for a specific user
// @Tags Booking
// @Accept json
// @Produce json
// @Param uid path string true "User ID"
// @Success 200 {array} map[string]interface{}
// @Router /booking/user/{uid} [get]
func (h *BookingHandler) GetBookingByUId(c *gin.Context) {
	_ = h.Service.GetBookingByUId(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// UpdateBooking godoc
// @Summary Update a booking
// @Description Update booking information
// @Tags Booking
// @Accept json
// @Produce json
// @Param id path string true "Booking ID"
// @Param input body map[string]interface{} true "Updated booking info"
// @Success 200 {object} map[string]string
// @Router /booking/{id} [put]
func (h *BookingHandler) UpdateBooking(c *gin.Context) {
	_ = h.Service.UpdateBooking(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// DeleteBooking godoc
// @Summary Delete a booking
// @Description Remove booking by ID
// @Tags Booking
// @Accept json
// @Produce json
// @Param id path string true "Booking ID"
// @Success 200 {object} map[string]string
// @Router /booking/{id} [delete]
func (h *BookingHandler) DeleteBooking(c *gin.Context) {
	_ = h.Service.DeleteBooking(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
