//go:build all || integration
// +build all integration

package handler

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "github.com/proullon/ramsql/driver"

	"lab1/internal/repository"
	"lab1/internal/usecase"
)

func Test_handler_CreateUser_Integration(t *testing.T) {
	t.Parallel()

	type fields struct {
		createUserUsecase    func(repo usecase.CreateUserRepository) CreateUserUsecase
		createUserRepository func(db *sql.DB) usecase.CreateUserRepository
		db                   func() *sql.DB
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		expect int
	}{
		{
			name: "1. Successful test.",
			fields: fields{
				db: func() *sql.DB {
					migration := `CREATE TABLE users (email text NOT NULL UNIQUE, name text NOT NULL);`

					db, err := sql.Open("ramsql", "1")
					if err != nil {
						t.Fatalf("can't open db: %v", err.Error())
					}

					_, err = db.Exec(migration)
					if err != nil {
						t.Errorf("can't run migrations: %v", err.Error())
					}
					return db
				},
				createUserRepository: func(db *sql.DB) usecase.CreateUserRepository {
					return repository.New(db)
				},
				createUserUsecase: func(repo usecase.CreateUserRepository) CreateUserUsecase {
					return usecase.New(repo)
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					"",
					"/test",
					strings.NewReader(`{"name":"vasya","email":"vasya@mail.ru"}`),
				),
			},
			expect: 200,
		},
		{
			name: "2. Successful test, user already exist.",
			fields: fields{
				db: func() *sql.DB {
					migration := `CREATE TABLE users (email text NOT NULL UNIQUE, name text NOT NULL);`

					db, err := sql.Open("ramsql", "2")
					if err != nil {
						t.Fatalf("can't open db: %v", err.Error())
					}

					_, err = db.Exec(migration)
					if err != nil {
						t.Errorf("can't run migrations: %v", err.Error())
					}

					query := `INSERT INTO users (name, email) VALUES ('vasya', 'vasya@mail.ru');`
					_, err = db.Exec(query)
					if err != nil {
						t.Errorf("can't insert into users: %v", err.Error())
					}

					return db
				},
				createUserRepository: func(db *sql.DB) usecase.CreateUserRepository {
					return repository.New(db)
				},
				createUserUsecase: func(repo usecase.CreateUserRepository) CreateUserUsecase {
					return usecase.New(repo)
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					"",
					"/test",
					strings.NewReader(`{"name":"vasya","email":"vasya@mail.ru"}`),
				),
			},
			expect: 200,
		},
		{
			name: "3. Broken input.",
			fields: fields{
				db: func() *sql.DB {
					return &sql.DB{}
				},
				createUserRepository: func(db *sql.DB) usecase.CreateUserRepository {
					return repository.New(db)
				},
				createUserUsecase: func(repo usecase.CreateUserRepository) CreateUserUsecase {
					return usecase.New(repo)
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					"",
					"/test",
					strings.NewReader(`not json`),
				),
			},
			expect: 400,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := New(tt.fields.createUserUsecase(tt.fields.createUserRepository(tt.fields.db())))
			h.CreateUser(tt.args.w, tt.args.r)

			if tt.expect != tt.args.w.Code {
				t.Errorf("expected != got, expected = %v, got = %v", tt.expect, tt.args.w.Code)
			}
		})
	}
}
