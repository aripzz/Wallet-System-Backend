package handlers

import (
	"Wallet-System-Backend/infra"
	"Wallet-System-Backend/internal/entity"
	"Wallet-System-Backend/internal/helpers"
	"Wallet-System-Backend/internal/repository"
	"Wallet-System-Backend/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ProductHandler struct {
	productUsecase usecase.ProductsUsecase
}

func NewProductHandler(app *fiber.App, db *gorm.DB, cache *infra.RedisClient) {
	repo := repository.NewProductRepository(db)
	usecase := usecase.NewProductsUsecase(repo, cache)
	handler := &ProductHandler{productUsecase: usecase}

	apiv1 := app.Group("/api/v1")

	apiv1.Post("/products", handler.Create)
	apiv1.Get("/products", handler.GetAll)
	apiv1.Get("/products/:id", handler.GetByID)
	apiv1.Put("/products/:id", handler.Update)
	apiv1.Delete("/products/:id", handler.Delete)
}

// @Summary Get all products
// @Description Get a list of all products
// @Tags products
// @Accept  json
// @Produce  json
// @Success 200 {object} entity.Products
// @Security BearerAuth
// @Router /api/v1/products [get]
func (h *ProductHandler) GetAll(c *fiber.Ctx) error {
	products, err := h.productUsecase.GetAll()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return helpers.SendResponse(c, fiber.StatusOK, "successfully retrieved products", products)
}

// @Summary Get a product by ID
// @Description Get a single product by ID
// @Tags products
// @Accept  json
// @Produce  json
// @Param id path int true "Product ID"
// @Success 200 {object} entity.Products
// @Security BearerAuth
// @Router /api/v1/products/{id} [get]
func (h *ProductHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID")
	}
	product, err := h.productUsecase.GetByID(id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return helpers.SendResponse(c, fiber.StatusOK, "successfully retrieved product", product)
}

// @Summary Create a new product
// @Description Add a new product to the database
// @Tags products
// @Accept  json
// @Produce  json
// @Param product body entity.CreateProducts true "Product data"
// @Success 201 {string} string "Product created"
// @Security BearerAuth
// @Router /api/v1/products [post]
func (h *ProductHandler) Create(c *fiber.Ctx) error {
	var product entity.CreateProducts
	if err := c.BodyParser(&product); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := h.productUsecase.Create(product); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return helpers.SendResponse(c, fiber.StatusCreated, "product created successfully", nil)
}

// @Summary Update an existing product
// @Description Update a product's data
// @Tags products
// @Accept  json
// @Produce  json
// @Param id path int true "Product ID"
// @Param product body entity.Products true "Product data"
// @Success 200 {string} string "Product updated"
// @Security BearerAuth
// @Router /api/v1/products/{id} [put]
func (h *ProductHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID")
	}

	var product entity.UpdateProducts
	if err := c.BodyParser(&product); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	product.ID = id
	if err := h.productUsecase.Update(product); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return helpers.SendResponse(c, fiber.StatusOK, "product updated successfully", nil)
}

// @Summary Update an existing product
// @Description Update a product's data
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body entity.UpdateProducts true "Product data"
// @Success 200 {object} helpers.StandardResponse "Product updated successfully"
// @Failure 400 {object} helpers.StandardResponse "Invalid ID or request body"
// @Failure 404 {object} helpers.StandardResponse "Data not found"
// @Security BearerAuth
// @Router /api/v1/products/{id} [patch]
func (h *ProductHandler) UpdatePatch(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID")
	}

	var req entity.UpdateProducts
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	req.ID = id
	err = h.productUsecase.Update(req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return helpers.SendResponse(c, fiber.StatusOK, "Product updated successfully", nil)
}

// @Summary Delete a product
// @Description Delete a product by ID
// @Tags products
// @Accept  json
// @Produce  json
// @Param id path int true "Product ID"
// @Success 200 {string} string "Product deleted"
// @Security BearerAuth
// @Router /api/v1/products/{id} [delete]
func (h *ProductHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID")
	}
	if err := h.productUsecase.Delete(id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return helpers.SendResponse(c, fiber.StatusOK, "product deleted successfully", nil)
}
