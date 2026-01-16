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
	ErrUserNotFound          = Repository.ErrUserNotFound
	ErrInvalidUserOrPassword = errors.New("手机号或密码错误")
)

type UserService interface {
	Signup(ctx context.Context, u Domain.User) error
	// GetProfile 获取用户资料
	GetProfile(ctx context.Context, id int64) (Domain.User, error)
	// EditProfile 编辑用户资料
	EditProfile(ctx context.Context, u Domain.User) error
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

func (svc *userService) GetProfile(ctx context.Context, id int64) (Domain.User, error) {
	// 直接透传到仓储层
	return svc.userRepo.FindById(ctx, id)
}

func (svc *userService) EditProfile(ctx context.Context, u Domain.User) error {
	// 只更新资料字段
	return svc.userRepo.UpdateProfile(ctx, u)
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
