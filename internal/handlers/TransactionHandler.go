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

type TransactionHandler struct {
	transactionUsecase usecase.TransactionUseCase
}

func NewTransactionHandler(app *fiber.App, db *gorm.DB, cache *infra.RedisClient) {
	transactionRepo := repository.NewTransactionRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	productRepo := repository.NewProductRepository(db)
	usecase := usecase.NewTransactionUseCase(transactionRepo, productRepo, walletRepo, db, cache)
	handler := &TransactionHandler{transactionUsecase: usecase}

	apiv1 := app.Group("/api/v1")
	apiv1.Post("/transactions", handler.CreateTransaction)
	apiv1.Get("/transactions", handler.GetAllTransactions)
	apiv1.Get("/transactions/:id", handler.GetTransactionByID)
	apiv1.Get("/users/transactions-user", handler.GetTransactionsByUserID)
	apiv1.Delete("/transactions/:id", handler.DeleteTransaction)
}

// @Summary Create a new transaction
// @Description Create a new transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body entity.RequestCreateTransactions true "Transaction data"
// @Success 201 {object} helpers.StandardResponse
// @Security BearerAuth
// @Router /api/v1/transactions [post]
func (h *TransactionHandler) CreateTransaction(c *fiber.Ctx) error {
	var transaction entity.RequestCreateTransactions
	if err := c.BodyParser(&transaction); err != nil {
		return helpers.SendResponse(c, fiber.StatusBadRequest, "Invalid input", nil)
	}
	userID, ok := c.Locals("userID").(uint64)
	if !ok {
		helpers.SendResponse(c, fiber.StatusInternalServerError, "user ID not found", nil)
	}

	if err := h.transactionUsecase.CreateTransaction(userID, transaction); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return helpers.SendResponse(c, fiber.StatusCreated, "Transaction created successfully", nil)
}

// @Summary Get all transactions
// @Description Get a list of all transactions
// @Tags transactions
// @Accept json
// @Produce json
// @Success 200 {array} entity.Transactions
// @Security BearerAuth
// @Router /api/v1/transactions [get]
func (h *TransactionHandler) GetAllTransactions(c *fiber.Ctx) error {
	transactions, err := h.transactionUsecase.GetAllTransactions()
	if err != nil {
		return helpers.SendResponse(c, fiber.StatusInternalServerError, "Failed to retrieve transactions", nil)
	}

	return helpers.SendResponse(c, fiber.StatusOK, "Transactions retrieved successfully", transactions)
}

// @Summary Get transaction by ID
// @Description Get a transaction by its ID
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} entity.Transactions
// @Security BearerAuth
// @Router /api/v1/transactions/{id} [get]
func (h *TransactionHandler) GetTransactionByID(c *fiber.Ctx) error {
	idStr := c.Params("id")                     // Get the ID as a string
	id, err := strconv.ParseUint(idStr, 10, 64) // Convert string to uint64
	if err != nil {
		return helpers.SendResponse(c, fiber.StatusBadRequest, "Invalid ID format", nil)
	}

	transaction, err := h.transactionUsecase.GetTransactionByID(id) // Use the uint64 ID
	if err != nil {
		return helpers.SendResponse(c, fiber.StatusNotFound, "Transaction not found", nil)
	}

	return helpers.SendResponse(c, fiber.StatusOK, "Transaction retrieved successfully", transaction)
}

// @Summary Get transactions by user ID
// @Description Get transactions for a specific user by their ID
// @Tags transactions
// @Accept json
// @Produce json
// @Success 200 {array} entity.Transactions
// @Security BearerAuth
// @Router /api/v1/users/transactions-user [get]
func (h *TransactionHandler) GetTransactionsByUserID(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uint64)
	if !ok {
		helpers.SendResponse(c, fiber.StatusInternalServerError, "user ID not found", nil)
	}

	transactions, err := h.transactionUsecase.GetTransactionsByUserID(userID)
	if err != nil {
		return helpers.SendResponse(c, fiber.StatusInternalServerError, "Failed to retrieve transactions", nil)
	}

	return helpers.SendResponse(c, fiber.StatusOK, "Transactions retrieved successfully", transactions)
}

// @Summary Delete a transaction
// @Description Delete a transaction by its ID
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 204 {object} nil
// @Security BearerAuth
// @Router /api/v1/transactions/{id} [delete]
func (h *TransactionHandler) DeleteTransaction(c *fiber.Ctx) error {
	idStr := c.Params("id")                     // Get the ID as a string
	id, err := strconv.ParseUint(idStr, 10, 64) // Convert string to uint64
	if err != nil {
		return helpers.SendResponse(c, fiber.StatusBadRequest, "Invalid ID format", nil)
	}

	if err := h.transactionUsecase.DeleteTransaction(id); err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Transaction not found")
	}

	return helpers.SendResponse(c, fiber.StatusNoContent, "Transaction deleted successfully", nil)
}
