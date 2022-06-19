package test

import (
	"errors"
	"fufeng.org/test/mock/db/application/cmd"
	"fufeng.org/test/mock/db/application/query"
	"fufeng.org/test/mock/db/domain/user"
	"fufeng.org/test/mock/db/infrastructure/db/mock"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// 创建仓储的mock实例
	userRepository := mock.NewMockUserRepository(ctrl)
	// 模拟接口参数信息，必须有该步骤，不然mock不知道该咋操作，类似于程序断言
	userEntity := user.UserEntity{Id: 1, Name: "fufeng"}
	gomock.InOrder(
		// 这里主要是描述该mock是做啥，比如下面的含义是：期望  创建user(这个动作来自自定义接口，UserRepository) 返回值
		userRepository.EXPECT().Create(userEntity).Return(nil),
	)
	// 通过mock实例，创建业务示例
	userCmd := cmd.NewUserCmd(userRepository)
	// 业务操作创建用户
	if err := userCmd.Create(userEntity); err != nil {
		t.Error(err)
	}
}

func TestUpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// 创建仓储的mock实例
	userRepository := mock.NewMockUserRepository(ctrl)
	// 模拟接口参数信息，必须有该步骤，不然mock不知道该咋操作，类似于程序断言
	userEntity := user.UserEntity{Id: 1, Name: "fufeng"}
	gomock.InOrder(
		// 这里主要是描述该mock是做啥，比如下面的含义是：期望  更新user(这个动作来自自定义接口，UserRepository) 返回值
		userRepository.EXPECT().Update(userEntity).Return(nil),
	)
	// 通过mock实例，修改业务示例
	userCmd := cmd.NewUserCmd(userRepository)
	// 业务操作修改用户
	if err := userCmd.Update(userEntity); err != nil {
		t.Error(err)
	}
}

func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// 创建仓储的mock实例
	userRepository := mock.NewMockUserRepository(ctrl)
	// 模拟接口参数信息，必须有该步骤，不然mock不知道该咋操作，类似于程序断言
	gomock.InOrder(
		// 这里主要是描述该mock是做啥，比如下面的含义是：期望  删除user(这个动作来自自定义接口，UserRepository) 返回值
		// 模拟删除失败
		userRepository.EXPECT().Delete(1).Return(errors.New("delete fail")),
	)
	// 通过mock实例，删除业务示例
	userCmd := cmd.NewUserCmd(userRepository)
	// 业务操作删除用户
	if err := userCmd.Delete(1); err != nil {
		t.Error(err)
	}
}

func TestGetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// 创建仓储的mock实例
	userRepository := mock.NewMockUserRepository(ctrl)
	// 模拟接口参数信息，必须有该步骤，不然mock不知道该咋操作，类似于程序断言
	gomock.InOrder(
		// 这里主要是描述该mock是做啥，比如下面的含义是：期望  查询user(这个动作来自自定义接口，UserRepository) 返回值
		// 模拟返回某一个实体对象
		userRepository.EXPECT().Get(1).Return(user.UserEntity{
			Id:   1,
			Name: "fufeng",
		}),
	)
	// 通过mock实例，查询业务示例
	userQuery := query.NewUserQuery(userRepository)
	// 业务操作查询用户
	t.Log(userQuery.Get(1))
}
