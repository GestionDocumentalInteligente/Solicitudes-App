package ablhdl

import (
	"errors"
	"net/http"

	"github.com/teamcubation/sg-file-manager-api/internal/manager"
	"github.com/teamcubation/sg-file-manager-api/internal/manager/docprocessor/file"
	"github.com/teamcubation/sg-file-manager-api/pkg/log"

	"github.com/gin-gonic/gin"
)

type ABLHandler struct {
	useCase manager.ABLUseCase
}

func NewABLHandler(uc manager.ABLUseCase) *ABLHandler {
	return &ABLHandler{
		useCase: uc,
	}
}

func (h *ABLHandler) VerifyDebt(c *gin.Context) {
	ctx := log.Context(c.Request)

	var requestData AblDataDTO

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	info, err := h.useCase.ValidateABLData(ctx, requestData.ABLNumber, requestData.Type)
	if err != nil {
		var customErr *file.CustomError
		if errors.As(err, &customErr) {
			c.JSON(customErr.StatusCode, gin.H{"error": customErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"dbt": info})
}
