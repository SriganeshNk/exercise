package Response

type BaseResponse struct {
	Status string
	Timestamp int64
	Result interface{}
}