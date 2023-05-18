package docs

type ErrorInternalServer struct {
	Title  string `json:"title" example:"Internal server error"`
	Status int    `json:"status" example:"500"`
	Detail string `json:"detail" example:"internal server error"`
}

type ErrorNotFound struct {
	Title  string `json:"title" example:"Not Found"`
	Status int    `json:"status" example:"404"`
	Detail string `json:"detail" example:"not found"`
}

type ErrorUserNotFound struct {
	Title  string `json:"title" example:"Not Found"`
	Status int    `json:"status" example:"404"`
	Detail string `json:"detail" example:"user not found"`
}

type ErrorCarNotFound struct {
	Title  string `json:"title" example:"Not Found"`
	Status int    `json:"status" example:"404"`
	Detail string `json:"detail" example:"car not found"`
}

type ErrorReservationNotFound struct {
	Title  string `json:"title" example:"Not Found"`
	Status int    `json:"status" example:"404"`
	Detail string `json:"detail" example:"reservation was not found"`
}

type ErrorEmailAlreadyRegistered struct {
	Title  string `json:"title" example:"Bad Request"`
	Status int    `json:"status" example:"400"`
	Detail string `json:"detail" example:"email already registered"`
}

type ErrorInvalidEmail struct {
	Title  string `json:"title" example:"Bad Request"`
	Status int    `json:"status" example:"400"`
	Detail string `json:"detail" example:"invalid email"`
}

type ErrorinvalidUUID struct {
	Title  string `json:"title" example:"Bad Request"`
	Status int    `json:"status" example:"400"`
	Detail string `json:"detail" example:"id could not be converted to uuid"`
}

type ErrorInvalidReservationTimeFrame struct {
	Title  string `json:"title" example:"Bad Request"`
	Status int    `json:"status" example:"400"`
	Detail string `json:"detail" example:"reservation time frame is invalid"`
}

type ErrorCarNotAvailable struct {
	Title  string `json:"title" example:"Bad Request"`
	Status int    `json:"status" example:"400"`
	Detail string `json:"detail" example:"car not available"`
}

type ErrorInvalidCarStatus struct {
	Title  string `json:"title" example:"Bad Request"`
	Status int    `json:"status" example:"400"`
	Detail string `json:"detail" example:"invalid car status"`
}

type ErrorCityQueryParamEmpty struct {
	Title  string `json:"title" example:"city query param can not be empty"`
	Status int    `json:"status" example:"400"`
	Detail string `json:"detail" example:"car not available"`
}

type ErrorInvalidReservationStatus struct {
	Title  string `json:"title" example:"Bad Request"`
	Status int    `json:"status" example:"400"`
	Detail string `json:"detail" example:"invalid reservation status"`
}

type ErrorMinimumReservationHours struct {
	Title  string `json:"title" example:"Bad Request"`
	Status int    `json:"status" example:"400"`
	Detail string `json:"detail" example:"period is shorter than minimun allowed (6 hours)"`
}

type ErrorInvalidCityName struct {
	Title  string `json:"title" example:"Bad request"`
	Status int    `json:"status" example:"400"`
	Detail string `json:"detail" example:"city name is not valid"`
}

type ErrorInvalidTimeFrame struct {
	Title  string `json:"title" example:"Bad Request"`
	Status int    `json:"status" example:"400"`
	Detail string `json:"detail" example:"invalid time frame"`
}
