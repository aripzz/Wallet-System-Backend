package entity

type Wallets struct {
	ID              uint64  `json:"id"`
	User_ID         uint64  `json:"user_id"`
	Wallet_Types_ID uint64  `json:"wallet_types_id"`
	Balance         float64 `json:"balance"`
	Active          bool    `json:"active"`
}

type CreateWallets struct {
	User_ID         uint64  `json:"user_id"`
	Wallet_Types_ID uint64  `json:"wallet_types_id"`
	Balance         float64 `json:"balance"`
	Active          bool    `json:"active"`
}

type UpdateWallets struct {
	ID              uint64   `json:"id" validate:"required"`
	User_ID         *uint64  `json:"user_id,omitempty"`
	Wallet_Types_ID *uint64  `json:"wallet_types_id,omitempty"`
	Balance         *float64 `json:"balance,omitempty"`
	Active          *bool    `json:"active,omitempty"`
}
