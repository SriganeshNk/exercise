package Response

import "time"

type BaseResponse struct {
	Status string
	Timestamp time.Time
	Result interface{}
}