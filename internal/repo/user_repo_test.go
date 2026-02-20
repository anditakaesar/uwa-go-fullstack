package repo_test

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/anditakaesar/uwa-go-fullstack/internal/common"
	"github.com/anditakaesar/uwa-go-fullstack/internal/domain"
	"github.com/anditakaesar/uwa-go-fullstack/internal/repo"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

type mockItems struct {
	ctx    context.Context
	mockDB pgxmock.PgxPoolIface
	now    time.Time
}

func setupMocks() (*mockItems, error) {
	mockDB, err := pgxmock.NewPool()
	if err != nil {
		return nil, err
	}

	return &mockItems{
		ctx:    context.Background(),
		mockDB: mockDB,
		now:    time.Now(),
	}, nil
}

func TestUserRepository_GetExecutor(test *testing.T) {
	test.Parallel()

	test.Run("success return from context", func(t *testing.T) {
		m, err := setupMocks()
		assert.NoError(t, err)
		defer m.mockDB.Close()

		m.mockDB.ExpectBegin()

		newTx, err := m.mockDB.Begin(m.ctx)
		assert.NoError(t, err)

		ctxWithValue := context.WithValue(m.ctx, common.TxKey, newTx)
		r := repo.NewUserRepository(m.mockDB)

		got := r.GetExecutor(ctxWithValue)
		assert.Equal(t, newTx, got)
	})

	test.Run("success return default", func(t *testing.T) {
		m, err := setupMocks()
		assert.NoError(t, err)
		defer m.mockDB.Close()

		r := repo.NewUserRepository(m.mockDB)
		got := r.GetExecutor(m.ctx)
		assert.Equal(t, m.mockDB, got)
	})
}

func TestUserRepository_CreateUser(test *testing.T) {
	test.Parallel()

	const query = `
			INSERT INTO users (username, password)
			VALUES ($1, $2)
			RETURNING id, username, created_at, updated_at, deleted_at;
		`
	newUser := domain.User{
		Username: "user1",
		Password: "password1",
	}

	test.Run("success", func(t *testing.T) {
		m, err := setupMocks()
		assert.NoError(t, err)
		defer m.mockDB.Close()

		rows := m.mockDB.NewRows([]string{"id", "username", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), newUser.Username, m.now, nil, nil)
		m.mockDB.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(newUser.Username, newUser.Password).
			WillReturnRows(rows)

		r := repo.NewUserRepository(m.mockDB)
		res, err := r.CreateUser(m.ctx, newUser)

		assert.NoError(t, err)
		assert.Equal(t, "user1", res.Username)
		assert.Equal(t, res.CreatedAt, m.now)
		assert.NoError(t, m.mockDB.ExpectationsWereMet())
	})

	test.Run("error", func(t *testing.T) {
		m, err := setupMocks()
		assert.NoError(t, err)
		defer m.mockDB.Close()

		m.mockDB.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(newUser.Username, newUser.Password).
			WillReturnError(errors.New("query_error"))

		r := repo.NewUserRepository(m.mockDB)
		res, err := r.CreateUser(m.ctx, newUser)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.NoError(t, m.mockDB.ExpectationsWereMet())
	})

}

func TestUserRepository_CreateAdminUser(test *testing.T) {
	test.Parallel()

	const query = `
			INSERT INTO users (username, password, role)
			VALUES ($1, $2, $3)
			RETURNING id, username, role, created_at, updated_at, deleted_at;
		`
	newUser := domain.User{
		Username: "user1",
		Password: "password1",
	}

	test.Run("success", func(t *testing.T) {
		m, err := setupMocks()
		assert.NoError(t, err)
		defer m.mockDB.Close()

		rows := m.mockDB.NewRows([]string{"id", "username", "role", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), newUser.Username, domain.RoleAdmin, m.now, nil, nil)
		m.mockDB.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(newUser.Username, newUser.Password, domain.RoleAdmin).
			WillReturnRows(rows)

		r := repo.NewUserRepository(m.mockDB)
		res, err := r.CreateUserAdmin(m.ctx, newUser)

		assert.NoError(t, err)
		assert.Equal(t, "user1", res.Username)
		assert.Equal(t, domain.RoleAdmin, res.Role)
		assert.Equal(t, res.CreatedAt, m.now)
		assert.NoError(t, m.mockDB.ExpectationsWereMet())
	})

	test.Run("error", func(t *testing.T) {
		m, err := setupMocks()
		assert.NoError(t, err)
		defer m.mockDB.Close()

		m.mockDB.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(newUser.Username, newUser.Password, domain.RoleAdmin).
			WillReturnError(errors.New("query_error"))

		r := repo.NewUserRepository(m.mockDB)
		res, err := r.CreateUserAdmin(m.ctx, newUser)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.NoError(t, m.mockDB.ExpectationsWereMet())
	})

}

func TestUserRepository_GetUser(test *testing.T) {
	test.Parallel()

	const query = `
		SELECT id, username, password, role, created_at, updated_at, deleted_at
        FROM users
        WHERE deleted_at IS NULL
		AND username = $1
	`

	expectUser := domain.User{
		Username: "user1",
		Password: "password1",
		Role:     domain.RoleUser,
	}

	test.Run("success", func(t *testing.T) {
		m, err := setupMocks()
		assert.NoError(t, err)
		defer m.mockDB.Close()

		rows := m.mockDB.NewRows([]string{"id", "username", "password", "role", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), expectUser.Username, expectUser.Password, expectUser.Role, m.now, nil, nil)
		m.mockDB.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(expectUser.Username).
			WillReturnRows(rows)

		r := repo.NewUserRepository(m.mockDB)
		res, err := r.GetUser(m.ctx, "user1")

		assert.NoError(t, err)
		assert.Equal(t, "user1", res.Username)
		assert.Equal(t, domain.RoleUser, res.Role)
		assert.Equal(t, res.CreatedAt, m.now)
		assert.NoError(t, m.mockDB.ExpectationsWereMet())
	})

	test.Run("error", func(t *testing.T) {
		m, err := setupMocks()
		assert.NoError(t, err)
		defer m.mockDB.Close()

		m.mockDB.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(expectUser.Username).
			WillReturnError(errors.New("query_error"))

		r := repo.NewUserRepository(m.mockDB)
		res, err := r.GetUser(m.ctx, "user1")

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.NoError(t, m.mockDB.ExpectationsWereMet())
	})
}

func TestUserRepository_GetUserByID(test *testing.T) {
	test.Parallel()

	const query = `
		SELECT id, username, role, created_at, updated_at, deleted_at
        FROM users
        WHERE deleted_at IS NULL
		AND id = $1
	`

	expectUser := domain.User{
		Username: "user1",
		Role:     domain.RoleUser,
	}

	test.Run("success", func(t *testing.T) {
		m, err := setupMocks()
		assert.NoError(t, err)
		defer m.mockDB.Close()

		rows := m.mockDB.NewRows([]string{"id", "username", "role", "created_at", "updated_at", "deleted_at"}).
			AddRow(int64(1), expectUser.Username, expectUser.Role, m.now, nil, nil)
		m.mockDB.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(int64(1)).
			WillReturnRows(rows)

		r := repo.NewUserRepository(m.mockDB)
		res, err := r.GetUserByID(m.ctx, int64(1))

		assert.NoError(t, err)
		assert.Equal(t, "user1", res.Username)
		assert.Equal(t, domain.RoleUser, res.Role)
		assert.Equal(t, res.CreatedAt, m.now)
		assert.NoError(t, m.mockDB.ExpectationsWereMet())
	})

	test.Run("error", func(t *testing.T) {
		m, err := setupMocks()
		assert.NoError(t, err)
		defer m.mockDB.Close()

		m.mockDB.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(int64(1)).
			WillReturnError(errors.New("query_error"))

		r := repo.NewUserRepository(m.mockDB)
		res, err := r.GetUserByID(m.ctx, int64(1))

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.NoError(t, m.mockDB.ExpectationsWereMet())
	})
}
