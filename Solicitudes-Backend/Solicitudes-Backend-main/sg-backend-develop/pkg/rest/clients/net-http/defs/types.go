package defs

type GenericTokenResponse struct {
	TokenData map[string]interface{}
}

func (g *GenericTokenResponse) GetAccessToken() string {
	if token, ok := g.TokenData["access_token"].(string); ok {
		return token
	}
	return ""
}
