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
	FindByPhone(ctx context.Context, phone string) (Domain.User, error)
	GetById(ctx context.Context, id int64) (Domain.User, error)
}

type CachedUserRepository struct {
	dao Dao.UserDao
}

func (repo *CachedUserRepository) Create(ctx context.Context, u Domain.User) error {
	return repo.dao.Insert(ctx, repo.toEntity(u))
}

func (repo *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (Domain.User, error) {
	entity, err := repo.dao.SelectByPhone(ctx, phone)
	if err != nil {
		return Domain.User{}, err
	}
	return repo.toDomain(entity), nil
}

func (repo *CachedUserRepository) GetById(ctx context.Context, id int64) (Domain.User, error) {
	entity, err := repo.dao.SelectById(ctx, id)
	if err != nil {
		return Domain.User{}, err
	}
	return repo.toDomain(entity), nil
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

// convert Dao.User to Domain.User
func (repo *CachedUserRepository) toDomain(entity Dao.User) Domain.User {
	return Domain.User{
		Id:        entity.Id,
		Name:      entity.Name,
		Info:      entity.Info,
		Password:  entity.Password,
		Identity:  entity.Identity,
		DoctorId:  entity.DoctorId,
		PatientId: entity.PatientId,
		Phone:     entity.Phone.String,
	}
}
