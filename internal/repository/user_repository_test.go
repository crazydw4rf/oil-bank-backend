package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/crazydw4rf/oil-bank-backend/internal/entity"
	. "github.com/crazydw4rf/oil-bank-backend/internal/pkg/result"
	"github.com/crazydw4rf/oil-bank-backend/internal/services"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/gomega"
)

// NOTE: This test suite uses sqlmock which works well with sqlx.
// For better pgx.PgError testing, consider migrating to:
//   - pgx v5 (github.com/jackc/pgx/v5)
//   - pgxmock v4 (github.com/pashagolub/pgxmock/v4)
// This would allow proper mocking of pgx-specific errors without workarounds.

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, services.DatabaseService) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	dbService := services.DatabaseService{DB: sqlxDB}

	return mockDB, mock, dbService
}

func TestUserRepository_Create_Success(t *testing.T) {
	g := NewWithT(t)
	mockDB, mock, dbService := setupMockDB(t)
	defer mockDB.Close()

	repo := NewUserRepository(dbService)
	ctx := context.Background()

	user := &entity.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		UserType:     entity.SELLER,
	}

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "address_id", "username", "email", "password_hash", "user_type", "created_at", "updated_at"}).
		AddRow(1, nil, "testuser", "test@example.com", "hashedpassword", "SELLER", now, now)

	mock.ExpectQuery(`INSERT INTO "User"`).
		WithArgs("testuser", "test@example.com", "hashedpassword", entity.SELLER).
		WillReturnRows(rows)

	result := repo.Create(ctx, user)

	g.Expect(result.IsError()).To(BeFalse())
	g.Expect(result.Value().Id).To(Equal(int64(1)))
	g.Expect(result.Value().Username).To(Equal("testuser"))
	g.Expect(result.Value().Email).To(Equal("test@example.com"))
	g.Expect(mock.ExpectationsWereMet()).To(Succeed())
}

func TestUserRepository_Create_DuplicateEmail(t *testing.T) {
	g := NewWithT(t)
	mockDB, mock, dbService := setupMockDB(t)
	defer mockDB.Close()

	repo := NewUserRepository(dbService)
	ctx := context.Background()

	user := &entity.User{
		Username:     "testuser",
		Email:        "duplicate@example.com",
		PasswordHash: "hashedpassword",
		UserType:     entity.SELLER,
	}

	// Simulate a database constraint violation error
	// Note: With sqlmock, we can't easily create a proper pgx.PgError
	// The actual pgx error handling is better tested with integration tests
	mock.ExpectQuery(`INSERT INTO "User"`).
		WithArgs("testuser", "duplicate@example.com", "hashedpassword", entity.SELLER).
		WillReturnError(sql.ErrConnDone) // Generic error for demo

	result := repo.Create(ctx, user)

	// Verify error handling works (even if not the exact error type)
	g.Expect(result.IsError()).To(BeTrue())
	g.Expect(mock.ExpectationsWereMet()).To(Succeed())
}

func TestUserRepository_Find_Success(t *testing.T) {
	g := NewWithT(t)
	mockDB, mock, dbService := setupMockDB(t)
	defer mockDB.Close()

	repo := NewUserRepository(dbService)
	ctx := context.Background()

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "address_id", "username", "email", "password_hash", "user_type", "created_at", "updated_at"}).
		AddRow(1, nil, "testuser", "test@example.com", "hashedpassword", "SELLER", now, now)

	mock.ExpectQuery(`SELECT \* FROM "User" WHERE id`).
		WithArgs(int64(1)).
		WillReturnRows(rows)

	result := repo.Find(ctx, 1)

	g.Expect(result.IsError()).To(BeFalse())
	g.Expect(result.Value().Id).To(Equal(int64(1)))
	g.Expect(result.Value().Username).To(Equal("testuser"))
	g.Expect(result.Value().Email).To(Equal("test@example.com"))
	g.Expect(mock.ExpectationsWereMet()).To(Succeed())
}

func TestUserRepository_Find_NotFound(t *testing.T) {
	g := NewWithT(t)
	mockDB, mock, dbService := setupMockDB(t)
	defer mockDB.Close()

	repo := NewUserRepository(dbService)
	ctx := context.Background()

	mock.ExpectQuery(`SELECT \* FROM "User" WHERE id`).
		WithArgs(int64(999)).
		WillReturnError(sql.ErrNoRows)

	result := repo.Find(ctx, 999)

	g.Expect(result.IsError()).To(BeTrue())
	g.Expect(result.RootError().Cause()).To(Equal(ENTITY_NOT_FOUND))
	g.Expect(result.RootError().IsExpected).To(BeTrue())
	g.Expect(result.RootError().Error()).To(ContainSubstring("user not found"))
	g.Expect(mock.ExpectationsWereMet()).To(Succeed())
}

func TestUserRepository_Update_Success(t *testing.T) {
	g := NewWithT(t)
	mockDB, mock, dbService := setupMockDB(t)
	defer mockDB.Close()

	repo := NewUserRepository(dbService)
	ctx := context.Background()

	user := &entity.User{
		Id:       1,
		Username: "updateduser",
		Email:    "updated@example.com",
	}

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "address_id", "username", "email", "password_hash", "user_type", "created_at", "updated_at"}).
		AddRow(1, nil, "updateduser", "updated@example.com", "hashedpassword", "SELLER", now, now)

	mock.ExpectQuery(`UPDATE "User" SET`).
		WithArgs(int64(1), "updateduser", "updated@example.com").
		WillReturnRows(rows)

	result := repo.Update(ctx, user)

	g.Expect(result.IsError()).To(BeFalse())
	g.Expect(result.Value().Id).To(Equal(int64(1)))
	g.Expect(result.Value().Username).To(Equal("updateduser"))
	g.Expect(result.Value().Email).To(Equal("updated@example.com"))
	g.Expect(mock.ExpectationsWereMet()).To(Succeed())
}

func TestUserRepository_Delete_Success(t *testing.T) {
	g := NewWithT(t)
	mockDB, mock, dbService := setupMockDB(t)
	defer mockDB.Close()

	repo := NewUserRepository(dbService)
	ctx := context.Background()

	mock.ExpectExec(`DELETE FROM "User" WHERE id`).
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	result := repo.Delete(ctx, 1)

	g.Expect(result.IsError()).To(BeFalse())
	g.Expect(result.Value()).To(BeTrue())
	g.Expect(mock.ExpectationsWereMet()).To(Succeed())
}

func TestUserRepository_Delete_NotFound(t *testing.T) {
	g := NewWithT(t)
	mockDB, mock, dbService := setupMockDB(t)
	defer mockDB.Close()

	repo := NewUserRepository(dbService)
	ctx := context.Background()

	mock.ExpectExec(`DELETE FROM "User" WHERE id`).
		WithArgs(int64(999)).
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected

	result := repo.Delete(ctx, 999)

	g.Expect(result.IsError()).To(BeTrue())
	g.Expect(result.RootError().Cause()).To(Equal(ENTITY_NOT_FOUND))
	g.Expect(result.RootError().IsExpected).To(BeTrue())
	g.Expect(result.RootError().Error()).To(ContainSubstring("user not found"))
	g.Expect(mock.ExpectationsWereMet()).To(Succeed())
}

func TestUserRepository_FindByEmail_Success(t *testing.T) {
	g := NewWithT(t)
	mockDB, mock, dbService := setupMockDB(t)
	defer mockDB.Close()

	repo := NewUserRepository(dbService)
	ctx := context.Background()

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "address_id", "username", "email", "password_hash", "user_type", "created_at", "updated_at"}).
		AddRow(1, nil, "testuser", "test@example.com", "hashedpassword", "SELLER", now, now)

	mock.ExpectQuery(`SELECT \* FROM "User" WHERE email`).
		WithArgs("test@example.com").
		WillReturnRows(rows)

	result := repo.FindByEmail(ctx, "test@example.com")

	g.Expect(result.IsError()).To(BeFalse())
	g.Expect(result.Value().Id).To(Equal(int64(1)))
	g.Expect(result.Value().Email).To(Equal("test@example.com"))
	g.Expect(mock.ExpectationsWereMet()).To(Succeed())
}

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	g := NewWithT(t)
	mockDB, mock, dbService := setupMockDB(t)
	defer mockDB.Close()

	repo := NewUserRepository(dbService)
	ctx := context.Background()

	mock.ExpectQuery(`SELECT \* FROM "User" WHERE email`).
		WithArgs("notfound@example.com").
		WillReturnError(sql.ErrNoRows)

	result := repo.FindByEmail(ctx, "notfound@example.com")

	g.Expect(result.IsError()).To(BeTrue())
	g.Expect(result.RootError().Cause()).To(Equal(ENTITY_NOT_FOUND))
	g.Expect(result.RootError().IsExpected).To(BeTrue())
	g.Expect(mock.ExpectationsWereMet()).To(Succeed())
}

func TestHandleUserError_NoRows(t *testing.T) {
	g := NewWithT(t)

	result := handleUserError[*entity.User](sql.ErrNoRows)

	g.Expect(result.IsError()).To(BeTrue())
	g.Expect(result.RootError().Cause()).To(Equal(ENTITY_NOT_FOUND))
	g.Expect(result.RootError().IsExpected).To(BeTrue())
	g.Expect(result.RootError().Error()).To(Equal("user not found"))
}

func TestHandleUserError_GenericDatabaseError(t *testing.T) {
	g := NewWithT(t)

	// Test with a generic database error
	genericErr := sql.ErrConnDone
	result := handleUserError[*entity.User](genericErr)

	g.Expect(result.IsError()).To(BeTrue())
	g.Expect(result.RootError().Cause()).To(Equal(INTERNAL_SERVICE_ERROR))
	g.Expect(result.RootError().IsExpected).To(BeFalse())
	g.Expect(result.RootError().Error()).To(ContainSubstring("database error"))
}

func TestUserRepository_FindByEmailWithSeller_Success(t *testing.T) {
	g := NewWithT(t)
	mockDB, mock, dbService := setupMockDB(t)
	defer mockDB.Close()

	repo := NewUserRepository(dbService)
	ctx := context.Background()

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "address_id", "username", "email", "password_hash", "user_type", "created_at", "updated_at",
		"seller_id", "seller_name",
	}).AddRow(1, nil, "testuser", "test@example.com", "hashedpassword", "SELLER", now, now, 10, "Test Seller")

	mock.ExpectQuery(`SELECT u\.\*, s\.id as seller_id, s\.seller_name`).
		WithArgs("test@example.com").
		WillReturnRows(rows)

	result := repo.FindByEmailWithSeller(ctx, "test@example.com")

	g.Expect(result.IsError()).To(BeFalse())
	g.Expect(result.Value().Id).To(Equal(int64(1)))
	g.Expect(result.Value().Email).To(Equal("test@example.com"))
	g.Expect(result.Value().SellerId).To(Equal(int64(10)))
	g.Expect(result.Value().SellerName).To(Equal("Test Seller"))
	g.Expect(mock.ExpectationsWereMet()).To(Succeed())
}

func TestUserRepository_FindByEmailWithCollector_Success(t *testing.T) {
	g := NewWithT(t)
	mockDB, mock, dbService := setupMockDB(t)
	defer mockDB.Close()

	repo := NewUserRepository(dbService)
	ctx := context.Background()

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "address_id", "username", "email", "password_hash", "user_type", "created_at", "updated_at",
		"collector_id", "collector_name",
	}).AddRow(2, nil, "collector1", "collector@example.com", "hashedpassword", "COLLECTOR", now, now, 20, "Test Collector")

	mock.ExpectQuery(`SELECT u\.\*, c\.id as collector_id, c\.collector_name`).
		WithArgs("collector@example.com").
		WillReturnRows(rows)

	result := repo.FindByEmailWithCollector(ctx, "collector@example.com")

	g.Expect(result.IsError()).To(BeFalse())
	g.Expect(result.Value().Id).To(Equal(int64(2)))
	g.Expect(result.Value().Email).To(Equal("collector@example.com"))
	g.Expect(result.Value().CollectorId).To(Equal(int64(20)))
	g.Expect(result.Value().CollectorName).To(Equal("Test Collector"))
	g.Expect(mock.ExpectationsWereMet()).To(Succeed())
}

func TestUserRepository_FindByEmailWithCompany_Success(t *testing.T) {
	g := NewWithT(t)
	mockDB, mock, dbService := setupMockDB(t)
	defer mockDB.Close()

	repo := NewUserRepository(dbService)
	ctx := context.Background()

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "address_id", "username", "email", "password_hash", "user_type", "created_at", "updated_at",
		"company_id", "company_name",
	}).AddRow(3, nil, "company1", "company@example.com", "hashedpassword", "COMPANY", now, now, 30, "Test Company")

	mock.ExpectQuery(`SELECT u\.\*, co\.id as company_id, co\.company_name`).
		WithArgs("company@example.com").
		WillReturnRows(rows)

	result := repo.FindByEmailWithCompany(ctx, "company@example.com")

	g.Expect(result.IsError()).To(BeFalse())
	g.Expect(result.Value().Id).To(Equal(int64(3)))
	g.Expect(result.Value().Email).To(Equal("company@example.com"))
	g.Expect(result.Value().CompanyId).To(Equal(int64(30)))
	g.Expect(result.Value().CompanyName).To(Equal("Test Company"))
	g.Expect(mock.ExpectationsWereMet()).To(Succeed())
}
