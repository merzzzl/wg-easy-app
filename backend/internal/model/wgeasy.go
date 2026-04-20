package model

type WGEasyClient struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type WGEasyCreateClientParams struct {
	Name      string  `json:"name"`
	ExpiresAt *string `json:"expiresAt"`
}

type WGEasyCreateClientResponse struct {
	Success  bool  `json:"success"`
	ClientID int64 `json:"clientId"`
}
