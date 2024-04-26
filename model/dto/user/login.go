package user

type (
	LoginRequest struct {
		Username string `json:"username" validate:"required" `
		Password string `json:"password" validate:"required"`
	}

	LoginResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
)
