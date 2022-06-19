## 操作流程

> go get github.com/golang/mock/gomock

> go get github.com/golang/mock/mockgen

### 1. 创建存储层接口 repository

```text
type UserEntity struct {
	Id   int
	Name string
}
```

```text
type UserRepository interface {
	Create(UserEntity) error
	Update(UserEntity) error
	Delete(id int) error
	Get(id int) UserEntity
}
```

### 2. mock文件生成

进入action/09_test

> mockgen -source=./mock/db/domain/user/repository.go -destination=./mock/db/infrastructure/db/mock/mock_repository.go -package=mock


### 3. 创建 CQRS（读写分离） 目录，并且创建对应 user_cmd.go user_query.go

```text
type UserCmd struct {
	userRepository user.UserRepository
}
```

```text
type UserQuery struct {
	userRepository user.UserRepository
}
```

实现对应的业务逻辑

### 4. 编写对应的测试代码 user_test.go 开启mock测试用例

### 5. 测试
进入测试用例目录：
go test .

### 7. 测试结果查看
生成测试覆盖率的 profile 文件：
> go test -coverprofile=cover.out .

利用 profile 文件生成可视化界面
> go tool cover -html=cover.out




