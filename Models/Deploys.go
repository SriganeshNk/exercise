package Models

import (
	"time"
)

type Deploys struct {
	Id int64
	Sha string
	Date time.Time
	Action string
	Engineer string
}
