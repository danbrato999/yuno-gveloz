package domain

type Dish struct {
	Name string `json:"name" binding:"required"`
}
