//go:build all || unit
// +build all unit

package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	_ "github.com/proullon/ramsql/driver"

	"lab1/internal/handler/mocks"
	"lab1/internal/model"
)

func Test_handler_CreateUser(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	type fields struct {
		createUserUsecase func() CreateUserUsecase
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
				createUserUsecase: func() CreateUserUsecase {
					m := mocks.NewMockCreateUserUsecase(ctrl)
					m.EXPECT().
						CreateUser(model.User{
							Name:  "vasya",
							Email: "vasya@mail.ru",
						}).
						Return(nil)
					return m
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
			name: "2. Internal error.",
			fields: fields{
				createUserUsecase: func() CreateUserUsecase {
					m := mocks.NewMockCreateUserUsecase(ctrl)
					m.EXPECT().
						CreateUser(model.User{
							Name:  "vasya",
							Email: "vasya@mail.ru",
						}).
						Return(errors.New("some internal error"))
					return m
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
			expect: 500,
		},
		{
			name: "3. Unmarshal json error.",
			fields: fields{
				createUserUsecase: func() CreateUserUsecase {
					return mocks.NewMockCreateUserUsecase(ctrl)
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					"",
					"/test",
					strings.NewReader(`not a json`),
				),
			},
			expect: 400,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := handler{
				createUserUsecase: tt.fields.createUserUsecase(),
			}
			h.CreateUser(tt.args.w, tt.args.r)

			if tt.expect != tt.args.w.Code {
				t.Errorf("expected != got, expected = %v, got = %v", tt.expect, tt.args.w.Code)
			}
		})
	}
}
