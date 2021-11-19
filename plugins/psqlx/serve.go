package psqlx

import (
	"context"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/jackc/pgx/v4"
	"github.com/miekg/dns"
)

func (px PSQLX) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	a := new(dns.Msg)
	a.SetReply(r)
	a.Compress = true
	a.Authoritative = true

	domain := domainNormalize(state.QName())

	domainInfo := DomainInfo{}

	for _, row := range domainNameGuesses(domain) {
		err := px.dbPool.QueryRow(ctx,
			"SELECT id,name,soa_mname,soa_rname,soa_serial,soa_refresh,soa_retry,soa_expire,soa_minimum FROM ccdns_domains WHERE name=$1 AND disabled=false",
			row,
		).Scan(
			&domainInfo.ID, &domainInfo.Name, &domainInfo.SoaMname, &domainInfo.SoaRname,
			&domainInfo.SoaSerial, &domainInfo.SoaRefresh, &domainInfo.SoaRetry, &domainInfo.SoaExpire, &domainInfo.SoaMinimum,
		)

		if err != nil && err == pgx.ErrNoRows {
			continue
		} else if err != nil {
			return 0, err
		} else {
			break
		}
	}

	if domainInfo.ID == 0 {
		return plugin.NextOrFailure(px.Name(), px.Next, ctx, w, r)
	}

	switch state.QType() {
	case dns.TypeSOA:
		soa := dns.SOA{
			Hdr: dns.RR_Header{
				Name:   domainZoneNormalize(domainInfo.Name),
				Rrtype: dns.TypeSOA,
				Class:  dns.ClassINET,
				Ttl:    domainInfo.SoaMinimum,
			},
			Ns:      domainZoneNormalize(domainInfo.SoaMname),
			Mbox:    domainZoneNormalize(domainInfo.SoaRname),
			Serial:  domainInfo.SoaSerial,
			Refresh: domainInfo.SoaRefresh,
			Retry:   domainInfo.SoaRetry,
			Expire:  domainInfo.SoaExpire,
			Minttl:  domainInfo.SoaMinimum,
		}

		packed, err := dns.NewRR(soa.String())
		if err != nil {
			return 0, err
		}

		if domain == domainInfo.Name {
			a.Answer = append(a.Answer, packed)
		} else {
			a.Ns = append(a.Ns, packed)
		}

	}

	return 0, w.WriteMsg(a)
}
