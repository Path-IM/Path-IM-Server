package repository

import (
	"github.com/showurl/Path-IM-Server/app/msg-transfer/cmd/persistent/internal/svc"
)

type Rep struct {
	svcCtx *svc.ServiceContext
}

var rep *Rep

func NewRep(svcCtx *svc.ServiceContext) *Rep {
	if rep != nil {
		return rep
	}
	rep = &Rep{
		svcCtx: svcCtx,
	}
	return rep
}
