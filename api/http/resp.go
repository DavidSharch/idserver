package api

type HttpResponse struct {
	Code int `json:"code"`
	Res  any `json:"res"`
}

func NewHttpResponse200(res any) *HttpResponse {
	return &HttpResponse{
		Code: 200,
		Res:  res,
	}
}
