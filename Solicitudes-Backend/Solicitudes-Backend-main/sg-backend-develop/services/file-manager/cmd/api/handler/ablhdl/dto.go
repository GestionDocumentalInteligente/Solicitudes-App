package ablhdl

type AblDataDTO struct {
	ABLNumber string `json:"abl_number" binding:"required"`
	Type      int    `json:"type" binding:"required"`
}
