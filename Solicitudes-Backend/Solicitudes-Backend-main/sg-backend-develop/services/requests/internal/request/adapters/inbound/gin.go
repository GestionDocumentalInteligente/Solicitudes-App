package inbound

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	sdkmwr "github.com/teamcubation/sg-backend/pkg/rest/middlewares/gin"
	sdkgin "github.com/teamcubation/sg-backend/pkg/rest/servers/gin"
	sdkdefs "github.com/teamcubation/sg-backend/pkg/rest/servers/gin/defs"

	config "github.com/teamcubation/sg-backend/services/requests/internal/config"
	transport "github.com/teamcubation/sg-backend/services/requests/internal/request/adapters/inbound/transport"
	ports "github.com/teamcubation/sg-backend/services/requests/internal/request/core/ports"
)

type GinHandler struct {
	ucs       ports.UseCases
	ginServer sdkdefs.Server
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

func (h *GinHandler) Start(ctx context.Context) error {
	h.routes()
	return h.ginServer.RunServer(ctx)
}

func (h *GinHandler) routes() {
	router := h.ginServer.GetRouter()

	// Definir prefijos de ruta
	apiVersion := h.ginServer.GetApiVersion()
	apiBase := "/api/" + apiVersion + "/requests"
	validatedPrefix := apiBase + "/validated"
	protectedPrefix := apiBase + "/protected"

	// Rutas públicas
	router.GET(apiBase+"/ping", h.Ping)
	router.GET(apiBase+"/address/autocomplete", h.GetSuggestions)
	router.GET(apiBase+"/abl/ownership", h.CheckAblOwnership)

	// Rutas validadas (requieren validación de credenciales)
	validated := router.Group(validatedPrefix)
	{
		// Aplicar middleware de validación de credenciales
		validated.Use(sdkmwr.ValidateCredentials())
		// Puedes añadir rutas aquí si es necesario
	}

	// Rutas protegidas (requieren JWT válido)
	protected := router.Group(protectedPrefix)
	{
		// Aplicar middleware de validación JWT
		protected.Use(sdkmwr.Validate(config.GetMiddlewareConfig().Auth))
		protected.GET("/ping", h.ProtectedPing)
		protected.POST("/create", h.CreateRequestByCuil)
		protected.PUT("/:id", h.UpdateRequestByCuil)
		protected.PUT("/verification/:id", h.VerifyRequest)
		protected.PUT("/validation/:id", h.ValidateRequest)
		protected.GET("/verification/owner", h.RequestsVerifications)
		protected.GET("/verifications", h.GetAllVerifications)
		protected.GET("/validations", h.GetAllValidations)
		protected.GET("/documents", h.GetDocumentsByFileNumber)
		protected.GET("/validations/documents", h.GetValidationDocumentsByFileNumber)
		protected.GET("/documents/:id", h.GetDocumentByID)
		protected.GET("/get/all/cuil", h.GetAllRequestsByCuil)
		protected.GET("get/all/id", h.GetAllRequestsByUserID)
		protected.GET("get/id", h.GetRequestByID)
		protected.GET("get/:code", h.GetRequestExpCode)
	}
}

func (h *GinHandler) ProtectedPing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Protected Pong!"})
}

func (h *GinHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Pong!"})
}

func (h *GinHandler) GetAllRequestsByUserID(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, transport.ErrorResponse{
			Error: transport.ErrInvalidUserID.Error(),
		})
		return
	}

	allRequests, err := h.ucs.GetAllRequestsByUserID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, transport.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, transport.RequestsResponse{
		Requests: transport.ToRequestListPresenter(allRequests),
	})
}

func (h *GinHandler) CreateRequestsByUserID(c *gin.Context) {
	var req transport.RequestJson
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, transport.ErrorResponse{
			Error: transport.ErrInvalidPayload.Error(),
		})
		return
	}

	err := h.ucs.CreateRequestByUserID(c, transport.ToRequestDomain(&req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, transport.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, transport.MessageResponse{
		Message: "Request created successfully",
	})
}

func (h *GinHandler) GetSuggestions(c *gin.Context) {
	inputText := c.Query("q")
	if len(inputText) == 0 {
		c.JSON(http.StatusBadRequest, transport.ErrorResponse{
			Error: transport.ErrMissingQueryParam.Error(),
		})
		return
	}

	list, err := h.ucs.GetSuggestions(c, inputText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, transport.ErrorResponse{
			Error: fmt.Sprintf("%v: %v", transport.ErrInternalServer, err),
		})
		return
	}

	if len(list) == 0 {
		c.JSON(http.StatusNotFound, transport.ErrorResponse{
			Error: transport.ErrNoSuggestionsFound.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, transport.SuggestionsResponse{
		Suggestions: transport.ToSuggestionPresenter(list),
	})
}

func (h *GinHandler) CheckAblOwnership(c *gin.Context) {
	cuil := c.Query("cuil")

	ablNumb, err := strconv.Atoi(c.Query("abl_number"))
	if err != nil {
		c.JSON(http.StatusBadRequest, transport.ErrorResponse{
			Error: fmt.Sprintf("Número ABL inválido: %v", err.Error()),
		})
		return
	}

	ownership, err := h.ucs.CheckAblOwnership(c, cuil, ablNumb)
	if err != nil {
		c.JSON(http.StatusInternalServerError, transport.ErrorResponse{
			// Error: transport.ErrInternalServer.Error() + ":" + err.Error(),
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, transport.AblOwnershipResponse{
		AblOwnership: ownership,
	})
}

func (h *GinHandler) GetAllVerifications(c *gin.Context) {
	verifications, err := h.ucs.AllRequestsVerifications(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, transport.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, transport.ToVerificationListPresenter(verifications))
}

func (h *GinHandler) GetAllValidations(c *gin.Context) {
	verifications, err := h.ucs.AllRequestsValidations(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, transport.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, transport.ToVerificationListPresenter(verifications))
}

func (h *GinHandler) GetDocumentsByFileNumber(c *gin.Context) {
	fileNumber := c.Query("recordNumber")

	if fileNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid 'recordNumber' query param",
		})
		return
	}

	request, documents, code, err := h.ucs.DocumentsByCode(c, fileNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, transport.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, transport.DocumentsResponse{
		Documents: transport.ToDocumentListPresenter(request, documents, code),
	})
}

func (h *GinHandler) GetValidationDocumentsByFileNumber(c *gin.Context) {
	fileNumber := c.Query("recordNumber")

	if fileNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid 'recordNumber' query param",
		})
		return
	}

	req, documents, ifDocs, err := h.ucs.ValidationDocumentsByCode(c, fileNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, transport.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"requiresInsurance": req.Insurance,
		"estimatedTime":     fmt.Sprintf("%d dias", req.EstimatedTime),
		"assignedTask":      req.Activities,
		"replacementIFs":    transport.ToReplacementDocumentList(ifDocs),
		"documents":         transport.ToDocumentList(documents),
	})
}

func (h *GinHandler) GetDocumentByID(c *gin.Context) {
	id := c.Param("id")

	document, err := h.ucs.DocumentByID(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, transport.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, transport.ToDocumentPresenter(document))
}

func (h *GinHandler) RequestsVerifications(c *gin.Context) {
	cuil, err := sdkmwr.ExtractClaim(c, "sub", "")
	if err != nil {
		c.JSON(http.StatusUnauthorized, transport.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	verification, err := h.ucs.RequestsVerifications(c, cuil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, transport.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, transport.VerificationResponse{
		Verification: transport.ToVerificationRequestPresenter(verification),
	})
}

func (h *GinHandler) CreateRequestByCuil(c *gin.Context) {
	cuil, err := sdkmwr.ExtractClaim(c, "sub", "")
	if err != nil {
		c.JSON(http.StatusUnauthorized, transport.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	var req transport.RequestJson
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, transport.ErrorResponse{
			Error: transport.ErrInvalidPayload.Error(),
		})
		return
	}

	req.Cuil = cuil
	ctx := context.Background()
	err = h.ucs.CreateRequestByCuil(ctx, transport.ToRequestDomain(&req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, transport.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Request created successfully"})
}

func (h *GinHandler) VerifyRequest(c *gin.Context) {
	cuil, err := sdkmwr.ExtractClaim(c, "sub", "")
	if err != nil {
		c.JSON(http.StatusUnauthorized, transport.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	var req transport.VerifiedRequestJson
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, transport.ErrorResponse{
			Error: transport.ErrInvalidPayload.Error(),
		})
		return
	}

	if req.VerificationType != "tasks" && req.VerificationType != "property" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid param"})
		return
	}

	req.Cuil = cuil
	req.FileNumber = c.Param("id")
	ctx := context.Background()
	err = h.ucs.UpdateRequest(ctx, transport.ToVerifiedRequestDomain(&req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, transport.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Request updated successfully"})
}

func (h *GinHandler) ValidateRequest(c *gin.Context) {
	cuil, err := sdkmwr.ExtractClaim(c, "sub", "")
	if err != nil {
		c.JSON(http.StatusUnauthorized, transport.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	var req transport.ValidateRequestJson
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, transport.ErrorResponse{
			Error: transport.ErrInvalidPayload.Error(),
		})
		return
	}

	req.Cuil = cuil
	req.FileNumber = c.Param("id")
	ctx := context.Background()
	err = h.ucs.ValidateRequest(ctx, transport.ToValidateRequestDomain(&req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, transport.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Request validated successfully"})
}

func (h *GinHandler) GetAllRequestsByCuil(c *gin.Context) {
	cuil, err := sdkmwr.ExtractClaim(c, "sub", "")
	if err != nil {
		c.JSON(http.StatusUnauthorized, transport.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	allRequests, err := h.ucs.GetAllRequestsByCuil(c, cuil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, transport.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, transport.RequestsResponse{
		Requests: transport.ToRequestListPresenter(allRequests),
	})

}

func (h *GinHandler) GetRequestByID(c *gin.Context) {
	reqID, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, transport.ErrorResponse{
			Error: transport.ErrInvalidUserID.Error(),
		})
		return
	}

	request, err := h.ucs.GetRequestByID(c, reqID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, transport.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"request": request})
}

func (h *GinHandler) GetRequestExpCode(c *gin.Context) {
	code := c.Param("code")
	request, err := h.ucs.GetRequestByExpCode(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, transport.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"request": request})
}

func (h *GinHandler) UpdateRequestByCuil(c *gin.Context) {
	cuil, err := sdkmwr.ExtractClaim(c, "sub", "")
	if err != nil {
		c.JSON(http.StatusUnauthorized, transport.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	var req transport.RequestJson
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, transport.ErrorResponse{
			Error: transport.ErrInvalidPayload.Error(),
		})
		return
	}

	request := transport.ToRequestDomain(&req)
	request.FileNumber = c.Param("id")
	request.Cuil = cuil
	ctx := context.Background()

	if err := h.ucs.UpdateRequestByFileNumber(ctx, request); err != nil {
		c.JSON(http.StatusInternalServerError, transport.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Request updated successfully"})
}
