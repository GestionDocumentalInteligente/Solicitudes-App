package ablhdl

import (
	"github.com/gin-gonic/gin"
)

type ABLHandlerRouter struct {
	ablhdl *ABLHandler
}

func NewRouter(hdl *ABLHandler) *ABLHandlerRouter {
	return &ABLHandlerRouter{
		ablhdl: hdl,
	}
}

func (s *ABLHandlerRouter) AddRoutesV1(v1 *gin.RouterGroup) {
	v1.POST("/address/debt", s.ablhdl.VerifyDebt)
}
