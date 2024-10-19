package response

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Meta struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type BOMResponse struct {
	Meta       Meta        `json:"meta"`
	DataBom    interface{} `json:"data_bom"`
}

func SuccessResponse(code int, message string, data interface{}) Response {
	return Response{
		Meta: Meta{
			Code:    code,
			Message: message,
		},
		Data: data,
	}
}

func ErrorResponse(code int, message string) Response {
	return Response{
		Meta: Meta{
			Code:    code,
			Message: message,
		},
		Data: nil,
	}
}

func SuccessResponseBom(code int, message string, data interface{}) BOMResponse {
	return BOMResponse{
		Meta: Meta{
			Code:    code,
			Message: message,
		},
		DataBom: data,
	}
}

func ErrorResponseBom(code int, message string) BOMResponse {
	return BOMResponse{
		Meta: Meta{
			Code:    code,
			Message: message,
		},
	}
}

