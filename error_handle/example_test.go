package error_handle

import (
	"log"
	"testing"
)

func Test_SimpleErrorHandler(t *testing.T) {
	tenantService := NewTenantService()
	tenant := tenantService.ActiveTenant(1)
	tenantService.InitBilling(tenant)

	if err := tenantService.Err(); err != nil {
		log.Printf("error: %s", err.Error())
	}
}
