package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateItem godoc
// @Summary Create a new inventory item
// @Description Adds a new item to inventory
// @Tags Inventory
// @Accept json
// @Produce json
// @Param input body map[string]interface{} true "Item info"
// @Success 200 {object} map[string]string
// @Router /inventory/item [post]
func (h *InventoryHandler) CreateItem(c *gin.Context) {
	_ = h.Service.CreateItem(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// GetItem godoc
// @Summary Get an inventory item
// @Description Retrieve an item by its ID
// @Tags Inventory
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Success 200 {object} map[string]interface{}
// @Router /inventory/item/{id} [get]
func (h *InventoryHandler) GetItem(c *gin.Context) {
	_ = h.Service.GetItem(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// ListAllItems godoc
// @Summary List all inventory items
// @Description Returns all items in inventory
// @Tags Inventory
// @Accept json
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /inventory/items [get]
func (h *InventoryHandler) ListAllItems(c *gin.Context) {
	_ = h.Service.ListAllItems(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// ListItemsByType godoc
// @Summary List inventory items by type
// @Description Returns items filtered by type
// @Tags Inventory
// @Accept json
// @Produce json
// @Param type path string true "Item type"
// @Success 200 {array} map[string]interface{}
// @Router /inventory/items/type/{type} [get]
func (h *InventoryHandler) ListItemsByType(c *gin.Context) {
	_ = h.Service.ListItemsByType(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// ListItemsByStatus godoc
// @Summary List inventory items by status
// @Description Returns items filtered by status
// @Tags Inventory
// @Accept json
// @Produce json
// @Param status path string true "Item status"
// @Success 200 {array} map[string]interface{}
// @Router /inventory/items/status/{status} [get]
func (h *InventoryHandler) ListItemsByStatus(c *gin.Context) {
	_ = h.Service.ListItemsByStatus(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// ListItemsByCategory godoc
// @Summary List inventory items by category
// @Description Returns items filtered by category
// @Tags Inventory
// @Accept json
// @Produce json
// @Param category path string true "Item category"
// @Success 200 {array} map[string]interface{}
// @Router /inventory/items/category/{category} [get]
func (h *InventoryHandler) ListItemsByCategory(c *gin.Context) {
	_ = h.Service.ListItemsByCategory(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// UpdateItem godoc
// @Summary Update an inventory item
// @Description Update item information
// @Tags Inventory
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Param input body map[string]interface{} true "Updated item info"
// @Success 200 {object} map[string]string
// @Router /inventory/item/{id} [put]
func (h *InventoryHandler) UpdateItem(c *gin.Context) {
	_ = h.Service.UpdateItem(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// DeleteItem godoc
// @Summary Delete an inventory item
// @Description Removes an item from inventory
// @Tags Inventory
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Success 200 {object} map[string]string
// @Router /inventory/item/{id} [delete]
func (h *InventoryHandler) DeleteItem(c *gin.Context) {
	_ = h.Service.DeleteItem(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}