package main

import (
	"flag"
	"net"

	"github.com/gorcon/rcon"
	"github.com/jszwec/csvutil"
	mp "github.com/mackerelio/go-mackerel-plugin"
)

func main() {
	plugin := new(Plugin)

	flag.StringVar(&plugin.host, "host", "localhost", "hostname")
	flag.StringVar(&plugin.port, "port", "25575", "RCON port")
	flag.StringVar(&plugin.password, "password", "", "administrator passowrd")
	flag.Parse()

	mackerelPlugin := mp.NewMackerelPlugin(plugin)
	mackerelPlugin.Run()
}

type Plugin struct {
	host     string
	port     string
	password string
}

var _ mp.PluginWithPrefix = new(Plugin)

type Player struct {
	Name      string `csv:"name"`
	PlayerUID int64  `csv:"playeruid"`
	SteamID   int64  `csv:"steamid"`
}

func (p *Plugin) FetchMetrics() (map[string]float64, error) {
	address := net.JoinHostPort(p.host, p.port)
	conn, err := rcon.Dial(address, p.password)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	response, err := conn.Execute("ShowPlayers")
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
