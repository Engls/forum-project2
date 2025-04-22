package entity

// RegisterRequest представляет запрос на регистрацию
type RegisterRequest struct {
	Username string `json:"username" example:"user123"`
	Password string `json:"password" example:"P@ssw0rd"`
	Role     string `json:"role" example:"user"`
}

// LoginRequest представляет запрос на вход
type LoginRequest struct {
	Username string `json:"username" example:"user123"`
	Password string `json:"password" example:"P@ssw0rd"`
}
