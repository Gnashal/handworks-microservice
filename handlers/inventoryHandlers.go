package handlers

import (
	"handworks-api/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateItem godoc
// @Summary Create a new inventory item
// @Description Adds a new item to inventory
// @Tags Inventory
// @Accept json
// @Produce json
// @Param input body types.CreateItemRequest true "Item info"
// @Success 200 {object} types.InventoryItem
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /inventory [post]
func (h *InventoryHandler) CreateItem(c *gin.Context) {
	var req types.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.NewErrorResponse(err))
		return
	}

	resp, err := h.Service.CreateItem(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetItem godoc
// @Summary Get an item by ID
// @Description Retrieve a single inventory item
// @Tags Inventory
// @Produce json
// @Param id path string true "Item ID"
// @Success 200 {object} types.InventoryItem
// @Failure 404 {object} types.ErrorResponse
// @Router /inventory/{id} [get]
func (h *InventoryHandler) GetItem(c *gin.Context) {
	id := c.Param("id")

	resp, err := h.Service.GetItem(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, types.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, resp)
}
// GetItem godoc
// @Summary Get all items
// @Description Retrieve all inventory items
// @Tags Inventory
// @Produce json
// @Param id path string true "Item ID"
// @Success 200 {object} []types.InventoryItem
// @Failure 404 {object} types.ErrorResponse
// @Router /inventory/ [get]
func (h *InventoryHandler) GetItems(c *gin.Context) {
	resp, err := h.Service.GetItems(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, types.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, resp)
}
// ListItemsByType godoc
// @Summary List items filtered by type
// @Description Get all inventory items matching the given type
// @Tags Inventory
// @Produce json
// @Param type path string true "Item type (resource, equipment)"
// @Success 200 {object} []types.InventoryItem
// @Failure 400 {object} types.ErrorResponse
// @Router /inventory/type/{type} [get]
func (h *InventoryHandler) ListItemsByType(c *gin.Context) {
	itemType := c.Param("type")

	resp, err := h.Service.ListItemsByType(c.Request.Context(), itemType)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListItemsByStatus godoc
// @Summary List items filtered by status
// @Description Get all inventory items with the given stock status
// @Tags Inventory
// @Produce json
// @Param status path string true "Status (high, low, danger, out_of_stock)"
// @Success 200 {object} []types.InventoryItem
// @Failure 400 {object} types.ErrorResponse
// @Router /inventory/status/{status} [get]
func (h *InventoryHandler) ListItemsByStatus(c *gin.Context) {
	status := c.Param("status")

	resp, err := h.Service.ListItemsByStatus(c.Request.Context(), status)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListItemsByCategory godoc
// @Summary List items filtered by category
// @Description Retrieve all inventory items in a category
// @Tags Inventory
// @Produce json
// @Param category path string true "Category (general, electronics, furniture, etc)"
// @Success 200 {object} []types.InventoryItem
// @Failure 400 {object} types.ErrorResponse
// @Router /inventory/category/{category} [get]
func (h *InventoryHandler) ListItemsByCategory(c *gin.Context) {
	category := c.Param("category")

	resp, err := h.Service.ListItemsByCategory(c.Request.Context(), category)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateItem godoc
// @Summary Update an inventory item
// @Description Modify fields of an existing inventory item
// @Tags Inventory
// @Accept json
// @Produce json
// @Param input body types.UpdateItemRequest true "Updated item info"
// @Success 200 {object} types.InventoryItem
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /inventory/ [put]
func (h *InventoryHandler) UpdateItem(c *gin.Context) {
	var req types.UpdateItemRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.NewErrorResponse(err))
		return
	}

	resp, err := h.Service.UpdateItem(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteItem godoc
// @Summary Delete an item
// @Description Remove inventory item by ID
// @Tags Inventory
// @Produce json
// @Param id path string true "Item ID"
// @Success 200 {object} []types.InventoryItem
// @Failure 500 {object} types.ErrorResponse
// @Router /inventory/{id} [delete]
func (h *InventoryHandler) DeleteItem(c *gin.Context) {
	id := c.Param("id")

	resp, err := h.Service.DeleteItem(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, resp)
}