package usergtw

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	sdkmwr "github.com/teamcubation/sg-users/pkg/middleware/gin"
	sdkgin "github.com/teamcubation/sg-users/pkg/rest/gin"
	sdkginports "github.com/teamcubation/sg-users/pkg/rest/gin/defs"

	trncreate "github.com/teamcubation/sg-users/internal/adapters/gateways/transport/create-user"
	trnget "github.com/teamcubation/sg-users/internal/adapters/gateways/transport/get-user"
	trnupdate "github.com/teamcubation/sg-users/internal/adapters/gateways/transport/update-user"

	"github.com/teamcubation/sg-users/internal/core/entities"
	ports "github.com/teamcubation/sg-users/internal/core/ports"
)

type GinHandler struct {
	ucs       ports.UseCases
	ginServer sdkginports.Server
}

func NewGinHandler(u ports.UseCases) (*GinHandler, error) {
	ginServer, err := sdkgin.Bootstrap()
	if err != nil {
		return nil, fmt.Errorf("gin Service error: %w", err)
	}

	return &GinHandler{
		ucs:       u,
		ginServer: ginServer,
	}, nil
}

func (h *GinHandler) GetRouter() *gin.Engine {
	return h.ginServer.GetRouter()
}

func (h *GinHandler) Start() error {
	// Cargar el secret solo cuando sea necesario
	secrets, err := getSecrets()
	if err != nil {
		return fmt.Errorf("failed to load secrets: %w", err)
	}

	// Configurar rutas
	h.routes(secrets)

	// Iniciar el servidor
	return h.ginServer.RunServer()
}

func (h *GinHandler) routes(secrets map[string]string) {
	router := h.ginServer.GetRouter()

	// Definir prefijos de ruta
	apiVersion := h.ginServer.GetApiVersion()
	apiBase := "/api/" + apiVersion + "/users"
	validatedPrefix := apiBase + "/validated"

	// Rutas para crear usuario
	// NOTE: eliminar, solo para pruebas, usar en protected
	// router.POST(apiBase, h.CreateUser)
	// router.PATCH(apiBase, h.UpdateUserByCuil)
	// router.GET(apiBase, h.GetUserByCuil)

	// Rutas públicas
	router.GET(apiBase+"/ping", h.Ping)

	// Rutas validadas (requieren validación de credenciales)
	validated := router.Group(validatedPrefix)
	{
		// Aplicar middleware de validación de credenciales
		validated.Use(sdkmwr.ValidateCredentials())
		// Puedes añadir rutas aquí si es necesario
	}

	// Rutas protegidas (requieren JWT válido)
	protected := router.Group(apiBase, sdkmwr.ValidateJwt(secrets["auth"]))
	{
		protected.GET("/protected-hi", h.ProtectedHi)
		protected.POST("", h.CreateUser)
		protected.PATCH("", h.UpdateUserByCuil)
		protected.GET("", h.GetUserByCuil)
	}
}

// ProtectedHi responde desde una ruta protegida
func (h *GinHandler) ProtectedHi(c *gin.Context) {
	cuil, err := getJwtData(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JWT data: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": []string{
			"Hi! From protected.",
			"CUIL: " + cuil,
		},
	})
}

// Ping responde con "Pong!" para verificar que el servidor está funcionando
func (h *GinHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Pong!"})
}

// CreateUser maneja la creación de un nuevo usuario
func (h *GinHandler) CreateUser(c *gin.Context) {
	var req *trncreate.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_REQUEST_BODY",
			"message": err.Error()})
		return
	}

	userUUID, err := h.ucs.CreateUser(c.Request.Context(), trncreate.ToUserDto(req))
	if err != nil {
		if errors.Is(err, entities.ErrUserAlreadyExists) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "USER_ALREADY_EXISTS",
				"message": "User with this email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating user: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created",
		"uuid":    userUUID,
	})
}

func (h *GinHandler) UpdateUserByCuil(c *gin.Context) {
	var req *trnupdate.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userUUID, err := h.ucs.UpdateUserByPersonCuil(c.Request.Context(), trnupdate.ToUserDto(req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user updated",
		"uuid":    userUUID,
	})
}

// GetUserByCuil handles the request to get a user by CUIL
func (h *GinHandler) GetUserByCuil(c *gin.Context) {
	// Get the CUIL from query parameters
	cuil := c.Query("cuil")
	if cuil == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "'cuil' parameter is required"})
		return
	}

	// Call the use case to get the user by CUIL
	user, err := h.ucs.FindUserByPersonCuil(c.Request.Context(), cuil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving user: " + err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Convert the user entity to a response DTO
	response := trnget.ToUserResponse(user)

	c.JSON(http.StatusOK, response)
}
