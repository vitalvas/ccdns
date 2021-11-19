package psqlx

import (
	"context"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/jackc/pgx/v4/pgxpool"
)

func init() {
	caddy.RegisterPlugin(Name, caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	px := PSQLX{}

	var dburl string
	var err error

	for c.Next() {
		for c.NextBlock() {
			switch c.Val() {
			case "db":
				args := c.RemainingArgs()

				if len(args) == 0 {
					return plugin.Error(Name, c.ArgErr())
				}

				if len(args) != 2 {
					return plugin.Error(Name, c.Errf("db requires 2 arguments"))
				}

				switch args[0] {
				case "url":
					dburl = args[1]

				default:
					return plugin.Error(Name, c.Errf("unknown property '%s'", args[0]))
				}

			default:
				if c.Val() != "}" {
					return plugin.Error(Name, c.Errf("unknown property '%s'", c.Val()))
				}
			}
		}
	}

	px.dbPool, err = pgxpool.Connect(context.Background(), dburl)
	if err != nil {
		return plugin.Error(Name, err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		px.Next = next
		return px
	})

	return nil
}
