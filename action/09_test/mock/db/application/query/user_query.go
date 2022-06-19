package query

import "fufeng.org/test/mock/db/domain/user"

type UserQuery struct {
	userRepository user.UserRepository
}

func NewUserQuery(userRepository user.UserRepository) *UserQuery {
	return &UserQuery{
		userRepository: userRepository,
	}
}

func (user UserQuery) Get(id int) user.UserEntity {
	// TODO 校验、业务逻辑
	return user.userRepository.Get(id)
}
