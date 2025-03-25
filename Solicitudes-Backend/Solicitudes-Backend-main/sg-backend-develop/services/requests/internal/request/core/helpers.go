package request

import (
	"strconv"
	"strings"

	sdktools "github.com/teamcubation/sg-backend/pkg/tools"

	domain "github.com/teamcubation/sg-backend/services/requests/internal/request/core/domain"
)

func parseInput(inputText string) domain.Suggestion {
	var addrNameParts []string
	var addrNum, ablNum int64
	var nameStr string

	tokens := strings.Fields(inputText)

	for _, token := range tokens {
		if sdktools.IsNumeric(token) {
			num, _ := strconv.ParseInt(token, 10, 64)
			if len(token) == 6 {
				ablNum = num
			} else if len(token) <= 4 {
				addrNum = num
			}
		} else {
			nameStr += sdktools.NormalizeString(token)
		}
	}

	if nameStr != "" {
		addrNameParts = append(addrNameParts, nameStr)
	}

	AddrStreet := strings.Join(addrNameParts, " ")

	return domain.Suggestion{
		AddrStreet: AddrStreet,
		AddrNum:    addrNum,
		AblNum:     ablNum,
	}
}
