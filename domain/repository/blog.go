package repository

import (
	"context"

	"github.com/yyh-gl/hobigon-golang-api-server/domain/model/blog"
)

// BlogRepository : ブログ用のリポジトリインターフェース
type BlogRepository interface {
	Create(ctx context.Context, blog blog.Blog) (*blog.Blog, error)
	SelectByTitle(ctx context.Context, title string) (*blog.Blog, error)
	Update(ctx context.Context, blog blog.Blog) (*blog.Blog, error)
}
