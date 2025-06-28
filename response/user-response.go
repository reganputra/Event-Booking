package response

type UserResponse struct {
	Id       int64
	Email    string `bind:"required"`
	Password string `bind:"required" json:"-"`
	Role     string
}
