package Dao

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	PhoneUniqueErr    = errors.New("邮箱唯一错误")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type UserDao interface {
	Insert(ctx context.Context, u User) error
}
type UserGormDao struct {
	db *gorm.DB
}

func (dao *UserGormDao) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now

	err := dao.db.WithContext(ctx).Create(&u).Error

	//将错误类型断言为unique键错误，并给出特定的错误处理
	var sqlErr *mysql.MySQLError
	if errors.As(err, &sqlErr) {
		const uniqueErrNum uint16 = 1062
		if sqlErr.Number == uniqueErrNum {
			return PhoneUniqueErr
		}
	}
	return err
}

func NewUserDao(db *gorm.DB) UserDao {
	return &UserGormDao{
		db: db,
	}
}

// User 映射数据库用户表
type User struct {
	Id       int64 `gorm:"primaryKey,autoFIncrement"`
	Name     string
	Info     string
	Password string

	//医生 Id和患者 Id互斥,身份由Identity字段决定
	//1为医生，2为患者,3为管理员
	Identity  int8
	DoctorId  string
	PatientId string
	Phone     sql.NullString `gorm:"unique"`

	//处理时区，统一用UTC 0时区的毫秒数
	Ctime int64
	Utime int64
}
