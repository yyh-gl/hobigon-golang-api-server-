package repository

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/repository"
)

type blogRepository struct {}

// TODO: 場所ここ？
func NewBlogRepository() repository.BlogRepository {
	return &blogRepository{}
}

func (bp blogRepository) SelectByTitle(ctx context.Context, title string) (blog model.Blog, err error) {
	db := ctx.Value("db").(*gorm.DB)
	err = db.First(&blog, "title=?", title).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return model.Blog{}, nil
		}
		return model.Blog{}, err
	}
	return blog, nil
}
