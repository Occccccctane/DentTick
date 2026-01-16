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
	ErrInvalidUserOrPassword = errors.New("手机号或密码错误")
)

type UserService interface {
	Signup(ctx context.Context, u Domain.User) error
	Login(ctx context.Context, phone string, password string) (Domain.User, error)
	GetUser(ctx context.Context, uid int64) (Domain.User, error)
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

func (svc *userService) Login(ctx context.Context, phone string, password string) (Domain.User, error) {
	u, err := svc.userRepo.FindByPhone(ctx, phone)
	if err != nil {
		return Domain.User{}, ErrInvalidUserOrPassword
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return Domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (svc *userService) GetUser(ctx context.Context, uid int64) (Domain.User, error) {
	return svc.userRepo.GetById(ctx, uid)
}

func NewUserService(userRepo Repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}
