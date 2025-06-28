package types

type Student struct {
	Id    int64
	Name  string `validate:"required,min=3,max=50"`
	Email string `validate:"required,email"`
	Age   int    `validate:"required,numeric,min=18,max=100"`
}
