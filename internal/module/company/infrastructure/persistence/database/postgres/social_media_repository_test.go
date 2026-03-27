package postgres

import (
	"context"
	"testing"
	"time"

	socialMediaDomain "github.com/paladignus/actajus/internal/module/social_media/domain"
	sharedpostgres "github.com/paladignus/actajus/internal/shared/infrastructure/persistence/database/postgres"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func TestSocialMediaRepositoryCreate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	now := time.Date(2026, 3, 27, 12, 0, 0, 0, time.UTC)
	item, err := socialMediaDomain.NewSocialMediaBuilder().
		WithIDCompany(55).
		WithPlatform("linkedin").
		WithURL("https://linkedin.com/company/acme").
		WithCreatedAt(now).
		WithUpdatedAt(now).
		Build()
	require.NoError(t, err)

	mock.ExpectQuery(`INSERT INTO social_media`).
		WithArgs(
			item.IDCompany(),
			item.Platform(),
			item.URL(),
			item.CreatedAt(),
			item.UpdatedAt(),
		).
		WillReturnRows(pgxmock.NewRows([]string{"idsocial_media"}).AddRow(int64(30)))

	repo := NewSocialMediaRepository(mock)
	err = repo.Create(ctx, item)
	require.NoError(t, err)
	require.Equal(t, int64(30), item.ID().Value())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSocialMediaRepositoryFindByIDCompany(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	now := time.Date(2026, 3, 27, 12, 0, 0, 0, time.UTC)
	mock.ExpectQuery(`SELECT\s+idsocial_media`).
		WithArgs(int64(55)).
		WillReturnRows(pgxmock.NewRows([]string{
			"idsocial_media", "platform", "url", "created_at", "updated_at",
		}).
			AddRow(int64(1), "linkedin", "https://linkedin.com/company/acme", now, now).
			AddRow(int64(2), "instagram", "https://instagram.com/acme", now, now))

	repo := NewSocialMediaRepository(mock)
	got, err := repo.FindByIDCompany(ctx, 55)
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.Equal(t, int64(1), got[0].ID().Value())
	require.Equal(t, int64(55), got[0].IDCompany().Value())
	require.Equal(t, "https://instagram.com/acme", got[1].URL().Value())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCompanyReadRepositoryList(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	now := time.Date(2026, 3, 27, 12, 0, 0, 0, time.UTC)
	companyID := int64(1)
	companyName := "Acme SA"
	tradeName := "Acme"
	cnpj := "12345678000199"
	addressID := int64(10)
	zip := "79000-000"
	title := "Matriz"
	street := "Rua A"
	complement := "Sala 1"
	reference := "Esquina"
	number := uint(42)
	neighborhood := "Centro"
	city := "Campo Grande"
	state := "MS"
	country := "BR"
	emailID := int64(20)
	emailAddress := "contato@acme.com"
	phoneID := int64(30)
	phoneNumber := "+5567999999999"
	phoneKind := "commercial"
	department := "sales"
	socialID := int64(40)
	platform := "linkedin"
	url := "https://linkedin.com/company/acme"
	registeredBy := "Jane Doe"
	rows := pgxmock.NewRows([]string{
		"idcompanies", "name", "trade_name", "cnpj", "created_at", "updated_at",
		"idaddresses", "zip", "title", "street", "complement", "reference", "number", "neighborhood",
		"city", "state", "country", "created_at", "updated_at",
		"idemails", "address", "created_at", "updated_at",
		"idphones", "number", "kind", "department", "created_at", "updated_at",
		"idsocial_media", "platform", "url", "created_at", "updated_at", "name",
	}).AddRow(
		&companyID, &companyName, &tradeName, &cnpj, &now, &now,
		&addressID, &zip, &title, &street, &complement, &reference, &number, &neighborhood,
		&city, &state, &country, &now, &now,
		&emailID, &emailAddress, &now, &now,
		&phoneID, &phoneNumber, &phoneKind, &department, &now, &now,
		&socialID, &platform, &url, &now, &now, &registeredBy,
	)

	mock.ExpectQuery(`SELECT\s+c.idcompanies`).
		WithArgs(11).
		WillReturnRows(rows)

	repo := NewCompanyReadRepository(mock)
	got, err := repo.List(ctx, nil, nil, 10, "http://localhost:8080/companies")
	require.NoError(t, err)
	require.Len(t, got.Data, 1)
	require.Equal(t, int64(1), got.Data[0].ID)
	require.Equal(t, "Jane Doe", got.Data[0].RegisteredByName)
	require.NotNil(t, got.Data[0].Address)
	require.Len(t, got.Data[0].Phones, 1)
	require.Len(t, got.Data[0].Emails, 1)
	require.Len(t, got.Data[0].SocialMedia, 1)
	require.False(t, got.PageInfo.HasNextPage)
	require.False(t, got.PageInfo.HasPreviousPage)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSocialMediaRepositoryUpdateAndDelete(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	now := time.Date(2026, 3, 27, 12, 0, 0, 0, time.UTC)
	deletedAt := now.Add(time.Minute)
	item, err := socialMediaDomain.NewSocialMediaBuilder().
		WithID(30).
		WithIDCompany(55).
		WithPlatform("linkedin").
		WithURL("https://linkedin.com/company/acme").
		WithUpdatedAt(now).
		WithDeletedAt(&deletedAt).
		Build()
	require.NoError(t, err)

	mock.ExpectExec(`UPDATE social_media SET platform =`).
		WithArgs(item.Platform(), item.URL(), item.UpdatedAt(), item.ID().Value()).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectExec(`UPDATE social_media SET updated_at =`).
		WithArgs(item.UpdatedAt(), item.DeletedAt(), item.ID().Value()).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectExec(`UPDATE social_media SET updated_at =`).
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), int64(55)).
		WillReturnResult(pgxmock.NewResult("UPDATE", 2))

	repo := NewSocialMediaRepository(mock)
	err = repo.Update(ctx, *item)
	require.NoError(t, err)
	err = repo.Delete(ctx, *item)
	require.NoError(t, err)
	err = repo.DeleteByIDCompany(ctx, 55)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCompanyReadRepositoryListWithAfterCursor(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	now := time.Date(2026, 3, 27, 12, 0, 0, 0, time.UTC)
	next := now.Add(time.Minute)
	after := sharedpostgres.EncodeCursor(now, "1")

	rows := companyListRows().
		AddRow(companyRow(next, int64(2), "Beta SA", "Beta", "22345678000199", "John Doe")...).
		AddRow(companyRow(next.Add(time.Minute), int64(3), "Gamma SA", "Gamma", "32345678000199", "John Doe")...)

	mock.ExpectQuery(`SELECT\s+c.idcompanies`).
		WithArgs(now.UTC(), "1", 2).
		WillReturnRows(rows)

	repo := NewCompanyReadRepository(mock)
	got, err := repo.List(ctx, &after, nil, 1, "http://localhost:8080/companies")
	require.NoError(t, err)
	require.Len(t, got.Data, 1)
	require.True(t, got.PageInfo.HasNextPage)
	require.True(t, got.PageInfo.HasPreviousPage)
	require.NotNil(t, got.PageInfo.NextURL)
	require.NotNil(t, got.PageInfo.PreviousURL)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCompanyReadRepositoryListWithBeforeCursor(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	now := time.Date(2026, 3, 27, 12, 0, 0, 0, time.UTC)
	before := sharedpostgres.EncodeCursor(now.Add(2*time.Minute), "3")

	rows := companyListRows().
		AddRow(companyRow(now.Add(time.Minute), int64(2), "Beta SA", "Beta", "22345678000199", "John Doe")...).
		AddRow(companyRow(now, int64(1), "Acme SA", "Acme", "12345678000199", "Jane Doe")...)

	mock.ExpectQuery(`SELECT\s+c.idcompanies`).
		WithArgs(now.Add(2*time.Minute).UTC(), "3", 2).
		WillReturnRows(rows)

	repo := NewCompanyReadRepository(mock)
	got, err := repo.List(ctx, nil, &before, 1, "http://localhost:8080/companies")
	require.NoError(t, err)
	require.Len(t, got.Data, 1)
	require.Equal(t, int64(2), got.Data[0].ID)
	require.True(t, got.PageInfo.HasNextPage)
	require.True(t, got.PageInfo.HasPreviousPage)
	require.NoError(t, mock.ExpectationsWereMet())
}

func companyListRows() *pgxmock.Rows {
	return pgxmock.NewRows([]string{
		"idcompanies", "name", "trade_name", "cnpj", "created_at", "updated_at",
		"idaddresses", "zip", "title", "street", "complement", "reference", "number", "neighborhood",
		"city", "state", "country", "created_at", "updated_at",
		"idemails", "address", "created_at", "updated_at",
		"idphones", "number", "kind", "department", "created_at", "updated_at",
		"idsocial_media", "platform", "url", "created_at", "updated_at", "name",
	})
}

func companyRow(now time.Time, companyID int64, companyName, tradeName, cnpj, registeredBy string) []any {
	addressID := int64(companyID * 10)
	zip := "79000-000"
	title := "Matriz"
	street := "Rua A"
	complement := "Sala 1"
	reference := "Esquina"
	number := uint(42)
	neighborhood := "Centro"
	city := "Campo Grande"
	state := "MS"
	country := "BR"
	emailID := int64(companyID * 20)
	emailAddress := "contato@example.com"
	phoneID := int64(companyID * 30)
	phoneNumber := "+5567999999999"
	phoneKind := "commercial"
	department := "sales"
	socialID := int64(companyID * 40)
	platform := "linkedin"
	url := "https://linkedin.com/company/example"

	return []any{
		&companyID, &companyName, &tradeName, &cnpj, &now, &now,
		&addressID, &zip, &title, &street, &complement, &reference, &number, &neighborhood,
		&city, &state, &country, &now, &now,
		&emailID, &emailAddress, &now, &now,
		&phoneID, &phoneNumber, &phoneKind, &department, &now, &now,
		&socialID, &platform, &url, &now, &now, &registeredBy,
	}
}
