package model

type WGEasyClient struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type WGEasyCreateClientParams struct {
	Name string `json:"name"`
}

type WGEasyCreateClientResponse struct {
	Success  bool   `json:"success"`
	ClientID string `json:"clientId"`
}
