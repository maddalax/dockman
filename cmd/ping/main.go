package main

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"
)

func main() {
	localIP, ipNet := getLocalIPAndSubnet()
	if localIP == "" || ipNet == nil {
		fmt.Println("Unable to determine local IP or subnet. Ensure your network connection is active.")
		return
	}

	ports := []int{22, 443, 3389, 5900, 80, 6380} // List of ports to scan

	fmt.Println("Scanning for devices on the network...")
	for ip := ipNet.IP.Mask(ipNet.Mask); ipNet.Contains(ip); incIP(ip) {
		for _, port := range ports {
			go func(ip string, port int) {
				if isPortOpen(ip, port) {
					deviceName := getDeviceName(ip)
					osInfo := detectOS(ip)
					mac, manufacturer := getMACAndManufacturer(ip)
					uptime := getUptime(ip)

					fmt.Printf("mac: %s, manufacturer: %s\n", mac, manufacturer)
					fmt.Printf("Device found: %s (port %d open)\n", ip, port)
					fmt.Printf("Hostname: %s, OS: %s, Manufacturer: %s, Uptime: %s\n", deviceName, osInfo, manufacturer, uptime)
				}
			}(ip.String(), port)
		}
	}

	// Wait to give the goroutines time to finish
	time.Sleep(30 * time.Second)
	fmt.Println("Scanning complete.")
}

// getLocalIPAndSubnet retrieves the local IP and subnet
func getLocalIPAndSubnet() (string, *net.IPNet) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Error getting network interfaces:", err)
		return "", nil
	}

	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String(), ipNet
		}
	}
	return "", nil
}

// incIP increments the IP address by one
func incIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// isPortOpen checks if a TCP port is open on the given IP address
func isPortOpen(ip string, port int) bool {
	address := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// getDeviceName performs a reverse DNS lookup to get the device name
func getDeviceName(ip string) string {
	names, err := net.LookupAddr(ip)
	if err != nil || len(names) == 0 {
		return "Unknown"
	}
	return names[0]
}

// detectOS attempts to detect the operating system using simple fingerprinting
func detectOS(ip string) string {
	// Use external tools like Nmap for more accurate OS detection
	out, err := exec.Command("nmap", "-O", ip).Output()
	if err != nil {
		return "Unknown"
	}
	output := string(out)
	if strings.Contains(output, "Windows") {
		return "Windows"
	} else if strings.Contains(output, "Linux") {
		return "Linux"
	} else if strings.Contains(output, "macOS") {
		return "macOS"
	}
	return "Unknown"
}

// getMACAndManufacturer retrieves the MAC address and manufacturer info
func getMACAndManufacturer(ip string) (string, string) {
	//handle, err := pcap.OpenLive("eth0", 1600, true, pcap.BlockForever)
	//if err != nil {
	//	return "Unknown", "Unknown"
	//}
	//defer handle.Close()
	//
	//packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	//for packet := range packetSource.Packets() {
	//	ethLayer := packet.Layer(gopacket.LayerTypeEthernet)
	//	if ethLayer != nil {
	//		eth, _ := ethLayer.(*gopacket.EthernetLayer)
	//		mac := eth.SrcMAC.String()
	//		manufacturer := lookupManufacturer(mac)
	//		return mac, manufacturer
	//	}
	//}
	return "Unknown", "Unknown"
}

// lookupManufacturer performs an OUI lookup for the MAC address
func lookupManufacturer(mac string) string {
	// You can use a local database or an API for OUI lookup
	return "Manufacturer Info"
}

// getUptime attempts to get the device uptime using SNMP (requires an SNMP library)
func getUptime(ip string) string {
	// Use an SNMP library to query the device for uptime
	// Placeholder for actual SNMP implementation
	return "Uptime Info"
}
