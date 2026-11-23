package discover

import "net"

func networks() ([]*net.IPNet, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var nets []*net.IPNet
	for _, iface := range ifaces {
		usable := iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagBroadcast != 0
		if !usable || iface.Flags&(net.FlagLoopback|net.FlagPointToPoint) != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok || ipnet.IP.To4() == nil {
				continue
			}
			nets = append(nets, &net.IPNet{IP: ipnet.IP.To4(), Mask: ipnet.Mask[len(ipnet.Mask)-4:]})
		}
	}
	return nets, nil
}

func broadcast(n *net.IPNet) net.IP {
	ip := make(net.IP, 4)
	for i := range ip {
		ip[i] = n.IP[i] | ^n.Mask[i]
	}
	return ip
}
