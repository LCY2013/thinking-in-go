package handler

import (
	"context"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/category/common"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/category/domain/model"
	"github.com/LCY2013/thinking-in-go/micro-with-containerization/category/domain/service"
	pb "github.com/LCY2013/thinking-in-go/micro-with-containerization/category/proto/category"
	"github.com/sirupsen/logrus"
)

type Category struct {
	categoryDataService service.ICategoryDataService
}

// CreateCategory 创建分类服务
func (c Category) CreateCategory(ctx context.Context,
	request *pb.CategoryRequest,
	response *pb.CreateCategoryResponse) error {
	category := &model.Category{}
	// 赋值
	err := common.SwapTo(request, category)
	if err != nil {
		return err
	}
	addCategoryId, err := c.categoryDataService.AddCategory(category)
	if err != nil {
		return err
	}

	response.Message = "添加分类成功"
	response.CategoryId = addCategoryId
	return nil
}

// UpdateCategory 更新分类服务
func (c Category) UpdateCategory(ctx context.Context,
	request *pb.CategoryRequest,
	response *pb.UpdateCategoryResponse) error {
	category := &model.Category{}
	// 赋值
	err := common.SwapTo(request, category)
	if err != nil {
		return err
	}
	err = c.categoryDataService.UpdateCategory(category)
	if err != nil {
		return err
	}

	response.Message = "更新分类成功"
	return nil
}

// DeleteCategory 删除分类信息
func (c Category) DeleteCategory(ctx context.Context,
	request *pb.DeleteCategoryRequest,
	response *pb.DeleteCategoryResponse) error {
	err := c.categoryDataService.DeleteCategory(request.CategoryId)
	if err != nil {
		return err
	}

	response.Message = "删除分类成功"
	return nil
}

// FindCategoryByName 根据名称查找分类信息
func (c Category) FindCategoryByName(ctx context.Context,
	request *pb.FindByNameRequest,
	response *pb.CategoryResponse) error {
	category, err := c.categoryDataService.FindCategoryByName(request.CategoryName)
	if err != nil {
		return err
	}

	return common.SwapTo(category, response)
}

// FindCategoryByID 根据id查询分类信息
func (c Category) FindCategoryByID(ctx context.Context,
	request *pb.FindByIdRequest,
	response *pb.CategoryResponse) error {
	category, err := c.categoryDataService.FindCategoryByID(request.CategoryId)
	if err != nil {
		return err
	}

	return common.SwapTo(category, response)
}

// FindCategoryByLevel 根据分类级别查询分类信息
func (c Category) FindCategoryByLevel(ctx context.Context,
	request *pb.FindByLevelRequest,
	response *pb.FindAllResponse) error {
	categorySlice, err := c.categoryDataService.FindCategoryByLevel(request.Level)
	if err != nil {
		return err
	}

	categoryToResponse(categorySlice, response)
	return nil
}

// categoryToResponse to response
func categoryToResponse(categorySlice []model.Category, response *pb.FindAllResponse) {
	for _, category := range categorySlice {
		cr := &pb.CategoryResponse{}
		err := common.SwapTo(category, cr)
		if err != nil {
			logrus.WithContext(context.TODO()).Error(err)
			break
		}
		response.Category = append(response.Category, cr)
	}
}

// FindCategoryByParent 根据分类父级查询子分类
func (c Category) FindCategoryByParent(ctx context.Context,
	request *pb.FindByParentRequest,
	response *pb.FindAllResponse) error {
	categorySlice, err := c.categoryDataService.FindCategoryByParent(request.ParentId)
	if err != nil {
		return err
	}

	categoryToResponse(categorySlice, response)
	return nil
}

// FindAllCategory 查询所有的分类信息
func (c Category) FindAllCategory(ctx context.Context,
	request *pb.FindAllRequest,
	response *pb.FindAllResponse) error {
	categorySlice, err := c.categoryDataService.FindAllCategory()
	if err != nil {
		return err
	}

	categoryToResponse(categorySlice, response)
	return nil
}

// New Return a new handler
func New(categoryDataService service.ICategoryDataService) *Category {
	return &Category{
		categoryDataService: categoryDataService,
	}
}
