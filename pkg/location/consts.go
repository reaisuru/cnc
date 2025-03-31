package location

import (
	"github.com/yl2chen/cidranger"
	"net"
	"sync"
)

var (
	ranger     cidranger.Ranger
	cache1     = make(map[string]GeoInfo)
	cache1Lock sync.RWMutex
)

type GeoData struct {
	Asn         string
	Org         string
	CountryCode string
	Latitude    string
	Longitude   string
}

type GeoInfo struct {
	CountryCode string
	Latitude    string
	Longitude   string
	Asn         string
	Org         string
}

type GeoEntry struct {
	network net.IPNet
	data    GeoData
}

func (ge GeoEntry) Network() net.IPNet {
	return ge.network
}
