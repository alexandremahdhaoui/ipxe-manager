package controller

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"text/template"

	"github.com/alexandremahdhaoui/ipxer/internal/adapter"
	"github.com/alexandremahdhaoui/ipxer/internal/types"
)

var (
	ErrIPXEFindProfileAndRender = errors.New("finding and rendering ipxe profile")

	errFallbackToDefaultAssignment = errors.New("fallback to default assignment")
	errSelectingAssignment         = errors.New("selecting assignment")
	errTemplatingIPXEProfile       = errors.New("templating ipxe profile")

	fmtCannotSelectAssignmentWithSelectors = "cannot select assignment with selectors: uuid=%q & buildarch=%q"
)

// ---------------------------------------------------- INTERFACES -------------------------------------------------- //

type IPXE interface {
	FindProfileAndRender(ctx context.Context, selectors types.IpxeSelectors) ([]byte, error)
	Boostrap() []byte
}

// --------------------------------------------------- CONSTRUCTORS ------------------------------------------------- //

func NewIPXE(assignment adapter.Assignment, profile adapter.Profile, mux ResolveTransformerMux) IPXE {
	return &ipxe{
		assignment: assignment,
		profile:    profile,
		mux:        mux,
	}
}

// -------------------------------------------------------- IPXE ---------------------------------------------------- //

type ipxe struct {
	assignment adapter.Assignment
	profile    adapter.Profile
	mux        ResolveTransformerMux

	bootstrap []byte
}

// -------------------------------------------------------- FindProfileAndRender ------------------------------------ //

func (svc *ipxe) FindProfileAndRender(ctx context.Context, selectors types.IpxeSelectors) ([]byte, error) {
	assignment, err := svc.assignment.FindProfileBySelectors(ctx, selectors)
	if errors.Is(err, adapter.ErrAssignmentNotFound) {
		// fallback to default profile
		defaultAssignment, defaultErr := svc.assignment.FindDefaultProfile(ctx, selectors.Buildarch)
		if defaultErr != nil {
			return nil, errors.Join(defaultErr,
				fmt.Errorf(fmtCannotSelectAssignmentWithSelectors, selectors.UUID, selectors.Buildarch),
				errFallbackToDefaultAssignment, errSelectingAssignment, ErrIPXEFindProfileAndRender)
		}

		assignment = defaultAssignment
	} else if err != nil {
		return nil, errors.Join(err, errSelectingAssignment, ErrIPXEFindProfileAndRender)
	}

	p, err := svc.profile.Get(ctx, assignment)
	if err != nil {
		return nil, errors.Join(err, ErrIPXEFindProfileAndRender)
	}

	data, err := svc.mux.ResolveAndTransformBatch(ctx, p.AdditionalContent)
	if err != nil {
		return nil, errors.Join(err, ErrIPXEFindProfileAndRender)
	}

	out, err := templateIPXEProfile(p.IPXETemplate, data)
	if err != nil {
		return nil, errors.Join(err, ErrIPXEFindProfileAndRender)
	}

	return out, nil
}

func templateIPXEProfile(ipxeTemplate string, data map[string][]byte) ([]byte, error) {
	tpl, err := template.New("").Parse(ipxeTemplate)
	if err != nil {
		return nil, errors.Join(err, errTemplatingIPXEProfile)
	}

	buf := bytes.NewBuffer(make([]byte, 0))
	if err := tpl.Execute(buf, data); err != nil {
		return nil, errors.Join(err, errTemplatingIPXEProfile)
	}

	return buf.Bytes(), nil
}

// -------------------------------------------------------- Bootstrap ----------------------------------------------- //

// TODO: mac should be `NETWORK_IFACE/mac`.

func (svc *ipxe) Boostrap() []byte {
	// init boostrap
	if len(svc.bootstrap) == 0 {
		params := ""
		for p, t := range allowedParamsWithType {
			if params != "" {
				params = fmt.Sprintf("%s&", params)
			}

			if t == none {
				params = fmt.Sprintf("%s%s=${%s}", params, p, p)
			} else {
				params = fmt.Sprintf("%s%s=${%s:%s}", params, p, p, t)
			}
		}

		svc.bootstrap = []byte(fmt.Sprintf(ipxeBootstrapFormat, params))
	}

	return svc.bootstrap
}

const ipxeBootstrapFormat = `#!ipxe
chain ipxe?%s
`

//#!ipxe
//chain ipxe?uuid=${uuid}&mac=${mac:hexhyp}&domain=${domain}&hostname=${hostname}&serial=${serial}&arch=${buildarch:uristring}

const (
	none      ipxeParamType = ""
	uriString ipxeParamType = "uristring"
)

type ipxeParamType string

var allowedParamsWithType = map[string]ipxeParamType{
	// types.Mac,
	// types.BusType,
	// types.BusLoc,
	// types.BusID,
	// types.Chip,
	// types.Ssid,
	// types.ActiveScan,
	// types.Key,

	// IPv4 settings

	// types.Ip,
	// types.Netmask,
	// types.Gateway,
	// types.Dns,
	// types.Domain,

	// Boot settings

	// types.Filename,
	// types.NextServer,
	// types.RootPath,
	// types.SanFilename,
	// types.InitiatorIqn,
	// types.KeepSan,
	// types.SkipSanBoot,

	// Host settings

	// types.Hostname,
	types.Uuid: none,
	// types.UserClass,
	// types.Manufacturer,
	// types.Product,
	// types.Serial,
	// types.Asset,

	// Authentication settings

	// types.Username,
	// types.Password,
	// types.ReverseUsername,
	// types.ReversePassword,

	// Cryptography settings

	// types.Crosscert,
	// types.Trust,
	// types.Cert,
	// types.Privkey,

	// Miscellaneous settings

	types.Buildarch: uriString,
	// types.Cpumodel,
	// types.Cpuvendor,
	// types.DhcpServer,
	// types.Keymap,
	// types.Memsize,
	// types.Platform,
	// types.Priority,
	// types.Scriptlet,
	// types.Syslog,
	// types.Syslogs,
	// types.Sysmac,
	// types.Unixtime,
	// types.UseCached,
	// types.Version,
	// types.Vram,
}
