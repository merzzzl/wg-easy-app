package model

import "time"

type Tunnel struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	WGClientName string    `json:"wg_client_name"`
	WGClientID   string    `json:"wg_client_id"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateTunnelParams struct {
	UserID int64
}
