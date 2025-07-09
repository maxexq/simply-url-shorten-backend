package model

type BaseResponse struct {
	Code int         `json:"error"`
	Data interface{} `json:"result"`
}

type BaseErrorResponse struct {
	Code     int    `json:"error" example:"4"`
	ErrTitle string `json:"result"`
}

type ValidatorErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func NewBaseResponse(code int, data interface{}) *BaseResponse {
	if data == nil {
		type Empty struct{}
		data = Empty{}
	}
	return &BaseResponse{
		Code: code,
		Data: data,
	}
}

func NewBaseErrorResponse(code int, title string) *BaseErrorResponse {
	return &BaseErrorResponse{
		Code:     code,
		ErrTitle: title,
	}
}
