// Package adapter
package adapter

import (
	"connectrpc.com/connect"
	"github.com/paladignus/actajus/internal/shared/domain"
	validationv1 "github.com/paladignus/actajus/proto/validation/v1"
)

func ValidationErrorDetail(err domain.ValidationError) *connect.Error {
	payload := ValidationPayload(err)
	detail := &validationv1.ValidationErrorDetail{
		Violations: make([]*validationv1.ValidationErrorDetail_Violation, 0, len(payload.Violations)),
	}
	for _, fe := range payload.Violations {
		detail.Violations = append(detail.Violations, &validationv1.ValidationErrorDetail_Violation{
			Path: fe.Path,
			Code: fe.Code,
			Meta: fe.Meta,
		})
	}
	cerr := connect.NewError(connect.CodeInvalidArgument, err)
	ed, e := connect.NewErrorDetail(detail)
	if e == nil {
		cerr.AddDetail(ed)
	}
	return cerr
}
