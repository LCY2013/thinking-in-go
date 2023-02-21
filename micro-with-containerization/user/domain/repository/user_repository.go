package repository

import (
	"github.com/jinzhu/gorm"
	"user/domain/model"
)

type IUserRepository interface {
	// InitTable 初始化数据表
	InitTable() error
	// FindUserByName 根据用户名称查找用户信息
	FindUserByName(string) (*model.User, error)
	// FindUserByID 根据用户ID查找用户信息
	FindUserByID(int64) (*model.User, error)
	// CreateUser 创建用户
	CreateUser(*model.User) (int64, error)
	// DeleteUserByID 根据用户ID删除用户
	DeleteUserByID(int64) error
	// UpdateUser 更新用户信息
	UpdateUser(*model.User) error
	// FindAll 查找所有用户信息
	FindAll() ([]*model.User, error)
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db: db}
}

type UserRepository struct {
	db *gorm.DB
}

func (u UserRepository) InitTable() error {
	return u.db.CreateTable(&model.User{}).Error
}

func (u UserRepository) FindUserByName(userName string) (user *model.User, err error) {
	user = &model.User{}
	return user, u.db.Where("user_name = ?", userName).Find(user).Error
}

func (u UserRepository) FindUserByID(userID int64) (user *model.User, err error) {
	user = &model.User{}
	return user, u.db.Where("id = ?", userID).Error
}

func (u UserRepository) CreateUser(user *model.User) (int64, error) {
	return user.ID, u.db.Create(user).Error
}

func (u UserRepository) DeleteUserByID(userID int64) error {
	return u.db.Where("id = ?", userID).Delete(&model.User{}).Error
}

func (u UserRepository) UpdateUser(user *model.User) error {
	return u.db.Updates(user).Error
}

func (u UserRepository) FindAll() (userAll []*model.User, err error) {
	return userAll, u.db.Find(&userAll).Error
}
