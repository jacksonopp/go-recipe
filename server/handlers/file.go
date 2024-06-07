package handlers

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/jacksonopp/go-recipe/services"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
	"log"
	"strconv"
)

type FileHandler struct {
	db            *gorm.DB
	ctx           context.Context
	r             fiber.Router
	bucketService services.BucketService
}

func NewFileHandler(r fiber.Router, minio *minio.Client, db *gorm.DB) *FileHandler {
	ctx := context.Background()
	subpath := r.Group("/file")
	bucketService := services.NewBucketService(db, minio)
	return &FileHandler{db: db, ctx: ctx, r: subpath, bucketService: bucketService}
}

func (h *FileHandler) RegisterRoutes() {
	h.r.Post("/", AuthMiddleware(h.db), h.UploadFile)
	h.r.Get("/:id", h.GetFile)
}

// POST /file
func (h *FileHandler) UploadFile(c *fiber.Ctx) error {
	user, err := getUserFromLocals(c)
	if err != nil {
		log.Println("failed to get user from locals", err)
		return SendError(c, InternalServerError())
	}
	file, err := c.FormFile("file")
	if err != nil {
		log.Println("failed to get file", err)
		return SendError(c, BadRequest("invalid file"))
	}

	dbFile, err := h.bucketService.UploadFile(user.ID, file)
	if err != nil {
		log.Println("failed to upload file", err)
		return SendError(c, InternalServerError())
	}

	return c.JSON(dbFile.ToDto())
}

// GET /file/:id
func (h *FileHandler) GetFile(c *fiber.Ctx) error {
	fileID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Println("failed to parse id", err)
		return SendError(c, BadRequest("invalid id"))
	}
	file, err := h.bucketService.GetFileByID(uint(fileID))
	if err != nil {
		log.Println("failed to get file", err)
		if errors.Is(err, services.ErrFileNotFound) {
			return SendError(c, NotFound(map[string]string{"file": "file not found"}))
		}
		return SendError(c, InternalServerError())
	}
	return c.JSON(file.ToDto())
}
