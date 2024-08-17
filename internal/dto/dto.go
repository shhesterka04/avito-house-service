package dto

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type DtoFlat struct {
	Id      int    `json:"id"`
	HouseId int    `json:"house_id"`
	Status  string `json:"status"`
	Number  int    `json:"number"`
	Rooms   int    `json:"rooms"`
	Price   int    `json:"price"`
}

type CreateFlatRequest struct {
	HouseID int `json:"house_id"`
	Number  int `json:"number"`
	Rooms   int `json:"rooms"`
	Price   int `json:"price"`
}
