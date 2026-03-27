package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/paladignus/actajus/internal/module/person/domain"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func TestPersonCreate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	now := time.Date(2026, 3, 27, 12, 0, 0, 0, time.UTC)
	person, err := domain.NewPersonBuilder().
		WithFirstName("Jane").
		WithLastName("Doe").
		WithGender(2).
		WithBirthday("01/01/1990").
		WithCreatedAt(now).
		WithUpdatedAt(now).
		Build()
	require.NoError(t, err)

	mock.ExpectQuery(`INSERT INTO people`).
		WithArgs(
			person.FirstName(),
			person.LastName(),
			person.IDGender(),
			person.Birthday().Value(),
			person.CreatedAt(),
			person.UpdatedAt(),
		).
		WillReturnRows(pgxmock.NewRows([]string{"idpeople"}).AddRow(uint(7)))

	repo := NewPerson(mock)
	err = repo.Create(ctx, person)
	require.NoError(t, err)
	require.Equal(t, uint(7), person.ID())
	require.NoError(t, mock.ExpectationsWereMet())
}
