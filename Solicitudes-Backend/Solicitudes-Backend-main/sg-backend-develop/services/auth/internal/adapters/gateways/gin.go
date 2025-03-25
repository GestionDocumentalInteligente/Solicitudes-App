package authgtw

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	sdkmwr "github.com/teamcubation/sg-auth/pkg/middleware/gin"
	sdkgin "github.com/teamcubation/sg-auth/pkg/rest/gin"
	sdkginports "github.com/teamcubation/sg-auth/pkg/rest/gin/ports"

	ports "github.com/teamcubation/sg-auth/internal/core/ports"
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
	apiBase := "/api/" + apiVersion + "/auth"
	protectedPrefix := apiBase + "/protected"

	// Rutas públicas
	router.GET(apiBase+"/ping", h.Ping)

	// Rutas protegidas (requieren JWT válido)
	authorized := router.Group(protectedPrefix)
	{
		// Aplicar middleware de validación JWT
		authorized.Use(sdkmwr.ValidateJwt(secrets["auth"]))
		authorized.GET("/hi", h.ProtectedHi)
		authorized.GET("/login", h.Login)
	}
}

func (h *GinHandler) ProtectedHi(c *gin.Context) {
	// cuil, err := getJwtData(c)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid jwt data: " + err.Error()})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		"message": []string{
			"hi! from protected.",
		},
	})

}

func (h *GinHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Pong!"})
}

func (h *GinHandler) Login(c *gin.Context) {
	// Obtener el CUIL del parámetro de consulta (query parameter)
	cuil := c.Query("cuil")
	if cuil == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cuil query parameter required"})
		return
	}

	// Llamar al caso de uso de Login con el CUIL proporcionado
	token, err := h.ucs.Login(c.Request.Context(), cuil)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Construir la respuesta
	response := gin.H{
		"message":      "Login successful",
		"cuil":         cuil,
		"access_token": token.AccessToken,
	}

	// Enviar la respuesta JSON
	c.JSON(http.StatusOK, response)

	// redirectURL := fmt.Sprintf("https://frontend-app.com/callback?token=%s", token.AccessToken)
	// c.Redirect(http.StatusFound, redirectURL)
}
