package usecase_test

import (
	"errors"
	"testing"

	"multi-stage-build/domain/model"
	mock "multi-stage-build/domain/repository/mock"
	"multi-stage-build/usecase"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

type mocks struct {
	userRepository *mock.MockUserRepository
}

func newWithMocks(t *testing.T) (usecase.UserUsecase, *mocks) {
	ctrl := gomock.NewController(t)
	userRepository := mock.NewMockUserRepository(ctrl)
	return usecase.NewUserUsecase(userRepository), &mocks{
		userRepository: userRepository,
	}
}

func Test_userUsecase_Create(t *testing.T) {
	type args struct {
		name string
	}
	type expected struct {
		user *model.User
		err  error
	}

	for name, tt := range map[string]struct {
		args     args
		prepare  func(f *mocks)
		expected expected
	}{
		"【正常系】ユーザの新規登録": {
			args: args{name: "sample"},
			prepare: func(f *mocks) {
				f.userRepository.EXPECT().Create(&model.User{Name: "sample"}).Return(&model.User{Name: "sample"}, nil).Times(1)
			},
			expected: expected{user: &model.User{ID: 0, Name: "sample"}, err: nil},
		},
		"【異常系】Nameが空の場合、ユーザの新規登録でエラーが発生する": {
			args:     args{name: ""},
			expected: expected{user: nil, err: errors.New("名前は必須です。")},
		},
	} {
		t.Run(name, func(t *testing.T) {
			u, m := newWithMocks(t)
			if tt.prepare != nil {
				tt.prepare(m)
			}
			got, err := u.Create(tt.args.name)
			assert.Equal(t, tt.expected.user, got)
			assert.Equal(t, tt.expected.err, err)
		})
	}
}
