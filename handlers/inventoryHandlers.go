package handlers

import (
	"context"
	"handworks-api/types"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateItem godoc
// @Summary Create a new inventory item
// @Description Adds a new item to inventory
// @Security BearerAuth
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := h.Service.CreateItem(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetItem godoc
// @Summary Get an item by ID
// @Description Retrieve a single inventory item
// @Security BearerAuth
// @Tags Inventory
// @Produce json
// @Param id path string true "Item ID"
// @Success 200 {object} types.InventoryItem
// @Failure 404 {object} types.ErrorResponse
// @Router /inventory/{id} [get]
func (h *InventoryHandler) GetItem(c *gin.Context) {
	id := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := h.Service.GetItem(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, types.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, resp)
}
// GetItem godoc
// @Summary Get all items
// @Description Retrieve all inventory items
// @Security BearerAuth
// @Tags Inventory
// @Produce json
// @Success 200 {object} []types.InventoryItem
// @Failure 404 {object} types.ErrorResponse
// @Router /inventory/ [get]
func (h *InventoryHandler) GetItems(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := h.Service.GetItems(ctx)
	if err != nil {
		c.JSON(http.StatusNotFound, types.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, resp)
}
// ListItemsByType godoc
// @Summary List items filtered by type
// @Description Get all inventory items matching the given type
// @Security BearerAuth
// @Tags Inventory
// @Produce json
// @Param type path string true "Item type (resource, equipment)"
// @Success 200 {object} []types.InventoryItem
// @Failure 400 {object} types.ErrorResponse
// @Router /inventory/type/{type} [get]
func (h *InventoryHandler) ListItemsByType(c *gin.Context) {
	itemType := c.Param("type")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := h.Service.ListItemsByType(ctx, itemType)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListItemsByStatus godoc
// @Summary List items filtered by status
// @Description Get all inventory items with the given stock status
// @Security BearerAuth
// @Tags Inventory
// @Produce json
// @Param status path string true "Status (high, low, danger, out_of_stock)"
// @Success 200 {object} []types.InventoryItem
// @Failure 400 {object} types.ErrorResponse
// @Router /inventory/status/{status} [get]
func (h *InventoryHandler) ListItemsByStatus(c *gin.Context) {
	status := c.Param("status")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := h.Service.ListItemsByStatus(ctx, status)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListItemsByCategory godoc
// @Summary List items filtered by category
// @Description Retrieve all inventory items in a category
// @Security BearerAuth
// @Tags Inventory
// @Produce json
// @Param category path string true "Category (general, electronics, furniture, etc)"
// @Success 200 {object} []types.InventoryItem
// @Failure 400 {object} types.ErrorResponse
// @Router /inventory/category/{category} [get]
func (h *InventoryHandler) ListItemsByCategory(c *gin.Context) {
	category := c.Param("category")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := h.Service.ListItemsByCategory(ctx, category)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateItem godoc
// @Summary Update an inventory item
// @Description Modify fields of an existing inventory item
// @Security BearerAuth
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
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := h.Service.UpdateItem(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteItem godoc
// @Summary Delete an item
// @Description Remove inventory item by ID
// @Security BearerAuth
// @Tags Inventory
// @Produce json
// @Param id path string true "Item ID"
// @Success 200 {object} []types.InventoryItem
// @Failure 500 {object} types.ErrorResponse
// @Router /inventory/{id} [delete]
func (h *InventoryHandler) DeleteItem(c *gin.Context) {
	id := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := h.Service.DeleteItem(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, resp)
}