package error_handle

import "errors"

type TenantService struct {
	err error
}

func NewTenantService() *TenantService {
	return new(TenantService)
}

func (s *TenantService) Err() error {
	return s.err
}

type Tenant struct {
	TenantId string
}

func (s *TenantService) ActiveTenant(userId int64) *Tenant {
	if s.err != nil {
		return nil
	}

	// todo

	return &Tenant{TenantId: "1"}
}

func (s *TenantService) InitBilling(tenant *Tenant) bool {
	if s.err != nil {
		return false
	}

	// todo

	s.err = errors.New("new error")
	return true
}
