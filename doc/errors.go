package docs

type ErrorResponseInternalServer struct {
	Title  string `json:"title" example:"Internal server error"`
	Status int    `json:"status" example:"500"`
	Detail string `json:"detail" example:"internal server error"`
}

type ErrorResponseNotFound struct {
	Title  string `json:"title" example:"Not Found"`
	Status int    `json:"status" example:"404"`
	Detail string `json:"detail" example:"not Found"`
}

type ErrorResponseBadRequest struct {
	Title  string `json:"title" example:"Bad Request"`
	Status int    `json:"status" example:"400"`
	Detail string `json:"detail" example:"bad request"`
}

type ErrorResponseInvCityName struct {
	Title  string `json:"title" example:"Bad request"`
	Status int    `json:"status" example:"400"`
	Detail string `json:"detail" example:"city name is not valid"`
}
