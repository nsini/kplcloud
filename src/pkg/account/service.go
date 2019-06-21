package account

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

var ErrInvalidArgument = errors.New("invalid argument")

type Service interface {
	Detail(ctx context.Context, id int64) (rs map[string]interface{}, err error)

}

type service struct {
	logger log.Logger
	config config.Config
}

/**
 * @Title 详情页
 */
func (c *service) Detail(ctx context.Context, id int64) (rs map[string]interface{}, err error) {
	detail, err := c.post.Find(id)
	if err != nil {
		return
	}

	if detail == nil {
		return nil, repository.PostNotFound
	}

	go func() {
		if err = c.post.SetReadNum(detail); err != nil {
			_ = c.logger.Log("post.SetReadNum", err.Error())
		}
	}()

	var headerImage string

	if image, err := c.image.FindByPostIdLast(id); err == nil && image != nil {
		headerImage = imageUrl(image.RealPath.String, c.config.Get("image-domain"))
	}

	return map[string]interface{}{
		"content":      detail.Content,
		"title":        detail.Title,
		"publish_at":   detail.PushTime.Time.Format("2006/01/02 15:04:05"),
		"updated_at":   detail.UpdatedAt,
		"author":       detail.User.Username,
		"comment":      detail.Reviews,
		"banner_image": headerImage,
	}, nil
}

/**
 * @Title 列表页
 */
func (c *service) List(ctx context.Context, order, by string, pageSize, offset int) (rs []map[string]interface{}, count uint64, err error) {
	// 取列表 判断搜索、分类、Tag条件
	// 取最多阅读

	posts, count, err := c.post.FindBy(order, by, pageSize, offset)
	if err != nil {
		return
	}

	var postIds []uint

	for _, post := range posts {
		postIds = append(postIds, post.Model.ID)
	}

	images, err := c.image.FindByPostIds(postIds)
	if err == nil && images == nil {
		_ = c.logger.Log("c.image.FindByPostIds", "postIds", "err", err)
	}

	imageMap := make(map[int64]string, len(images))
	for _, image := range images {
		imageMap[image.PostID] = imageUrl(image.RealPath.String, c.config.Get("image-domain"))
	}

	_ = c.logger.Log("count", count)

	for _, val := range posts {
		imageUrl, ok := imageMap[int64(val.Model.ID)]
		if !ok {
			_ = c.logger.Log("postId", val.Model.ID, "image", ok)
		}
		rs = append(rs, map[string]interface{}{
			"id":         strconv.FormatUint(uint64(val.Model.ID), 10),
			"title":      val.Title,
			"desc":       val.Description,
			"publish_at": val.PushTime.Time.Format("2006/01/02 15:04:05"),
			"image_url":  imageUrl,
			"comment":    val.Reviews,
			"author":     val.User.Username,
		})
	}

	return
}


func NewService(logger log.Logger, cf config.Config, post repository.PostRepository, user repository.UserRepository, image repository.ImageRepository) Service {
	return &service{
		post:   post,
		user:   user,
		image:  image,
		logger: logger,
		config: cf,
	}
}
