package filehdl

import (
	"github.com/teamcubation/sg-file-manager-api/internal/manager/docprocessor/file"
	"github.com/teamcubation/sg-file-manager-api/internal/manager/docprocessor/file/user"
)

type DocumentDTO struct {
	Name        string              `json:"name"`
	Type        file.DocumentTypeID `json:"type" binding:"required"`
	Content     string              `json:"content" binding:"required"`
	ContentType string              `json:"contentType"`
}

type RequestDataDTO struct {
	ID                 int64         `json:"id"`
	Cuil               int64         `json:"cuil" binding:"required"`
	DocumentNumber     int64         `json:"dni" binding:"required"`
	FirstName          string        `json:"first_name" binding:"required"`
	LastName           string        `json:"last_name" binding:"required"`
	Email              string        `json:"email" binding:"required"`
	Phone              int64         `json:"phone" binding:"required"`
	UserType           user.UserType `json:"user_type" binding:"required"`
	Address            AddressDTO    `json:"address" binding:"required"`
	SelectedActivities []int         `json:"tasks"`
	Activities         []string      `json:"activities"`
	ProjectDesc        string        `json:"description"`
	Observations       string        `json:"observations"`
	ObservationsTasks  string        `json:"observations_tasks"`
	EstimatedTime      int64         `json:"time"`
	Insurance          bool          `json:"insurance"`
	Documents          []DocumentDTO `json:"documents" binding:"required"`
}

type ValidatedRequestDataDTO struct {
	Cuil      int64  `json:"cuil" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Address   struct {
		Street string `json:"street"`
		Number string `json:"number"`
	} `json:"address" binding:"required"`
	Activities    []string `json:"activities"`
	ProjectDesc   string   `json:"description"`
	EstimatedTime int64    `json:"time"`
	Insurance     bool     `json:"insurance"`
}

type AddressDTO struct {
	ABLNumber int64  `json:"abl_number" binding:"required"`
	ZipCode   int64  `json:"zc" binding:"required"`
	Street    string `json:"street" binding:"required"`
	Number    string `json:"number" binding:"required"`
}
