package tldstats

import (
	"context"
	"strings"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"golang.org/x/net/publicsuffix"
)

type TLDStats struct {
	Next plugin.Handler
}

func (s *TLDStats) Name() string { return pluginName }

func (s *TLDStats) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	eTLD, _ := publicsuffix.PublicSuffix(strings.TrimSuffix(state.Name(), "."))

	qType, ok := dns.TypeToString[state.QType()]
	if !ok {
		qType = "unknown"
	}

	RequestCount.WithLabelValues(eTLD, qType).Inc()

	return plugin.NextOrFailure(s.Name(), s.Next, ctx, w, r)
}
