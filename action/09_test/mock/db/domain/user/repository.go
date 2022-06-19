package user

type UserRepository interface {
	Create(UserEntity) error
	Update(UserEntity) error
	Delete(id int) error
	Get(id int) UserEntity
}
