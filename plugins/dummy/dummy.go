package dummy

import (
	"context"

	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

type Dummy struct {
	Next               plugin.Handler
	AnswerMinimisation AnswerMinimisationType
	ResponseRecords    map[string]Records
}

func New() *Dummy {
	return &Dummy{
		ResponseRecords: make(map[string]Records),
	}
}

func (d *Dummy) Name() string { return pluginName }

func (d *Dummy) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {

	return plugin.NextOrFailure(d.Name(), d.Next, ctx, w, r)

	// return dns.RcodeSuccess, nil
}
