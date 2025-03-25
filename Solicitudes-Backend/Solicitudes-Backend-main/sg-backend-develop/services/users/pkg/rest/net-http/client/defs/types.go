package defs

type GenericTokenResponse struct {
	Data map[string]interface{}
}

func (g *GenericTokenResponse) GetAccessToken() string {
	if token, ok := g.Data["access_token"].(string); ok {
		return token
	}
	return ""
}
