package types

import (
	"encoding/hex"
	"net"
	"strings"

	"github.com/google/uuid"
)

// -------------------------------------------------- PARAMETERS ---------------------------------------------------- //

const (
	Mac        = "mac"         //	MAC address
	BusType    = "bustype"     // Bus type
	BusLoc     = "busloc"      // Bus location
	BusID      = "busid"       // Bus ExposedConfigID
	Chip       = "chip"        // Chip type
	Ssid       = "ssid"        // Wireless SSID
	ActiveScan = "active-scan" // Actively scan for wireless orks
	Key        = "key"         // Wireless encryption key

	// IPv4 settings

	Ip      = "ip"      // IP address
	Netmask = "netmask" // Subnet mask
	Gateway = "gateway" // Default gateway
	Dns     = "dns"     // DNS server
	Domain  = "domain"  // DNS domain

	// Boot settings

	Filename     = "filename"      // Boot filename
	NextServer   = "next-server"   // TFTP server
	RootPath     = "root-path"     // SAN root path
	SanFilename  = "scan-filename" // SAN filename
	InitiatorIqn = "initiator-iqn" // iSCSI initiator name
	KeepSan      = "keep-san"      // Preserve SAN connection
	SkipSanBoot  = "skip-san-boot" // Do not boot from SAN device

	// Host settings

	Hostname     = "hostname"     // Host name
	Uuid         = "uuid"         // UUID
	UserClass    = "user-class"   // DHCP user class
	Manufacturer = "manufacturer" // Manufacturer
	Product      = "product"      // Product name
	Serial       = "serial"       // Serial number
	Asset        = "asset"        // Asset tag

	// Authentication settings

	Username        = "username"         // User name
	Password        = "password"         // Password
	ReverseUsername = "reverse-username" // Reverse user name
	ReversePassword = "reverse-password" // Reverse password

	// Cryptography settings

	Crosscert = "crosscert" // Cross-signed certificate source
	Trust     = "trust"     // Trusted root certificate fingerprints
	Cert      = "cert"      // Client certificate
	Privkey   = "privkey"   // Client private key

	// Miscellaneous settings

	Buildarch  = "buildarch"   // Build architecture
	Cpumodel   = "cpumodel"    // CPU model
	Cpuvendor  = "cpuvendor"   // CPU vendor
	DhcpServer = "dhcp-server" // DHCP server
	Keymap     = "keymap"      // Keyboard layout
	Memsize    = "memsize"     // Memory size
	Platform   = "platform"    // Firmware platform
	Priority   = "priority"    // Settings priority
	Scriptlet  = "scriptlet"   // Boot scriptlet
	Syslog     = "syslog"      // Syslog server
	Syslogs    = "syslogs"     // Encrypted syslog server
	Sysmac     = "sysmac"      // System MAC address
	Unixtime   = "unixtime"    // Seconds since the Epoch
	UseCached  = "use-cached"  // Use cached settings
	Version    = "version"     // iPXE version
	Vram       = "vram"        // Video RAM contents
)

// --- PARAMS --- //

type IpxeParams struct {
	Mac        *hexa   //	MAC address
	BusType    *string // Bus type
	BusLoc     *uint32 // Bus location
	BusID      *hexa   // Bus ExposedConfigID
	Chip       *string // Chip type
	Ssid       *string // Wireless SSID
	ActiveScan *int8   // Actively scan for wireless orks
	Key        *string // Wireless encryption key

	// IPv4 settings

	Ip      *net.IP // IP address
	Netmask *net.IP // Subnet mask
	Gateway *net.IP // Default gateway
	Dns     *net.IP // DNS server
	Domain  *string // DNS domain

	// Boot settings

	Filename     *string // Boot filename
	NextServer   *net.IP // TFTP server
	RootPath     *string // SAN root path
	SanFilename  *string // SAN filename
	InitiatorIqn *string // iSCSI initiator name
	KeepSan      *int8   // Preserve SAN connection
	SkipSanBoot  *int8   // Do not boot from SAN device

	// Host settings

	Hostname     *string    // Host name
	UUID         *uuid.UUID // UUID
	UserClass    *string    // DHCP user class
	Manufacturer *string    // Manufacturer
	Product      *string    // Product name
	Serial       *string    // Serial number
	Asset        *string    // Asset tag

	// Authentication settings

	Username        *string // User name
	Password        *string // Password
	ReverseUsername *string // Reverse user name
	ReversePassword *string // Reverse password

	// Cryptography settings

	Crosscert *string // Cross-signed certificate source
	Trust     *hexa   // Trusted root certificate fingerprints
	Cert      *hexa   // Client certificate
	Privkey   *hexa   // Client private key

	// Miscellaneous settings

	Buildarch  *string // Build architecture
	Cpumodel   *string // CPU model
	Cpuvendor  *string // CPU vendor
	DhcpServer *net.IP // DHCP server
	Keymap     *string // Keyboard layout
	Memsize    *int32  // Memory size
	Platform   *string // Firmware platform
	Priority   *int8   // Settings priority
	Scriptlet  *string // Boot scriptlet
	Syslog     *net.IP // Syslog server
	Syslogs    *string // Encrypted syslog server
	Sysmac     *hexa   // System MAC address
	Unixtime   *uint32 // Seconds since the Epoch
	UseCached  *uint8  // Use cached settings
	Version    *string // iPXE version
	Vram       *[]byte // Video RAM contents
}

type hexa []byte

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
func (b *hexa) UnmarshalText(text []byte) error {
	*b = make(hexa, 0)

	for _, s := range strings.Split(string(text), ":") {
		decoded, err := hex.DecodeString(s)
		if err != nil {
			return err // TODO: write this err.
		}

		*b = append(*b, decoded...)
	}

	return nil
}

// ------------------------------------------------ LABEL SELECTORS ------------------------------------------------- //

type IPXESelectors struct {
	Buildarch string
	UUID      uuid.UUID
}
