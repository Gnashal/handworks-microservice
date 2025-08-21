package types

type ApiGatewayServices struct {
	AccountService      IAccountService
	AdminService        IAdminService
	BookingService      IBookingService
	EmployeeService     IEmployeeService
	InventoryService    IInventoryService
	NotificationService INotificationService
	PaymentService      IPaymentService
}
type IAccountService interface {
}
type IAdminService interface {
}
type IBookingService interface {
}
type IEmployeeService interface {
}
type IInventoryService interface {
}
type INotificationService interface {
}
type IPaymentService interface {
}
