package Repository

import (
	"DentTick/Domain"
	"DentTick/Repository/Dao"
	"context"
	"database/sql"
	"time"
)

var (
	ErrUserNotFound = Dao.ErrRecordNotFound
	ErrUserUnique   = Dao.PhoneUniqueErr
)

type UserRepository interface {
	Create(ctx context.Context, u Domain.User) error
	// FindById 查询用户
	FindById(ctx context.Context, id int64) (Domain.User, error)
	// UpdateProfile 更新用户资料
	UpdateProfile(ctx context.Context, u Domain.User) error
}

type CachedUserRepository struct {
	dao Dao.UserDao
}

func (repo *CachedUserRepository) Create(ctx context.Context, u Domain.User) error {
	return repo.dao.Insert(ctx, repo.toEntity(u))
}

func (repo *CachedUserRepository) FindById(ctx context.Context, id int64) (Domain.User, error) {
	// Dao -> Domain 映射
	u, err := repo.dao.FindById(ctx, id)
	return repo.toDomain(u), err
}

func (repo *CachedUserRepository) UpdateProfile(ctx context.Context, u Domain.User) error {
	// Domain -> Dao 映射
	return repo.dao.UpdateProfile(ctx, repo.toEntity(u))
}

func NewUserRepository(dao Dao.UserDao) UserRepository {
	return &CachedUserRepository{
		dao: dao,
	}
}
func (repo *CachedUserRepository) toEntity(u Domain.User) Dao.User {
	return Dao.User{
		Avatar:    u.Avatar,
		NickName:  u.NickName,
		Name:      u.Name,
		Info:      u.Info,
		Password:  u.Password,
		Identity:  u.Identity,
		DoctorId:  u.DoctorId,
		PatientId: u.PatientId,
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Utime: time.Now().UnixMilli(),
	}
}

func (repo *CachedUserRepository) toDomain(u Dao.User) Domain.User {
	// 数据库实体转领域对象
	return Domain.User{
		Id:        u.Id,
		Avatar:    u.Avatar,
		NickName:  u.NickName,
		Name:      u.Name,
		Info:      u.Info,
		Password:  u.Password,
		Identity:  u.Identity,
		DoctorId:  u.DoctorId,
		PatientId: u.PatientId,
		Phone:     u.Phone.String,
	}
}
