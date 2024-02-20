package adapter

import (
	"context"
	"errors"

	"github.com/alexandremahdhaoui/ipxe-api/internal/dsa"
	"github.com/alexandremahdhaoui/ipxe-api/pkg/v1alpha1"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Profile interface {
	FindByLabelSelectorsAndRender(ctx context.Context, selectors dsa.IpxeSelectors) ([]byte, error)

	// NB: CRUD operations are done via the controller-runtime client.Client; only FindByLabelSelectorsAndRender is
	// required.
}

func NewProfile(c client.Client) Profile {
	return profile{client: c}
}

type profile struct {
	client client.Client
}

func (p profile) FindByLabelSelectorsAndRender(ctx context.Context, selectors dsa.IpxeSelectors) ([]byte, error) {
	// list assignment
	list := new(v1alpha1.AssignmentList)
	if err := p.client.List(ctx, list, toListOptions(selectors)...); err != nil {
		return nil, err //TODO: wrap this err.
	}

	if list == nil || len(list.Items) == 0 {
		// Get the list of default matching the buildarch
		if err := p.client.List(ctx, list, toDefaultListOptions(selectors.Buildarch)...); err != nil {
			return nil, err //TODO: wrap this err.
		}

		if list == nil || len(list.Items) == 0 {
			return nil, errors.New("TODO") //TODO: err
		}
	}

	//TODO: 1. Select the right Profile or choose the Default assignment;
	//      2. Render the Profile associated with that Assignment. (render recursively if required).

	return nil, nil
}

func toListOptions(selectors dsa.IpxeSelectors) []client.ListOption {
	return []client.ListOption{
		client.MatchingLabels{
			v1alpha1.BuildarchAssignmentLabel: selectors.Buildarch,
		},
		client.HasLabels{
			v1alpha1.UUIDLabelSelector(selectors.Uuid),
		},
	}
}

func toDefaultListOptions(buildarch string) []client.ListOption {
	return []client.ListOption{
		client.MatchingLabels{
			v1alpha1.BuildarchAssignmentLabel: buildarch,
		},
		client.HasLabels{
			v1alpha1.DefaultAssignmentLabel,
		},
	}
}
