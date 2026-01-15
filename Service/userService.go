package Service

import (
	"DentTick/Domain"
	"DentTick/Repository"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserUnique            = Repository.ErrUserUnique
	ErrInvalidUserOrPassword = errors.New("账号或密码错误")
)

type UserService interface {
	Signup(ctx context.Context, u Domain.User) error
}
type userService struct {
	userRepo Repository.UserRepository
}

func (svc *userService) Signup(ctx context.Context, u Domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	//没报错就将密码加密为哈希，将哈希存入数据库中
	u.Password = string(hash)
	return svc.userRepo.Create(ctx, u)
}

func NewUserService(userRepo Repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}
