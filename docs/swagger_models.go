package docs

// User request dan response models untuk Swagger

// RegisterRequest model info
// @Description Register user request payload
type RegisterRequest struct {
	Username string `json:"username" example:"johndoe"`
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"password123"`
}

// LoginRequest model info
// @Description Login user request payload
type LoginRequest struct {
	Username string `json:"username" example:"johndoe"`
	Password string `json:"password" example:"password123"`
}

// UserResponse model info
// @Description User response payload
type UserResponse struct {
	ID       uint   `json:"id" example:"1"`
	Username string `json:"username" example:"johndoe"`
	Email    string `json:"email" example:"john@example.com"`
	Role     string `json:"role" example:"user"`
}

// TokenResponse model info
// @Description Token response payload
type TokenResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresIn    int    `json:"expires_in" example:"3600"`
}

// PostRequest model info
// @Description Post request payload
type PostRequest struct {
	Title  string `json:"title" example:"Judul Post"`
	Body   string `json:"body" example:"Isi konten post"`
	UserID uint   `json:"user_id" example:"1"`
}

// PostResponse model info
// @Description Post response payload
type PostResponse struct {
	ID        uint         `json:"id" example:"1"`
	Title     string       `json:"title" example:"Judul Post"`
	Body      string       `json:"body" example:"Isi konten post"`
	UserID    uint         `json:"user_id" example:"1"`
	User      UserResponse `json:"user"`
	CreatedAt string       `json:"created_at" example:"2023-01-01T12:00:00Z"`
	UpdatedAt string       `json:"updated_at" example:"2023-01-01T12:00:00Z"`
}

// ErrorResponse model info
// @Description Error response payload
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid credentials"`
}