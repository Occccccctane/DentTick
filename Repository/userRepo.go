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
}

type CachedUserRepository struct {
	dao Dao.UserDao
}

func (repo *CachedUserRepository) Create(ctx context.Context, u Domain.User) error {
	return repo.dao.Insert(ctx, repo.toEntity(u))
}
func NewUserRepository(dao Dao.UserDao) UserRepository {
	return &CachedUserRepository{
		dao: dao,
	}
}
func (repo *CachedUserRepository) toEntity(u Domain.User) Dao.User {
	return Dao.User{
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
