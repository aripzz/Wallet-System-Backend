package entity

type Transactions struct {
	ID          uint64  `json:"id"`
	UsersID     uint64  `json:"users_id"`
	ProductsID  uint64  `json:"products_id"`
	Amount      float64 `json:"amount"`
	Status      string  `json:"status"`
	Username    string  `json:"username"`
	ProductName string  `json:"product_name"`
}
type RequestCreateTransactions struct {
	ProductID    uint64 `json:"product_id" validate:"required"`
	WalletTypeID uint64 `json:"wallet_type_id" validate:"required"`
}
