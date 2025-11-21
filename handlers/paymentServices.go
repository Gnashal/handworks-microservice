package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MakeQuotation godoc
// @Summary Create a quotation
// @Description Generate a new quotation for a customer
// @Tags Payment
// @Accept json
// @Produce json
// @Param input body map[string]interface{} true "Quotation details"
// @Success 200 {object} map[string]string
// @Router /payment/quotation [post]
func (h *PaymentHandler) MakeQuotation(c *gin.Context) {
	_ = h.Service.MakeQuotation(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// GetAllQuotesFromCustomer godoc
// @Summary Get all quotations for a customer
// @Description Retrieve all quotations associated with a specific customer
// @Tags Payment
// @Accept json
// @Produce json
// @Param customer_id path string true "Customer ID"
// @Success 200 {array} map[string]interface{}
// @Router /payment/quotations/customer/{customer_id} [get]
func (h *PaymentHandler) GetAllQuotesFromCustomer(c *gin.Context) {
	_ = h.Service.GetAllQuotesFromCustomer(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}