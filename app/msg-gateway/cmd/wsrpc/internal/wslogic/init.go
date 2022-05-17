package wslogic

import (
	"github.com/go-playground/validator/v10"
	"sync"
)

var (
	rwLock          *sync.RWMutex
	validate        *validator.Validate
	sendMsgAllCount uint64
	userCount       uint64
)

func init() {
	rwLock = new(sync.RWMutex)
	validate = validator.New()
}
