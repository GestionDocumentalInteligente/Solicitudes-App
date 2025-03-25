package mailgtw

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	trp "github.com/teamcubation/sg-mailing-api/internal/adapters/gateways/transport/email-verification"
	sdkmwr "github.com/teamcubation/sg-mailing-api/pkg/middleware/gin"
	sdkgin "github.com/teamcubation/sg-mailing-api/pkg/rest/gin"
	sdkgindefs "github.com/teamcubation/sg-mailing-api/pkg/rest/gin/defs"

	"github.com/teamcubation/sg-mailing-api/internal/core/entities"
	ports "github.com/teamcubation/sg-mailing-api/internal/core/ports"
)

type GinHandler struct {
	ucs       ports.UseCases
	ginServer sdkgindefs.Server
}

func NewGinHandler(u ports.UseCases) (*GinHandler, error) {
	// Aquí aceptará el tipo ports.UseCases
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
	apiBase := "/api/" + apiVersion + "/mailing"
	protectedPrefix := apiBase + "/protected"

	// Rutas públicas
	router.GET(apiBase+"/ping", h.Ping)
	router.GET(apiBase+"/activate-account", h.ConfirmEmailVerification)
	router.POST(apiBase+"/resend-activation-email", h.ResendActivationEmail)
	router.POST(apiBase+"/request-new-activation", h.NewActivationEmail)

	router.POST(apiBase+"/new-request", h.SendNewRequestMessage)
	router.POST(apiBase+"/update-request", h.SendUpdateRequestMessage)
	router.POST(apiBase+"/update-request-code", h.SendUpdateRequestByCodeMessage)
	router.POST(apiBase+"/validate-request", h.SendValidateRequestMessage)

	// Rutas protegidas (requieren JWT válido)
	protected := router.Group(protectedPrefix)
	{
		protected.Use(sdkmwr.ValidateJwt(secrets["jwt"]))
		protected.POST("/email", h.InitiateEmailVerification)
		protected.GET("/hi", h.ProtectedHi)
	}
}

func (h *GinHandler) ProtectedHi(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": []string{
			"hi! from protected.",
		},
	})

}

func (h *GinHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Pong!"})
}

func (h *GinHandler) InitiateEmailVerification(c *gin.Context) {
	var req trp.Req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format: " + err.Error()})
		return
	}

	err := h.ucs.InitiateEmailVerification(c.Request.Context(), req.Email, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to initiate email verification: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "verification email sent"})
}

type emailReq struct {
	Email        string `json:"email" binding:"required,email"`
	Code         string `json:"code" binding:"required"`
	Observations string `json:"observations"`
}

func (h *GinHandler) SendNewRequestMessage(c *gin.Context) {
	var req emailReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format: " + err.Error()})
		return
	}

	err := h.ucs.SendNewRequestMessage(c.Request.Context(), req.Code, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to initiate email verification: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "new request email sent"})
}

func (h *GinHandler) ConfirmEmailVerification(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token no proporcionado"})
		return
	}

	err := h.ucs.ActivateAccount(c.Request.Context(), token)
	if err != nil {
		switch err.(type) {
		case *entities.ErrorTokenInvalid:
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "email verified"})
}

type tokenReq struct {
	Token string
}

func (h *GinHandler) ResendActivationEmail(c *gin.Context) {
	var req *tokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.ucs.ResendActivationEmail(c.Request.Context(), req.Token)
	if err != nil {
		switch err.(type) {
		case *entities.ErrorTokenInvalid:
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		case *entities.NotFoundInDatabase:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		case *entities.UserAlreadyActive:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "email forwarded"})
}

type emailActivationReq struct {
	Email string
}

func (h *GinHandler) NewActivationEmail(c *gin.Context) {
	var req *emailActivationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email not provided"})
		return
	}

	err := h.ucs.ResendActivationEmailExistingUser(c.Request.Context(), req.Email)
	if err != nil {
		switch err.(type) {
		case *entities.ErrorTokenInvalid:
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		case *entities.NotFoundInDatabase:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		case *entities.UserAlreadyActive:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "email forwarded"})
}

func (h *GinHandler) SendUpdateRequestByCodeMessage(c *gin.Context) {
	var req emailReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format: " + err.Error()})
		return
	}

	err := h.ucs.SendUpdateRequestByCodeMessage(c.Request.Context(), req.Code, req.Email, req.Observations)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to initiate email verification: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "new request email sent"})
}

func (h *GinHandler) SendUpdateRequestMessage(c *gin.Context) {
	var req emailReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format: " + err.Error()})
		return
	}

	err := h.ucs.SendUpdateRequestMessage(c.Request.Context(), req.Code, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to initiate email verification: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "new request email sent"})
}

func (h *GinHandler) SendValidateRequestMessage(c *gin.Context) {
	var req emailReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format: " + err.Error()})
		return
	}

	err := h.ucs.SendValidateRequestMessage(c.Request.Context(), req.Code, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to initiate email verification: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "new request email sent"})
}
