package main

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/coremain"

	_ "github.com/coredns/coredns/plugin/acl"
	_ "github.com/coredns/coredns/plugin/any"
	_ "github.com/coredns/coredns/plugin/bind"
	_ "github.com/coredns/coredns/plugin/bufsize"
	_ "github.com/coredns/coredns/plugin/cache"
	_ "github.com/coredns/coredns/plugin/cancel"
	_ "github.com/coredns/coredns/plugin/debug"
	_ "github.com/coredns/coredns/plugin/errors"
	_ "github.com/coredns/coredns/plugin/file"
	_ "github.com/coredns/coredns/plugin/forward"
	_ "github.com/coredns/coredns/plugin/grpc"
	_ "github.com/coredns/coredns/plugin/header"
	_ "github.com/coredns/coredns/plugin/health"
	_ "github.com/coredns/coredns/plugin/hosts"
	_ "github.com/coredns/coredns/plugin/loadbalance"
	_ "github.com/coredns/coredns/plugin/log"
	_ "github.com/coredns/coredns/plugin/loop"
	_ "github.com/coredns/coredns/plugin/metadata"
	_ "github.com/coredns/coredns/plugin/metrics"
	_ "github.com/coredns/coredns/plugin/minimal"
	_ "github.com/coredns/coredns/plugin/nsid"
	_ "github.com/coredns/coredns/plugin/pprof"
	_ "github.com/coredns/coredns/plugin/reload"
	_ "github.com/coredns/coredns/plugin/template"
	_ "github.com/coredns/coredns/plugin/tls"
	_ "github.com/coredns/coredns/plugin/whoami"
	_ "github.com/miekg/redis"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// Directives are registered in the order they should be executed.
// Ordering is VERY important. Every plugin will feel the effects of all other plugin below
// (after) them during a request, but they must not care what plugin above them are doing.
var directives = []string{
	"metadata",
	"cancel",
	"tls",
	"reload",
	"nsid",
	"bufsize",
	"bind",
	"debug",
	"health",
	"pprof",
	"prometheus",
	"metrics",
	"errors",
	"log",
	"acl",
	"any",
	"loadbalance",
	"cache",
	"redisc",
	"header",
	"minimal",
	"template",
	"hosts",
	"file",
	"loop",
	"forward",
	"grpc",
	"whoami",
}

func init() {
	dnsserver.Directives = directives
	caddy.AppVersion = version
}

func main() {
	coremain.Run()
}
