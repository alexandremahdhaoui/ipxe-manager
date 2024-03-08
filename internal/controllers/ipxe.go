package controllers

import (
	"bytes"
	"context"
	"fmt"
	"github.com/alexandremahdhaoui/ipxe-api/internal/adapter"
	"github.com/alexandremahdhaoui/ipxe-api/internal/types"
	"text/template"
)

// ---------------------------------------------------- INTERFACES -------------------------------------------------- //

type IPXE interface {
	FindProfileAndRender(ctx context.Context, selectors types.IpxeSelectors) ([]byte, error)
	Boostrap() []byte
}

// --------------------------------------------------- CONSTRUCTORS ------------------------------------------------- //

func NewIPXE(profile adapter.Profile, mux ResolveTransformerMux) IPXE {
	return &ipxe{
		profile: profile,
		mux:     mux,
	}
}

// -------------------------------------------------------- IPXE ---------------------------------------------------- //

type ipxe struct {
	profile adapter.Profile
	mux     ResolveTransformerMux

	bootstrap []byte
}

// -------------------------------------------------------- FindProfileAndRender ------------------------------------ //

func (svc *ipxe) FindProfileAndRender(ctx context.Context, selectors types.IpxeSelectors) ([]byte, error) {
	p, err := svc.profile.FindBySelectors(ctx, selectors)
	if err != nil {
		return nil, err //TODO: wrap
	}

	data, err := svc.mux.ResolveAndTransformBatch(ctx, p.AdditionalContent)
	if err != nil {
		return nil, err //TODO: wrap
	}

	output, err := templateIPXEProfile(p.IPXETemplate, data)
	if err != nil {
		return nil, err //TODO: wrap
	}

	return output, nil
}

func templateIPXEProfile(ipxeTemplate string, data map[string][]byte) ([]byte, error) {
	tpl, err := template.New("").Parse(ipxeTemplate)
	if err != nil {
		return nil, err //TODO: wrap
	}

	buf := bytes.NewBuffer(make([]byte, 0))
	if err := tpl.Execute(buf, data); err != nil {
		return nil, err //TODO: wrap
	}

	return buf.Bytes(), nil
}

// -------------------------------------------------------- Bootstrap ----------------------------------------------- //

func (svc *ipxe) Boostrap() []byte {
	// init boostrap
	if len(svc.bootstrap) == 0 {
		params := ""
		for _, param := range allowedParams {
			if params != "" {
				params = fmt.Sprintf("%s&", params)
			}

			params = fmt.Sprintf("%s%s=${%s}", params, param, param)
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

var (
	allowedParams = []string{
		types.Mac,
		types.BusType,
		types.BusLoc,
		types.BusID,
		types.Chip,
		types.Ssid,
		types.ActiveScan,
		types.Key,
		// IPv4 settings

		types.Ip,
		types.Netmask,
		types.Gateway,
		types.Dns,
		types.Domain,

		//Boot settings

		types.Filename,
		types.NextServer,
		types.RootPath,
		types.SanFilename,
		types.InitiatorIqn,
		types.KeepSan,
		types.SkipSanBoot,

		// Host settings

		types.Hostname,
		types.Uuid,
		types.UserClass,
		types.Manufacturer,
		types.Product,
		types.Serial,
		types.Asset,

		//Authentication settings

		types.Username,
		types.Password,
		types.ReverseUsername,
		types.ReversePassword,

		//Cryptography settings

		types.Crosscert,
		types.Trust,
		types.Cert,
		types.Privkey,

		//Miscellaneous settings

		types.Buildarch,
		types.Cpumodel,
		types.Cpuvendor,
		types.DhcpServer,
		types.Keymap,
		types.Memsize,
		types.Platform,
		types.Priority,
		types.Scriptlet,
		types.Syslog,
		types.Syslogs,
		types.Sysmac,
		types.Unixtime,
		types.UseCached,
		types.Version,
		types.Vram,
	}
)
