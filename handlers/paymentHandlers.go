package handlers

import (
	"handworks-api/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

// MakeQuotation godoc
// @Summary Create a quotation
// @Description Generate a new quotation for a customer
// @Tags Payment
// @Accept json
// @Produce json
// @Param input body types.QuoteRequest true "Quote details"
// @Success 200 {object} types.QuoteResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /payment/quote [post]
func (h *PaymentHandler) MakeQuotation(c *gin.Context) {
	var req types.QuoteRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.NewErrorResponse(err))
		return
	}
	res, err := h.Service.MakeQuotation(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.NewErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, res)
}

// GetAllQuotesFromCustomer godoc
// @Summary Get all quotations for a customer
// @Security BearerAuth
// @Description Retrieve all quotations associated with a specific customer
// @Tags Payment
// @Accept json
// @Produce json
// @Param customer_id path string true "Customer ID"
// @Success 200 {array} types.QuotesResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /payment/quote/{customerId} [get]
func (h *PaymentHandler) GetAllQuotesFromCustomer(c *gin.Context) {
	customerId := c.Param("customerId")

	res, err := h.Service.GetAllQuotesFromCustomer(c.Request.Context(), customerId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.NewErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, res)
}