package main

import (
	"flag"
	"log"
	"net"
	"strings"
	"time"

	"github.com/gorcon/rcon"
	"github.com/jszwec/csvutil"
	mp "github.com/mackerelio/go-mackerel-plugin"
)

func main() {
	plugin := new(Plugin)
	var (
		timeoutStr string
		err        error
	)

	flag.StringVar(&plugin.host, "host", "localhost", "hostname")
	flag.StringVar(&plugin.port, "port", "25575", "RCON port")
	flag.StringVar(&plugin.password, "password", "", "administrator passowrd")
	flag.StringVar(&timeoutStr, "timeout", "5s", "dial timeout seconds")
	flag.Parse()

	plugin.timeout, err = time.ParseDuration(timeoutStr)
	if err != nil {
		log.Fatalln(err)
	}
	mackerelPlugin := mp.NewMackerelPlugin(plugin)
	mackerelPlugin.Run()
}

type Plugin struct {
	host     string
	port     string
	password string
	timeout  time.Duration
}

var _ mp.PluginWithPrefix = new(Plugin)

type Player struct {
	Name      string `csv:"name"`
	PlayerUID int64  `csv:"playeruid"`
	SteamID   int64  `csv:"steamid"`
}

func (p *Plugin) FetchMetrics() (map[string]float64, error) {
	response, err := p.getShowPlayers()
	if err != nil {
		return nil, err
	}
	var players []Player
	if err := csvutil.Unmarshal([]byte(response), &players); err != nil {
		return nil, err
	}

	metrics := make(map[string]float64, 1)
	// why is the key not "players.num"?
	metrics["num"] = float64(len(players))
	return metrics, nil
}

func (p *Plugin) getShowPlayers() (string, error) {
	address := net.JoinHostPort(p.host, p.port)
	conn, err := rcon.Dial(address, p.password, rcon.SetDialTimeout(p.timeout))
	if err != nil {
		return "", err
	}
	defer conn.Close()

	// palworld's rcon server sometimes times out accidentally
	// nolint
	response, _ := conn.Execute("ShowPlayers")
	response = strings.ReplaceAll(response, "\x00", "")
	return response, nil
}

func (p *Plugin) GraphDefinition() map[string]mp.Graphs {
	return map[string]mp.Graphs{
		"players": {
			Label: "PalServer Players",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "num", Label: "Number"},
			},
		},
	}
}

func (p *Plugin) MetricKeyPrefix() string {
	return "palworld"
}
