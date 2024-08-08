package api_key

type ScopeFile struct {
	UUID  string `json:"uuid"`
	Title string `json:"title"`
}

type ProductScope struct {
	UUID   string      `json:"uuid"`
	Name   string      `json:"name"`
	Scopes []ScopeFile `json:"scopes"`
}
