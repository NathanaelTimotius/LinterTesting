package models

type User struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Address  string `json:"address"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type UserResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    User   `json:"data"`
}

type UsersResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []User `json:"data"`
}

type SuccessResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type Products struct {
	ID    int    `json:"id" gorm:"primaryKey"`
	Name  string `json:"Name"`
	Price int    `json:"Price"`
}

type ProductResponse struct {
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Data    Products `json:"data"`
}

type ProductsResponse struct {
	Status  int        `json:"status"`
	Message string     `json:"message"`
	Data    []Products `json:"data"`
}

type Transactions struct {
	ID        int `json:"ID"`
	UserID    int `json:"UserID"`
	ProductID int `json:"ProductID"`
	Quantity  int `json:"QUantity"`
}

type TransactionResponse struct {
	Status  int          `json:"status"`
	Message string       `json:"message"`
	Data    Transactions `json:"data"`
}

type TransactionsResponse struct {
	Status  int            `json:"status"`
	Message string         `json:"message"`
	Data    []Transactions `json:"data"`
}

type DetailTransactions struct {
	ID           int            `json:"id"`
	UserID       int            `json:"-"`
	ProductID    int            `json:"-"`
	Transactions []Transactions `json:"transactions"`
	User         User           `gorm:"foreignKey:UserID" json:"user"`
	Product      Products       `gorm:"foreignKey:ProductID" json:"product"`
	Quantity     int            `json:"quantity"`
}

type DetailTransaction struct {
	ID        int      `json:"id"`
	UserID    int      `json:"-"`
	ProductID int      `json:"-"`
	User      User     `gorm:"foreignKey:UserID" json:"user"`
	Product   Products `gorm:"foreignKey:ProductID" json:"product"`
	Quantity  int      `json:"quantity"`
}

type DetailTransactionsResponse struct {
	Status  int                 `json:"status"`
	Message string              `json:"message"`
	Data    []DetailTransaction `json:"data"`
}
