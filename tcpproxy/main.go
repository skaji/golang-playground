package main

import (
	"encoding/json"
	"fmt"
	"os"

	"inet.af/tcpproxy"
)

type Config struct {
	HTTPRoutes  [][3]string
	HTTPSRoutes [][3]string
	Routes      [][2]string
}

func loadConfig(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c *Config
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return c, nil
}

func main() {
	c, err := loadConfig(os.Args[1])
	if err != nil {
		panic(err)
	}

	var p tcpproxy.Proxy
	for _, e := range c.HTTPRoutes {
		p.AddHTTPHostRoute(e[0], e[1], tcpproxy.To(e[2]))
	}
	for _, e := range c.HTTPSRoutes {
		p.AddSNIRoute(e[0], e[1], tcpproxy.To(e[2]))
	}
	for _, e := range c.Routes {
		p.AddRoute(e[0], tcpproxy.To(e[1]))
	}
	if err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
