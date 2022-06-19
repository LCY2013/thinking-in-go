package cmd

import "fufeng.org/test/mock/db/domain/user"

type UserCmd struct {
	userRepository user.UserRepository
}

func NewUserCmd(userRepository user.UserRepository) *UserCmd {
	return &UserCmd{
		userRepository: userRepository,
	}
}

func (user UserCmd) Create(userEntity user.UserEntity) error {
	// TODO 校验、业务逻辑
	return user.userRepository.Create(userEntity)
}

func (user UserCmd) Update(userEntity user.UserEntity) error {
	// TODO 校验、业务逻辑
	return user.userRepository.Update(userEntity)
}

func (user UserCmd) Delete(id int) error {
	// TODO 校验、业务逻辑
	return user.userRepository.Delete(id)
}
