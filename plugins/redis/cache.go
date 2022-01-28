package redis

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"strconv"
	"time"

	"github.com/coredns/coredns/plugin/pkg/response"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
)

// Return key under which we store the message.
// Currently we do not cache Truncated, errors, zone transfers or dynamic update messages.
func key(dontUseHash bool, m *dns.Msg, t response.Type, do bool) (bool, string) {
	// We don't store truncated responses.
	if m.Truncated {
		return false, ""
	}
	// Nor errors or Meta or Update
	if t == response.OtherError || t == response.Meta || t == response.Update {
		return false, ""
	}

	if dontUseHash {
		return true, hashString(m.Question[0].Name, m.Question[0].Qtype, do)
	} else {
		return true, hashInt(m.Question[0].Name, m.Question[0].Qtype, do)
	}
}

var (
	one  = []byte("1")
	zero = []byte("0")
)

func hashString(qname string, qtype uint16, do bool) string {
	doSet := 0
	if do {
		doSet = 1
	}

	return fmt.Sprintf(
		"%s|%d|%d",
		qname, qtype, doSet,
	)
}

func hashInt(qname string, qtype uint16, do bool) string {
	h := fnv.New32()

	if do {
		h.Write(one)
	} else {
		h.Write(zero)
	}

	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, qtype)
	h.Write(b)

	for i := range qname {
		c := qname[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		h.Write([]byte{c})
	}

	return strconv.Itoa(int(h.Sum32()))
}

// ResponseWriter is a response writer that caches the reply message in Redis.
type ResponseWriter struct {
	dns.ResponseWriter
	state request.Request
	*Redis
	server      string
	DontUseHash bool
}

// WriteMsg implements the dns.ResponseWriter interface.
func (w *ResponseWriter) WriteMsg(res *dns.Msg) error {
	do := false
	mt, opt := response.Typify(res, w.now().UTC())
	if opt != nil {
		do = opt.Do()
	}

	// key returns empty string for anything we don't want to cache.
	cache, key := key(w.DontUseHash, res, mt, do)

	duration := w.pttl
	if mt == response.NameError || mt == response.NoData {
		duration = w.nttl
	}

	msgTTL := minMsgTTL(res, mt)
	if msgTTL < duration {
		duration = msgTTL
	}

	if cache && duration > 0 {
		if w.state.Match(res) {
			w.set(res, key, cache, mt, duration)
		} else {
			// Don't log it, but increment counter
			cacheDrops.WithLabelValues(w.server).Inc()
		}
	}

	// Apply capped TTL to this reply to avoid jarring TTL experience 1799 -> 8 (e.g.)
	ttl := uint32(duration.Seconds())
	for i := range res.Answer {
		res.Answer[i].Header().Ttl = ttl
	}
	for i := range res.Ns {
		res.Ns[i].Header().Ttl = ttl
	}
	for i := range res.Extra {
		if res.Extra[i].Header().Rrtype != dns.TypeOPT {
			res.Extra[i].Header().Ttl = ttl
		}
	}
	return w.ResponseWriter.WriteMsg(res)
}

func (w *ResponseWriter) set(m *dns.Msg, key string, cache bool, mt response.Type, duration time.Duration) {
	if !cache || duration == 0 {
		return
	}

	switch mt {
	case response.NoError, response.Delegation:
		fallthrough

	case response.NameError, response.NoData:
		if err := Add(w.pool, key, m, duration); err != nil {
			log.Debugf("Failed to add response to Redis cache: %s", err)

			redisErr.WithLabelValues(w.server).Inc()
		}

	case response.OtherError:
		// don't cache these
	default:
		log.Warningf("Redis called with unknown typification: %d", mt)
	}
}

// Write implements the dns.ResponseWriter interface.
func (w *ResponseWriter) Write(buf []byte) (int, error) {
	log.Warningf("Redis called with Write: not caching reply")
	n, err := w.ResponseWriter.Write(buf)
	return n, err
}

const (
	maxTTL      = 1 * time.Hour
	maxNTTL     = 30 * time.Minute
	failSafeTTL = 5 * time.Second

	// Success is the class for caching positive caching.
	Success = "success"
	// Denial is the class defined for negative caching.
	Denial = "denial"
)
