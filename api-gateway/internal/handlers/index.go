package handlers

import (
	"handworks/api-gateway/internal/types"
	"handworks/common/utils"
)

func NewGatewayServices(log *utils.Logger) *types.ApiGatewayServices {
	return &types.ApiGatewayServices{
		AccountService:      nil, // Initialize with actual service implementation
		AdminService:        nil, // Initialize with actual service implementation
		BookingService:      nil, // Initialize with actual service implementation
		EmployeeService:     nil, // Initialize with actual service implementation
		InventoryService:    nil, // Initialize with actual service implementation
		NotificationService: nil, // Initialize with actual service implementation
		PaymentService:      nil, // Initialize with actual service implementation

	}
}
