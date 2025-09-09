package service

import "handworks/common/grpc/account"

func DetermineStatusType(status string) string {
	if status == account.EmployeeStatus_ACTIVE.String() {
		return account.EmployeeStatus_ACTIVE.String()
	} else if status == account.EmployeeStatus_INACTIVE.String() {
		return account.EmployeeStatus_INACTIVE.String()
	}
	return account.EmployeeStatus_ONDUTY.String()
}
