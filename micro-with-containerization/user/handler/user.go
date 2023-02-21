package handler

import (
	"context"
	"user/domain/model"
	"user/domain/service"
	"user/proto/user"
)

type User struct {
	userService service.IUserDataService
}

func (u User) Register(ctx context.Context,
	request *user.UserRegisterRequest,
	response *user.UserRegisterResponse) error {
	userRegister := &model.User{
		UserName:     request.UserName,
		FirstName:    request.FirstName,
		HashPassword: request.Pwd,
	}

	_, err := u.userService.AddUser(userRegister)
	if err != nil {
		return err
	}

	response.Message = "注册成功"
	return nil
}

func (u User) Login(ctx context.Context,
	request *user.UserLoginRequest,
	response *user.UserLoginResponse) error {
	ok, err := u.userService.CheckPwd(request.GetUserName(), request.GetPwd())
	if err != nil {
		return err
	}
	response.IsSuccess = ok
	return nil
}

func (u User) GetUserInfo(ctx context.Context,
	request *user.UserInfoRequest,
	response *user.UserInfoResponse) error {
	user, err := u.userService.FindUserByName(request.GetUserName())
	if err != nil {
		return err
	}

	response = UserForResponse(user)
	return nil
}

// UserForResponse 类型转化
func UserForResponse(userModel *model.User) *user.UserInfoResponse {
	return &user.UserInfoResponse{
		UserId:    userModel.ID,
		UserName:  userModel.UserName,
		FirstName: userModel.FirstName,
	}
}

// New Return a new handler
func New(userService service.IUserDataService) *User {
	return &User{
		userService: userService,
	}
}
