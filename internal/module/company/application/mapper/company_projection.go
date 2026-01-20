// Package mapper
package mapper

import (
	addrDomain "github.com/paladignus/actajus/internal/module/address/domain"
	"github.com/paladignus/actajus/internal/module/company/domain"
)

type CompanyProjection struct {
	Company *domain.Company
	Address *addrDomain.Address
}
