package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jacksonopp/go-recipe/domain"
	"github.com/jacksonopp/go-recipe/services"
	"gorm.io/gorm"
)

type TagHandler struct {
	r          fiber.Router
	db         *gorm.DB
	tagService services.TagService
}

func NewTagHandler(r fiber.Router, db *gorm.DB) *TagHandler {
	subpath := r.Group("/tag")
	tagService := services.NewTagService(db)
	return &TagHandler{r: subpath, db: db, tagService: tagService}
}

func (h *TagHandler) RegisterRoutes() {
	h.r.Get("/", h.GetTags)
	h.r.Post("/", h.CreateTag)
}

func (h *TagHandler) GetTags(c *fiber.Ctx) error {
	tags, err := h.tagService.GetAllTags()
	if err != nil {
		//TODO handle error codes
		return SendError(c, InternalServerError())
	}

	tagDtos := make([]domain.TagDto, 0, len(tags))
	for _, tag := range tags {
		tagDtos = append(tagDtos, tag.ToDto().(domain.TagDto))
	}

	return c.JSON(tagDtos)
}

func (h *TagHandler) CreateTag(c *fiber.Ctx) error {
	tag := struct {
		Tag string `json:"tag"`
	}{}
	err := c.BodyParser(&tag)
	if err != nil {
		return SendError(c, BadRequest("invalid request body"))
	}

	createdTag, err := h.tagService.CreateTag(tag.Tag)
	if err != nil {
		//TODO handle error codes
		return SendError(c, InternalServerError())
	}

	return c.JSON(createdTag.ToDto())
}
