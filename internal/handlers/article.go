package handlers

import (
	"strconv"

	"sv-app/backend/internal/models"
	"sv-app/backend/internal/validators"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ArticleHandler holds a DB reference so handlers stay stateless functions.
type ArticleHandler struct {
	DB *gorm.DB
}

// CreateArticle handles POST /article/
func (h *ArticleHandler) CreateArticle(c *fiber.Ctx) error {
	var req validators.ArticleRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON body"})
	}

	if err := validators.Validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": validators.FormatErrors(err)})
	}

	post := models.Post{
		Title:    req.Title,
		Content:  req.Content,
		Category: req.Category,
		Status:   req.Status,
	}

	if err := h.DB.Create(&post).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create article"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{})
}

// ListArticles handles GET /article/:limit/:offset
// Optional query param: ?status=publish|draft|thrash
func (h *ArticleHandler) ListArticles(c *fiber.Ctx) error {
	limit, err := strconv.Atoi(c.Params("limit"))

	if err != nil || limit < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "limit must be a non-negative integer"})
	}

	offset, err := strconv.Atoi(c.Params("offset"))

	if err != nil || offset < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "offset must be a non-negative integer"})
	}

	var posts []models.Post

	query := h.DB.Limit(limit).Offset(offset)

	// Optional status filter — makes the Dashboard tabs and Preview pagination work correctly.
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&posts).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch articles"})
	}

	return c.JSON(posts)
}

// GetArticle handles GET /article/:id
func (h *ArticleHandler) GetArticle(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var post models.Post

	if err := h.DB.First(&post, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "article not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch article"})
	}

	return c.JSON(post)
}

// UpdateArticle handles PUT /article/:id
func (h *ArticleHandler) UpdateArticle(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var req validators.ArticleRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON body"})
	}

	if err := validators.Validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": validators.FormatErrors(err)})
	}

	// Ensure record exists before updating.
	var post models.Post
	if err := h.DB.First(&post, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "article not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch article"})
	}

	post.Title = req.Title
	post.Content = req.Content
	post.Category = req.Category
	post.Status = req.Status

	if err := h.DB.Save(&post).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update article"})
	}

	return c.JSON(fiber.Map{})
}

// DeleteArticle handles DELETE /article/:id
func (h *ArticleHandler) DeleteArticle(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	result := h.DB.Delete(&models.Post{}, id)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete article"})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "article not found"})
	}

	return c.JSON(fiber.Map{})
}
