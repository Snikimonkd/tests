//go:build all || unit
// +build all unit

package repository

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"lab1/internal/model"
)

func Test_repository_CheckUserExist(t *testing.T) {
	t.Parallel()

	type fields struct {
		db func() *sql.DB
	}
	type args struct {
		email string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "1. Successful test, user does not exist.",
			fields: fields{
				db: func() *sql.DB {
					db, mock, err := sqlmock.New()
					if err != nil {
						t.Errorf("can't mock db: %v", err.Error())
					}

					mock.
						ExpectQuery(`SELECT email FROM users WHERE email = `).
						WithArgs("email@email.ru").
						WillReturnRows(sqlmock.NewRows([]string{}))

					return db
				},
			},
			args: args{
				email: "email@email.ru",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "2. Successful test, user not exist.",
			fields: fields{
				db: func() *sql.DB {
					db, mock, err := sqlmock.New()
					if err != nil {
						t.Errorf("can't mock db: %v", err.Error())
					}

					mock.
						ExpectQuery(`SELECT email FROM users WHERE email = `).
						WithArgs("email@email.ru").
						WillReturnRows(sqlmock.NewRows(
							[]string{"email"},
						).AddRow(
							"email",
						))

					return db
				},
			},
			args: args{
				email: "email@email.ru",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "3. Database internal error.",
			fields: fields{
				db: func() *sql.DB {
					db, mock, err := sqlmock.New()
					if err != nil {
						t.Errorf("can't mock db: %v", err.Error())
					}

					mock.
						ExpectQuery(`SELECT email FROM users WHERE email = $1;`).
						WithArgs("email@email.ru").
						WillReturnError(errors.New("some internal error"))

					return db
				},
			},
			args: args{
				email: "email@email.ru",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := repository{
				db: tt.fields.db(),
			}

			got, err := r.CheckUserExist(tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("repository.CheckUserExist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("repository.CheckUserExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_repository_CreateUser(t *testing.T) {
	t.Parallel()

	type fields struct {
		db func() *sql.DB
	}
	type args struct {
		user model.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "1. Successful test.",
			fields: fields{
				db: func() *sql.DB {
					db, mock, err := sqlmock.New()
					if err != nil {
						t.Errorf("can't mock db: %v", err.Error())
					}

					mock.
						ExpectExec(`INSERT INTO users`).
						WillReturnResult(sqlmock.NewResult(1, 1))

					return db
				},
			},
			args: args{
				user: model.User{
					Name:  "name",
					Email: "email",
				},
			},
			wantErr: false,
		},
		{
			name: "2. Database error.",
			fields: fields{
				db: func() *sql.DB {
					db, mock, err := sqlmock.New()
					if err != nil {
						t.Errorf("can't mock db: %v", err.Error())
					}

					mock.
						ExpectExec(`INSERT INTO users`).
						WillReturnError(errors.New("some internal error"))

					return db
				},
			},
			args: args{
				user: model.User{
					Name:  "name",
					Email: "email",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := repository{
				db: tt.fields.db(),
			}
			if err := r.CreateUser(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("repository.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
