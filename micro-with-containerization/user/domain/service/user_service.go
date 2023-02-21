package service

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"user/domain/model"
	"user/domain/repository"
)

type IUserDataService interface {
	AddUser(user *model.User) (int64, error)
	DeleteUser(userID int64) error
	UpdateUser(user *model.User, isChangedPwd bool) error
	FindUserByName(userName string) (*model.User, error)
	CheckPwd(userName string, pwd string) (isOk bool, err error)
}

// NewUserDataService 创建用户数据服务
func NewUserDataService(userRepository repository.IUserRepository) IUserDataService {
	return &UserService{
		userRepository: userRepository,
	}
}

type UserService struct {
	userRepository repository.IUserRepository
}

// GeneratePassword 加密用户密码
func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

// ValidatePassword 校验用户密码
func ValidatePassword(hashPassword string, userPassword string) (isOk bool, err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(userPassword)); err != nil {
		return false, errors.WithMessage(err, "密码校验失败")
	}
	return true, nil
}

func (u UserService) AddUser(user *model.User) (int64, error) {
	pwdByte, err := GeneratePassword(user.HashPassword)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	user.HashPassword = string(pwdByte)

	// MQ
	return u.userRepository.CreateUser(user)
}

func (u UserService) DeleteUser(userID int64) error {
	return u.userRepository.DeleteUserByID(userID)
}

func (u UserService) UpdateUser(user *model.User, isChangedPwd bool) error {
	// 判断是否更新了密码
	if isChangedPwd {
		pwdByte, err := GeneratePassword(user.HashPassword)
		if err != nil {
			return errors.WithStack(err)
		}

		user.HashPassword = string(pwdByte)
	}

	// log
	return u.userRepository.UpdateUser(user)
}

func (u UserService) FindUserByName(userName string) (*model.User, error) {
	return u.userRepository.FindUserByName(userName)
}

func (u UserService) CheckPwd(userName string, pwd string) (isOk bool, err error) {
	name, err := u.userRepository.FindUserByName(userName)
	if err != nil {
		return false, errors.WithStack(err)
	}
	return ValidatePassword(name.HashPassword, pwd)
}
