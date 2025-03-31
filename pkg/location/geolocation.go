package location

import (
	"encoding/csv"
	"github.com/yl2chen/cidranger"
	"io"
	"net"
	"os"
)

func Load(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	_, err = reader.Read()
	if err != nil {
		return err
	}

	ranger = cidranger.NewPCTrieRanger()
	line := 2
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			line++
			continue
		}

		if len(record) < 6 {
			line++
			continue
		}

		_, network, err := net.ParseCIDR(record[0])
		if err != nil {
			line++
			continue
		}

		geoData := GeoData{
			Asn:         record[1],
			Org:         record[2],
			CountryCode: record[3],
			Latitude:    record[4],
			Longitude:   record[5],
		}

		ranger.Insert(GeoEntry{
			network: *network,
			data:    geoData,
		})

		line++
	}
	return nil
}

func FindGeolocation(ip string) *GeoInfo {
	cache1Lock.RLock()

	if geoInfo, found := cache1[ip]; found {
		cache1Lock.RUnlock()
		return &geoInfo
	}

	cache1Lock.RUnlock()

	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return nil
	}

	entries, err := ranger.ContainingNetworks(ipAddr)
	if err != nil || len(entries) == 0 {
		return nil
	}

	entry := entries[0].(GeoEntry)
	geoData := entry.data

	geoInfo := GeoInfo{
		CountryCode: geoData.CountryCode,
		Latitude:    geoData.Latitude,
		Longitude:   geoData.Longitude,
		Asn:         geoData.Asn,
		Org:         geoData.Org,
	}

	cache1Lock.Lock()
	cache1[ip] = geoInfo
	cache1Lock.Unlock()

	return &geoInfo
}
