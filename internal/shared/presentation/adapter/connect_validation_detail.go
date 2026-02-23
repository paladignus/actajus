// Package adapter
package adapter

import (
	"connectrpc.com/connect"
	"github.com/paladignus/actajus/internal/shared/domain"
	validationv1 "github.com/paladignus/actajus/proto/validation/v1"
)

func ValidationErrorDetail(err domain.ValidationError) *connect.Error {
	detail := &validationv1.ValidationErrorDetail{
		Violations: make([]*validationv1.ValidationErrorDetail_Violation, 0, len(err.Violations)),
	}
	for _, fe := range err.Violations {
		meta := fe.Meta
		if meta == nil {
			meta = map[string]string{}
		}
		detail.Violations = append(detail.Violations, &validationv1.ValidationErrorDetail_Violation{
			Path: fe.Path,
			Code: string(fe.Code),
			Meta: meta,
		})
	}
	cerr := connect.NewError(connect.CodeInvalidArgument, err)
	ed, e := connect.NewErrorDetail(detail)
	if e == nil {
		cerr.AddDetail(ed)
	}
	return cerr
}
