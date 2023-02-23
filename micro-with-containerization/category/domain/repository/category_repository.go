package repository

import (
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/category/domain/model"
	"github.com/jinzhu/gorm"
)

type ICategoryRepository interface {
	// InitTable 初始化数据表
	InitTable() error
	// FindCategoryByID 根据CategoryID查找Category信息
	FindCategoryByID(int64) (*model.Category, error)
	// CreateCategory 创建Category
	CreateCategory(*model.Category) (int64, error)
	// DeleteCategoryByID 根据CategoryID删除Category
	DeleteCategoryByID(int64) error
	// UpdateCategory 更新Category信息
	UpdateCategory(*model.Category) error
	// FindAll 查找所有Category信息
	FindAll() ([]model.Category, error)
	FindCategoryByName(string) (*model.Category, error)
	FindCategoryByLevel(uint32) ([]model.Category, error)
	FindCategoryByParent(int64) ([]model.Category, error)
}

func NewCategoryRepository(db *gorm.DB) ICategoryRepository {
	return &CategoryRepository{db: db}
}

type CategoryRepository struct {
	db *gorm.DB
}

func (rep *CategoryRepository) InitTable() error {
	return rep.db.CreateTable(&model.Category{}).Error
}

// FindCategoryByID 根据ID查找Category信息
func (rep *CategoryRepository) FindCategoryByID(categoryID int64) (category *model.Category, err error) {
	category = &model.Category{}
	return category, rep.db.First(category, categoryID).Error
}

// CreateCategory 创建Category信息
func (rep *CategoryRepository) CreateCategory(category *model.Category) (int64, error) {
	return category.ID, rep.db.Create(category).Error
}

// DeleteCategoryByID 根据ID删除Category信息
func (rep *CategoryRepository) DeleteCategoryByID(categoryID int64) error {
	return rep.db.Where("id = ?", categoryID).Delete(&model.Category{}).Error
}

// UpdateCategory 更新Category信息
func (rep *CategoryRepository) UpdateCategory(category *model.Category) error {
	return rep.db.Model(category).Update(category).Error
}

// FindAll 获取结果集
func (rep *CategoryRepository) FindAll() (categoryAll []model.Category, err error) {
	return categoryAll, rep.db.Find(&categoryAll).Error
}

// FindCategoryByName 根据分类名称进行查找
func (rep *CategoryRepository) FindCategoryByName(categoryName string) (category *model.Category, err error) {
	category = &model.Category{}
	return category, rep.db.Where("category_name = ?", categoryName).Find(category).Error
}

func (rep *CategoryRepository) FindCategoryByLevel(level uint32) (categorySlice []model.Category, err error) {
	return categorySlice, rep.db.Where("category_level = ?", level).Find(categorySlice).Error
}

func (rep *CategoryRepository) FindCategoryByParent(parent int64) (categorySlice []model.Category, err error) {
	return categorySlice, rep.db.Where("category_parent = ?", parent).Find(categorySlice).Error
}
