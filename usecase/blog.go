package usecase

import (
	"context"

	"github.com/pkg/errors"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/gateway"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model/entity"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/repository"
)

//////////////////////////////////////////////////
// NewBlogUseCase
//////////////////////////////////////////////////

// BlogUseCase : ブログ用のユースケースインターフェース
type BlogUseCase interface {
	Create(context.Context, string) (*entity.Blog, error)
	Show(context.Context, string) (*entity.Blog, error)
	Like(context.Context, string) (*entity.Blog, error)
}

type blogUseCase struct {
	br repository.BlogRepository
	sg gateway.SlackGateway
}

// NewBlogUseCase : ブログ用のユースケースを取得
func NewBlogUseCase(br repository.BlogRepository, sg gateway.SlackGateway) BlogUseCase {
	return &blogUseCase{
		br: br,
		sg: sg,
	}
}

//////////////////////////////////////////////////
// Create
//////////////////////////////////////////////////

// Create : ブログ情報を新規作成
func (bu blogUseCase) Create(ctx context.Context, title string) (*entity.Blog, error) {
	blog := entity.NewBlog(title)
	createdBlog, err := bu.br.Create(ctx, *blog)
	if err != nil {
		return nil, errors.Wrap(err, "blogRepository.Create()内でのエラー")
	}

	return createdBlog, nil
}

//////////////////////////////////////////////////
// Show
//////////////////////////////////////////////////

// Show : ブログ情報を1件取得
func (bu blogUseCase) Show(ctx context.Context, title string) (*entity.Blog, error) {
	blog, err := bu.br.SelectByTitle(ctx, title)
	if err != nil {
		switch err.Error() {
		case "record not found":
			return nil, err
		default:
			return nil, errors.Wrap(err, "blogRepository.SelectByTitle()内でのエラー")
		}
	}

	return blog, nil
}

//////////////////////////////////////////////////
// Like
//////////////////////////////////////////////////

// Like : 指定ブログにいいねをプラス1
func (bu blogUseCase) Like(ctx context.Context, title string) (*entity.Blog, error) {
	blog, err := bu.br.SelectByTitle(ctx, title)
	if err != nil {
		switch err.Error() {
		case "record not found":
			return nil, err
		default:
			return nil, errors.Wrap(err, "blogRepository.SelectByTitle()内でのエラー")
		}
	}

	// Count をプラス1
	blog.CountUp()
	blog, err = bu.br.Update(ctx, *blog)
	if err != nil {
		return nil, errors.Wrap(err, "blogRepository.Update()内でのエラー")
	}

	// Slack に通知
	err = bu.sg.SendLikeNotify(ctx, *blog)
	if err != nil {
		return nil, errors.Wrap(err, "slackGateway.SendLikeNotify()内でのエラー")
	}

	return blog, nil
}
