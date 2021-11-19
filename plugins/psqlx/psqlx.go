package psqlx

import (
	"github.com/coredns/coredns/plugin"
	"github.com/jackc/pgx/v4/pgxpool"
)

const Name = "psqlx"

type PSQLX struct {
	Next   plugin.Handler
	dbPool *pgxpool.Pool
}

func (px PSQLX) Name() string {
	return Name
}

type DomainInfo struct {
	ID         uint64
	Name       string
	SoaMname   string
	SoaRname   string
	SoaSerial  uint32
	SoaRefresh uint32
	SoaRetry   uint32
	SoaExpire  uint32
	SoaMinimum uint32
}
