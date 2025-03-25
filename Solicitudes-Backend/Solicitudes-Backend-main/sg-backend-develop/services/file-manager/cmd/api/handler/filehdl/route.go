package filehdl

import (
	"github.com/gin-gonic/gin"
)

type FileHandlerRouter struct {
	filehdl *FileHandler
}

func NewRouter(hdl *FileHandler) *FileHandlerRouter {
	return &FileHandlerRouter{
		filehdl: hdl,
	}
}

func (s *FileHandlerRouter) AddRoutesV1(v1 *gin.RouterGroup) {
	v1.POST("/record", s.filehdl.CreateRecord)
	v1.POST("/record/:id/documents", s.filehdl.SendDocumentToRecord)
	v1.PUT("/record/:id/documents", s.filehdl.UpdateDocumentsInRecord)
	v1.GET("/ping", s.filehdl.Ping)
}
