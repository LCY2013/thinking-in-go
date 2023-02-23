package service

import (
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/category/domain/model"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/category/domain/repository"
)

type ICategoryDataService interface {
	AddCategory(category *model.Category) (int64, error)
	DeleteCategory(categoryID int64) error
	UpdateCategory(category *model.Category) error
	FindCategoryByID(categoryID int64) (category *model.Category, err error)
	FindAllCategory() ([]model.Category, error)
	FindCategoryByName(string) (*model.Category, error)
	FindCategoryByLevel(uint32) ([]model.Category, error)
	FindCategoryByParent(int64) ([]model.Category, error)
}

// NewCategoryDataService 创建用户数据服务
func NewCategoryDataService(categoryRepository repository.ICategoryRepository) ICategoryDataService {
	return &CategoryService{
		categoryRepository: categoryRepository,
	}
}

type CategoryService struct {
	categoryRepository repository.ICategoryRepository
}

func (ser CategoryService) FindAllCategory() ([]model.Category, error) {
	return ser.categoryRepository.FindAll()
}

func (ser CategoryService) FindCategoryByName(s string) (*model.Category, error) {
	//TODO implement me
	panic("implement me")
}

func (ser CategoryService) FindCategoryByLevel(u uint32) ([]model.Category, error) {
	//TODO implement me
	panic("implement me")
}

func (ser CategoryService) FindCategoryByParent(i int64) ([]model.Category, error) {
	//TODO implement me
	panic("implement me")
}

func (ser CategoryService) FindCategoryByID(categoryID int64) (category *model.Category, err error) {
	category = &model.Category{}
	return ser.categoryRepository.FindCategoryByID(categoryID)
}

func (ser CategoryService) AddCategory(category *model.Category) (int64, error) {
	// MQ
	return ser.categoryRepository.CreateCategory(category)
}

func (ser CategoryService) DeleteCategory(categoryID int64) error {
	return ser.categoryRepository.DeleteCategoryByID(categoryID)
}

func (ser CategoryService) UpdateCategory(category *model.Category) error {
	// log
	return ser.categoryRepository.UpdateCategory(category)
}
