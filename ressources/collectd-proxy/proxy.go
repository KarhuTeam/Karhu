package collectdproxy

import (
	"github.com/gotoolz/env"
	"github.com/karhuteam/karhu/models"
	"github.com/patrickmn/go-cache"
	"log"
	"net"
	"time"
)

var proxy *net.UDPConn
var target *net.UDPConn

// host cache
// Store allowed host with 5 min ttl
// Score disallowed host with 30s ttl
var hostCache *cache.Cache

type cacheEntry struct {
	T       time.Time
	Allowed bool
}

func (e cacheEntry) Expired() bool {

	var t time.Duration
	switch e.Allowed {
	case true:
		t = time.Minute * 5
	case false:
		t = time.Second * 30
	}

	return e.T.Add(t).Before(time.Now())
}

func init() {

	hostCache = cache.New(5*time.Minute, 30*time.Second)

	if err := setup(); err != nil {
		log.Fatalln("collectproxy:", err)
	}

	go runProxy()
}

func setup() error {

	// Resolve bind address
	saddr, err := net.ResolveUDPAddr("udp", env.GetDefault("COLLECTD_PROXY_BIND", "127.0.0.1:25827"))
	if err != nil {
		return err
	}

	// Listen udp socket
	proxy, err = net.ListenUDP("udp", saddr)
	if err != nil {
		return err
	}

	// Resolve target address
	taddr, err := net.ResolveUDPAddr("udp", env.GetDefault("COLLECTD_PROXY_TARGET", "127.0.0.1:25826"))
	if err != nil {
		return err
	}

	target, err = net.DialUDP("udp", nil, taddr)
	if err != nil {
		return err
	}

	return nil
}

func isAllowed(ip string) (bool, error) {

	c, found := hostCache.Get(ip)
	if !found || found && c.(cacheEntry).Expired() {

		allowed, err := isAllowedNoCache(ip)
		if err != nil {
			return false, err
		}

		entry := cacheEntry{
			Allowed: allowed,
			T:       time.Now(),
		}

		hostCache.Set(ip, entry, cache.DefaultExpiration)
		c = entry
	}

	return c.(cacheEntry).Allowed, nil
}

func isAllowedNoCache(ip string) (bool, error) {

	node, err := models.NodeMapper.FetchOneIP(ip)
	if err != nil {
		return false, err
	}

	return node != nil, nil
}

func runProxy() {

	var buffer [1500]byte
	for {

		n, cliaddr, err := proxy.ReadFromUDP(buffer[0:])
		if err != nil {
			log.Println("collectdproxy: ReadFromUDP:", err)
			continue
		}

		// log.Printf("Read %d from client %s\n", n, cliaddr.IP.String())

		allowed, err := isAllowed(cliaddr.IP.String())
		if err != nil {
			log.Println("collectdproxy: isAllowed:", err)
			continue
		}

		if !allowed {
			log.Println("collectdproxy: disalowed:", cliaddr.IP.String())
			continue
		}

		_, err = target.Write(buffer[0:n])
		if err != nil {
			log.Println("collectdproxy: Write:", err)
			continue
		}

		// log.Printf("Writed %d from client %s\n", n, cliaddr.String())
	}
}
