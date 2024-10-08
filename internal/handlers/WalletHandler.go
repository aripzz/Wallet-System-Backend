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

type WalletHandler struct {
	walletUsecase usecase.WalletsUsecase
}

func NewWalletHandler(app *fiber.App, db *gorm.DB, cache *infra.RedisClient) {
	repo := repository.NewWalletRepository(db)
	usecase := usecase.NewWalletsUsecase(repo, cache) // Assuming you have a Redis client, pass it here
	handler := &WalletHandler{walletUsecase: usecase}

	apiv1 := app.Group("/api/v1")

	apiv1.Post("/wallets", handler.CreateWallet)
	apiv1.Get("/wallets", handler.GetAllWallets)
	apiv1.Get("/wallets/:id", handler.GetWalletByID)
	apiv1.Put("/wallets/:id", handler.UpdateWallet)
	apiv1.Delete("/wallets/:id", handler.DeleteWallet)
}

// @Summary Create a new wallet
// @Description Create a new wallet
// @Tags wallets
// @Accept json
// @Produce json
// @Param wallet body entity.CreateWallets true "Wallet data"
// @Success 201 {object} entity.Wallets
// @Security BearerAuth
// @Router /api/v1/wallets [post]
func (h *WalletHandler) CreateWallet(c *fiber.Ctx) error {
	var wallet entity.CreateWallets
	if err := c.BodyParser(&wallet); err != nil {
		return helpers.SendResponse(c, fiber.StatusBadRequest, "Invalid input", nil)
	}

	if err := h.walletUsecase.Create(wallet); err != nil {
		return helpers.SendResponse(c, fiber.StatusInternalServerError, "Failed to create wallet", nil)
	}

	return helpers.SendResponse(c, fiber.StatusCreated, "Wallet created successfully", nil)
}

// @Summary Get all wallets
// @Description Get a list of all wallets
// @Tags wallets
// @Accept json
// @Produce json
// @Success 200 {array} entity.Wallets
// @Security BearerAuth
// @Router /api/v1/wallets [get]
func (h *WalletHandler) GetAllWallets(c *fiber.Ctx) error {
	wallets, err := h.walletUsecase.GetAll()
	if err != nil {
		return helpers.SendResponse(c, fiber.StatusInternalServerError, "Failed to retrieve wallets", nil)
	}

	return helpers.SendResponse(c, fiber.StatusOK, "Wallets retrieved successfully", wallets)
}

// @Summary Get wallet by ID
// @Description Get a wallet by its ID
// @Tags wallets
// @Accept json
// @Produce json
// @Param id path int true "Wallet ID"
// @Success 200 {object} entity.Wallets
// @Security BearerAuth
// @Router /api/v1/wallets/{id} [get]
func (h *WalletHandler) GetWalletByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return helpers.SendResponse(c, fiber.StatusBadRequest, "Invalid ID format", nil)
	}
	wallet, err := h.walletUsecase.GetByID(id)
	if err != nil {
		return helpers.SendResponse(c, fiber.StatusNotFound, "Wallet not found", nil)
	}

	return helpers.SendResponse(c, fiber.StatusOK, "Wallet retrieved successfully", wallet)
}

// @Summary Update a wallet
// @Description Update a wallet by its ID
// @Tags wallets
// @Accept json
// @Produce json
// @Param id path int true "Wallet ID"
// @Param wallet body entity.UpdateWallets true "Updated wallet data"
// @Success 200 {object} entity.Wallets
// @Security BearerAuth
// @Router /api/v1/wallets/{id} [put]
func (h *WalletHandler) UpdateWallet(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return helpers.SendResponse(c, fiber.StatusBadRequest, "Invalid ID format", nil)
	}
	var wallet entity.UpdateWallets
	if err := c.BodyParser(&wallet); err != nil {
		return helpers.SendResponse(c, fiber.StatusBadRequest, "Invalid input", nil)
	}
	wallet.ID = id

	if err := h.walletUsecase.Update(wallet); err != nil {
		return helpers.SendResponse(c, fiber.StatusInternalServerError, "Failed to update wallet", nil)
	}

	return helpers.SendResponse(c, fiber.StatusOK, "Wallet updated successfully", nil)
}

// @Summary Delete a wallet
// @Description Delete a wallet by its ID
// @Tags wallets
// @Accept json
// @Produce json
// @Param id path int true "Wallet ID"
// @Success 204 {object} nil
// @Security BearerAuth
// @Router /api/v1/wallets/{id} [delete]
func (h *WalletHandler) DeleteWallet(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return helpers.SendResponse(c, fiber.StatusBadRequest, "Invalid ID format", nil)
	}
	if err := h.walletUsecase.Delete(id); err != nil {
		return helpers.SendResponse(c, fiber.StatusNotFound, "Wallet not found", nil)
	}

	return helpers.SendResponse(c, fiber.StatusNoContent, "Wallet deleted successfully", nil)
}
