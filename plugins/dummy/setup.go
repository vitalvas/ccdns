package dummy

import (
	"errors"
	"fmt"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
)

const pluginName = "dummy"

var log = clog.NewWithPlugin(pluginName)

func init() { plugin.Register(pluginName, setup) }

func setup(c *caddy.Controller) error {
	p, err := parse(c)
	if err != nil {
		return plugin.Error(pluginName, err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		p.Next = next
		return p
	})

	return nil
}

func parse(c *caddy.Controller) (*Dummy, error) {
	p := New()

	for c.Next() {
		// args := c.RemainingArgs()
		// if len(args) > 0 {
		// 	fmt.Println(args)
		// }

		for c.NextBlock() {
			switch c.Val() {
			case "answer_minimisation":
				args := c.RemainingArgs()
				if len(args) == 1 {
					switch args[0] {
					case "none":
						p.AnswerMinimisation = AnswerMinimisationNone

					case "chash":
						p.AnswerMinimisation = AnswerMinimisationCHash

					default:
						return nil, errors.New("unknown answer_minimisation: " + args[0])
					}
				}

			case "answer":
				args := c.RemainingArgs()
				fmt.Println(args)

			default:
				return nil, c.ArgErr()
			}
		}
	}

	return p, nil
}
