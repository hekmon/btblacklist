package ripe

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strings"
)

func extractRange(payload result) (ranges []Range, err error) {
	// How many results ?
	nbInets, nbRoutes, err := nbResults(payload)
	if err != nil || (nbInets == 0 && nbRoutes == 0) {
		return
	}
	ranges = make([]Range, 0, nbInets+nbRoutes)
	// Search for inetnums
	if err = searchInets(payload, &ranges); err != nil {
		err = fmt.Errorf("error while searching for inetnums: %v", err)
		return
	}
	// Search for routes
	if err = searchRoutes(payload, &ranges); err != nil {
		err = fmt.Errorf("error while searching for routes: %v", err)
		return
	}
	// Done
	return
}

func nbResults(payload result) (nbInets, nbRoutes int, err error) {
	var found bool
	for _, list := range payload.Lists {
		if list.Name != "facet_counts" {
			continue
		}
		for _, listList := range list.Lists {
			if listList.Name != "facet_fields" {
				continue
			}
			for _, listListList := range listList.Lists {
				if listListList.Name != "object-type" {
					continue
				}
				found = true
				if listListList.Ints == nil {
					err = fmt.Errorf("object-type ints map is nil")
					return
				}
				// inets
				inetnumList, ok := listListList.Ints["inetnum"]
				if ok {
					for _, inetnum := range inetnumList {
						nbInets += inetnum
					}
				}
				// routes
				routeList, ok := listListList.Ints["route"]
				if ok {
					for _, routenum := range routeList {
						nbRoutes += routenum
					}
				}
			}
		}
	}
	if !found {
		err = errors.New("can't find 'facet_counts > facet_fields > object-type'")
	}
	return
}

func searchInets(payload result, ranges *[]Range) (err error) {
	var (
		inetnum, netname, desc []string
		ok                     bool
	)
	for _, doc := range payload.Result.Docs {
		// Only if inetnum
		if inetnum, ok = doc.Strings["inetnum"]; !ok {
			continue
		}
		if len(inetnum) != 1 {
			err = fmt.Errorf("inetnum has %d values (expecting only 1): %v", len(inetnum), inetnum)
			return
		}
		// With which netname ?
		if netname, ok = doc.Strings["netname"]; !ok {
			err = fmt.Errorf("can't find a 'netname' for 'inetnum': %v", inetnum)
			return
		}
		if len(netname) != 1 {
			err = fmt.Errorf("inetnum '%s' has multiples 'netname': %v", inetnum[0], netname)
			return
		}
		// With which desc ?
		if desc, ok = doc.Strings["descr"]; !ok {
			err = fmt.Errorf("can't find a 'descr' for 'inetnum': %v", inetnum)
			return
		}
		if len(desc) != 1 {
			err = fmt.Errorf("inetnum '%s' has multiples 'descr': %v", inetnum[0], desc)
			return
		}
		// Save range
		*ranges = append(*ranges, Range{
			Name:  fmt.Sprintf("%s (%s)", netname[0], desc[0]),
			Range: strings.Replace(inetnum[0], " ", "", -1),
		})
	}
	return
}

func searchRoutes(payload result, ranges *[]Range) (err error) {
	var (
		route, desc []string
		network     *net.IPNet
		broadcastIP net.IP
		ok          bool
	)
	for _, doc := range payload.Result.Docs {
		// Only if inetnum
		if route, ok = doc.Strings["route"]; !ok {
			continue
		}
		if len(route) != 1 {
			err = fmt.Errorf("route has %d values (expecting only 1): %v", len(route), route)
			return
		}
		// With which desc ?
		if desc, ok = doc.Strings["descr"]; !ok {
			err = fmt.Errorf("can't find a 'descr' for 'inetnum': %v", route)
			return
		}
		if len(desc) != 1 {
			err = fmt.Errorf("inetnum '%s' has multiples 'descr': %v", route[0], desc)
			return
		}
		// Get range
		if _, network, err = net.ParseCIDR(route[0]); err != nil {
			err = fmt.Errorf("can't parse route '%s' as network: %v", route[0], err)
			return
		}
		if broadcastIP, err = lastAddr(network); err != nil {
			err = fmt.Errorf("can't get last address of %s: %v", network, err)
			return
		}
		// Save range
		*ranges = append(*ranges, Range{
			Name:  desc[0],
			Route: route[0],
			Range: fmt.Sprintf("%s-%s", network.IP, broadcastIP),
		})
	}
	return
}

// https://stackoverflow.com/questions/36166791/how-to-get-broadcast-address-of-ipv4-net-ipnet
func lastAddr(n *net.IPNet) (ip net.IP, err error) { // works when the n is a prefix, otherwise...
	nv4 := n.IP.To4()
	if nv4 == nil {
		err = errors.New("does not support IPv6 addresses")
		return
	}
	ip = make(net.IP, len(nv4))
	binary.BigEndian.PutUint32(ip, binary.BigEndian.Uint32(n.IP.To4())|^binary.BigEndian.Uint32(net.IP(n.Mask).To4()))
	return
}
