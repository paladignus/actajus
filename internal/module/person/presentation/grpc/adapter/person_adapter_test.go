package adapter_test

import (
	"testing"

	"github.com/paladignus/actajus/internal/module/person/presentation/grpc/adapter"
	personv1 "github.com/paladignus/actajus/proto/person/v1"
)

func TestProtoToCreatePersonCommand(t *testing.T) {
	t.Parallel()

	req := &personv1.CreatePersonRequest{
		Name:     "Jane Doe",
		Gender:   2,
		Birthday: "01/01/1990",
	}

	got := adapter.ProtoToCreatePersonCommand(req)

	if got.Name != req.Name || got.Gender != uint(req.Gender) || got.Birthday != req.Birthday {
		t.Fatalf("unexpected mapping: %+v", got)
	}
}
