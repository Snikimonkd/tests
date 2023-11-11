//go:build all || unit
// +build all unit

package usecase

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"

	"lab1/internal/model"
	"lab1/internal/usecase/mocks"
)

func Test_usecase_CreateUser(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	type fields struct {
		createUserRepository func(user model.User) CreateUserRepository
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
				createUserRepository: func(user model.User) CreateUserRepository {
					m := mocks.NewMockCreateUserRepository(ctrl)

					m.EXPECT().
						CheckUserExist(user.Email).
						Return(false, nil)

					m.EXPECT().
						CreateUser(user).
						Return(nil)

					return m
				},
			},
			args: args{
				user: model.User{
					Name:  "vasya",
					Email: "email@mail.ru",
				},
			},
			wantErr: false,
		},
		{
			name: "2. CreateUser() error.",
			fields: fields{
				createUserRepository: func(user model.User) CreateUserRepository {
					m := mocks.NewMockCreateUserRepository(ctrl)

					m.EXPECT().
						CheckUserExist(user.Email).
						Return(false, nil)

					m.EXPECT().
						CreateUser(user).
						Return(errors.New("some error"))

					return m
				},
			},
			args: args{
				user: model.User{
					Name:  "vasya",
					Email: "email@mail.ru",
				},
			},
			wantErr: true,
		},
		{
			name: "3. User already exist.",
			fields: fields{
				createUserRepository: func(user model.User) CreateUserRepository {
					m := mocks.NewMockCreateUserRepository(ctrl)

					m.EXPECT().
						CheckUserExist(user.Email).
						Return(true, nil)

					return m
				},
			},
			args: args{
				user: model.User{
					Name:  "vasya",
					Email: "email@mail.ru",
				},
			},
			wantErr: false,
		},
		{
			name: "4. CheckUserExist() error.",
			fields: fields{
				createUserRepository: func(user model.User) CreateUserRepository {
					m := mocks.NewMockCreateUserRepository(ctrl)

					m.EXPECT().
						CheckUserExist(user.Email).
						Return(true, errors.New("some error"))

					return m
				},
			},
			args: args{
				user: model.User{
					Name:  "vasya",
					Email: "email@mail.ru",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			u := usecase{
				createUserRepository: tt.fields.createUserRepository(tt.args.user),
			}
			if err := u.CreateUser(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("usecase.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
