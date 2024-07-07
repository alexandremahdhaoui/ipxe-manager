package types

import (
	"encoding"
	"encoding/hex"
	"errors"
	"net"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
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

var errUnsupportedParameterType = errors.New("unsupported parameter type")

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
	Uuid         *uuid.UUID // UUID
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

func NewIpxeParamsFromContext(c echo.Context) (IpxeParams, error) {
	ipxeParams := IpxeParams{}
	var err error

	if ipxeParams.Mac, err = getParam[hexa](c, Mac); err != nil {
		return IpxeParams{}, err
	}
	if ipxeParams.BusType, err = getParam[string](c, BusType); err != nil {
		return IpxeParams{}, err
	}
	if ipxeParams.BusLoc, err = getParam[uint32](c, BusLoc); err != nil {
		return IpxeParams{}, err
	}
	if ipxeParams.BusID, err = getParam[hexa](c, BusID); err != nil {
		return IpxeParams{}, err
	}
	if ipxeParams.Chip, err = getParam[string](c, Chip); err != nil {
		return IpxeParams{}, err
	}
	if ipxeParams.Ssid, err = getParam[string](c, Ssid); err != nil {
		return IpxeParams{}, err
	}
	if ipxeParams.ActiveScan, err = getParam[int8](c, ActiveScan); err != nil {
		return IpxeParams{}, err
	}
	if ipxeParams.Key, err = getParam[string](c, Key); err != nil {
		return IpxeParams{}, err
	}

	// IPv4 settings

	if ipxeParams.Ip, err = getParam[net.IP](c, Ip); err != nil {
		return IpxeParams{}, err
	}
	if ipxeParams.Netmask, err = getParam[net.IP](c, Netmask); err != nil {
		return IpxeParams{}, err
	}
	if ipxeParams.Gateway, err = getParam[net.IP](c, Gateway); err != nil {
		return IpxeParams{}, err
	}
	if ipxeParams.Dns, err = getParam[net.IP](c, Dns); err != nil {
		return IpxeParams{}, err
	}
	if ipxeParams.Domain, err = getParam[string](c, Domain); err != nil {
		return IpxeParams{}, err
	}

	// Boot settings

	if ipxeParams.Filename, err = getParam[string](c, Filename); err != nil {
		return IpxeParams{}, err
	}
	if ipxeParams.NextServer, err = getParam[net.IP](c, NextServer); err != nil {
		return IpxeParams{}, err
	}
	if ipxeParams.RootPath, err = getParam[string](c, RootPath); err != nil {
		return IpxeParams{}, err
	}
	if ipxeParams.SanFilename, err = getParam[string](c, SanFilename); err != nil {
		return IpxeParams{}, err
	}
	if ipxeParams.InitiatorIqn, err = getParam[string](c, InitiatorIqn); err != nil {
		return IpxeParams{}, err
	}
	if ipxeParams.KeepSan, err = getParam[int8](c, KeepSan); err != nil {
		return IpxeParams{}, err
	}
	if ipxeParams.SkipSanBoot, err = getParam[int8](c, SkipSanBoot); err != nil {
		return IpxeParams{}, err
	}

	// Host settings

	if ipxeParams.Hostname, err = getParam[string](c, Hostname); err != nil {
		return IpxeParams{}, err
	}
	if ipxeParams.Uuid, err = getParam[uuid.UUID](c, Uuid); err != nil {
		return IpxeParams{}, err
	}
	if ipxeParams.UserClass, err = getParam[string](c, UserClass); err != nil {
		return IpxeParams{}, err
	}
	if ipxeParams.Manufacturer, err = getParam[string](c, Manufacturer); err != nil {
		return IpxeParams{}, err
	}
	if ipxeParams.Product, err = getParam[string](c, Product); err != nil {
		return IpxeParams{}, err
	}
	if ipxeParams.Serial, err = getParam[string](c, Serial); err != nil {
		return IpxeParams{}, err
	}
	if ipxeParams.Asset, err = getParam[string](c, Asset); err != nil {
		return IpxeParams{}, err
	}

	// Authentication settings

	if ipxeParams.Username, err = getParam[string](c, Username); err != nil {
		return IpxeParams{}, err
	}
	if ipxeParams.Password, err = getParam[string](c, Password); err != nil {
		return IpxeParams{}, err
	}
	if ipxeParams.ReverseUsername, err = getParam[string](c, ReverseUsername); err != nil {
		return IpxeParams{}, err
	}
	if ipxeParams.ReversePassword, err = getParam[string](c, ReversePassword); err != nil {
		return IpxeParams{}, err
	}

	// Miscellaneous settings

	if ipxeParams.Crosscert, err = getParam[string](c, Crosscert); err != nil {
		return IpxeParams{}, err
	}

	if ipxeParams.Trust, err = getParam[hexa](c, Trust); err != nil {
		return IpxeParams{}, err
	}

	if ipxeParams.Cert, err = getParam[hexa](c, Cert); err != nil {
		return IpxeParams{}, err
	}

	if ipxeParams.Privkey, err = getParam[hexa](c, Privkey); err != nil {
		return IpxeParams{}, err
	}

	if ipxeParams.Buildarch, err = getParam[string](c, Buildarch); err != nil {
		return IpxeParams{}, err
	}

	if ipxeParams.Cpumodel, err = getParam[string](c, Cpumodel); err != nil {
		return IpxeParams{}, err
	}

	if ipxeParams.Cpuvendor, err = getParam[string](c, Cpuvendor); err != nil {
		return IpxeParams{}, err
	}

	if ipxeParams.DhcpServer, err = getParam[net.IP](c, DhcpServer); err != nil {
		return IpxeParams{}, err
	}

	if ipxeParams.Keymap, err = getParam[string](c, Keymap); err != nil {
		return IpxeParams{}, err
	}

	if ipxeParams.Memsize, err = getParam[int32](c, Memsize); err != nil {
		return IpxeParams{}, err
	}

	if ipxeParams.Platform, err = getParam[string](c, Platform); err != nil {
		return IpxeParams{}, err
	}

	if ipxeParams.Priority, err = getParam[int8](c, Priority); err != nil {
		return IpxeParams{}, err
	}

	if ipxeParams.Scriptlet, err = getParam[string](c, Scriptlet); err != nil {
		return IpxeParams{}, err
	}

	if ipxeParams.Syslog, err = getParam[net.IP](c, Syslog); err != nil {
		return IpxeParams{}, err
	}

	if ipxeParams.Syslogs, err = getParam[string](c, Syslogs); err != nil {
		return IpxeParams{}, err
	}

	if ipxeParams.Sysmac, err = getParam[hexa](c, Sysmac); err != nil {
		return IpxeParams{}, err
	}

	if ipxeParams.Unixtime, err = getParam[uint32](c, Unixtime); err != nil {
		return IpxeParams{}, err
	}

	if ipxeParams.UseCached, err = getParam[uint8](c, UseCached); err != nil {
		return IpxeParams{}, err
	}

	if ipxeParams.Version, err = getParam[string](c, Version); err != nil {
		return IpxeParams{}, err
	}

	if ipxeParams.Vram, err = getParam[[]byte](c, Vram); err != nil {
		return IpxeParams{}, err
	}

	return ipxeParams, nil
}

func getParam[T any](c echo.Context, key string) (*T, error) {
	out := new(T)
	if !c.QueryParams().Has(key) {
		return nil, nil
	}

	s := c.QueryParam(key)

	if v, ok := any(out).(encoding.TextUnmarshaler); ok {
		if err := v.UnmarshalText([]byte(s)); err != nil {
			return nil, err // wrap error
		}

		return out, nil
	}

	switch any(out).(type) {
	case string:
		return any(&s).(*T), nil
	case int32:
		i, err := strconv.Atoi(s)
		if err != nil {
			return nil, err // TODO: wrap this err
		}

		return any(Ptr(int32(i))).(*T), nil
	case uint32:

		i, err := strconv.Atoi(s)
		if err != nil {
			return nil, err // TODO: wrap this err
		}

		return any(Ptr(uint32(i))).(*T), nil
	case int8:

		i, err := strconv.Atoi(s)
		if err != nil {
			return nil, err // TODO: wrap this err
		}

		return any(Ptr(int8(i))).(*T), nil
	case uint8:
		i, err := strconv.Atoi(s)
		if err != nil {
			return nil, err // TODO: wrap this err
		}

		return any(Ptr(uint8(i))).(*T), nil
	default:
		return nil, errUnsupportedParameterType
	}
}

func Ptr[T any](v T) *T {
	return &v
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

type IpxeSelectors struct {
	UUID      uuid.UUID
	Buildarch string
}

func NewIpxeSelectors(params IpxeParams) (IpxeSelectors, error) {
	if params.Uuid == nil {
		return IpxeSelectors{}, errors.New("TODO") // TODO: define this err.
	}

	if params.Buildarch == nil {
		return IpxeSelectors{}, errors.New("TODO") // TODO: define this err.
	}

	return IpxeSelectors{
		UUID:      *params.Uuid,
		Buildarch: *params.Buildarch,
	}, nil
}

func NewIpxeSelectorsFromContext(c echo.Context) (IpxeSelectors, error) {
	params, err := NewIpxeParamsFromContext(c)
	if err != nil {
		return IpxeSelectors{}, err // TODO: wrap
	}

	selectors, err := NewIpxeSelectors(params)
	if err != nil {
		return IpxeSelectors{}, err // TODO: wrap
	}

	return selectors, nil
}
