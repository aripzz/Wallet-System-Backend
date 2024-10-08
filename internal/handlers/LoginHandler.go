package handlers

import (
	"Wallet-System-Backend/internal/entity"
	"Wallet-System-Backend/internal/helpers"
	"Wallet-System-Backend/internal/repository"
	"Wallet-System-Backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type LoginHandler struct {
	loginUsecase usecase.LoginUsecase
}

func NewLoginHandler(app *fiber.App, db *gorm.DB) {
	repo := repository.NewUserRepository(db)
	usecase := usecase.NewLoginUsecase(repo)
	handler := &LoginHandler{loginUsecase: usecase}

	apiv1 := app.Group("/api/v1")

	apiv1.Post("/login", handler.Login)
}

// @Summary User Login
// @Description Authenticate a user and return a JWT token
// @Tags Login
// @Accept  json
// @Produce  json
// @Param login body entity.UserLogin true "Login"
// @Success 200 {array} helpers.StandardResponse
// @Router /api/v1/login [post]
func (h *LoginHandler) Login(c *fiber.Ctx) error {
	var user entity.UserLogin
	if err := c.BodyParser(&user); err != nil {
		return helpers.SendResponse(c, fiber.StatusBadRequest, "Invalid input", nil)
	}

	token, err := h.loginUsecase.Authenticate(user)
	if err != nil || token == "" {
		return helpers.SendResponse(c, fiber.StatusUnauthorized, "Invalid credentials", nil)
	}

	return helpers.SendResponse(c, fiber.StatusOK, "Logged in successfully", map[string]interface{}{"token": token})
}
