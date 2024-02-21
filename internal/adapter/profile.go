package adapter

import (
	"context"
	"errors"
	"k8s.io/apimachinery/pkg/types"

	"github.com/alexandremahdhaoui/ipxe-api/pkg/v1alpha1"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ProfileModel struct{}

type Profile interface {
	FindBySelectors(ctx context.Context, selectors IpxeSelectors) (ProfileModel, error)

	// NB: CRUD operations are done via the controller-runtime client.Client; only FindBySelectorsAndRender is
	// required.
}

func NewProfile(c client.Client, namespace string) Profile {
	return profile{
		client:    c,
		namespace: namespace,
	}
}

type profile struct {
	client    client.Client
	namespace string
}

func (p profile) FindBySelectors(ctx context.Context, selectors IpxeSelectors) (ProfileModel, error) {
	// list assignment
	list := new(v1alpha1.AssignmentList)
	if err := p.client.List(ctx, list, toListOptions(selectors)...); err != nil {
		return ProfileModel{}, err //TODO: wrap this err.
	}

	if list == nil || len(list.Items) == 0 {
		// Get the list of default matching the buildarch
		if err := p.client.List(ctx, list, toDefaultListOptions(selectors.Buildarch)...); err != nil {
			return ProfileModel{}, err //TODO: wrap this err.
		}

		if list == nil || len(list.Items) == 0 {
			return ProfileModel{}, errors.New("TODO") //TODO: err
		}
	}

	profileName := list.Items[0].Spec.ProfileName
	prof := new(v1alpha1.Profile)

	if err := p.client.Get(ctx, types.NamespacedName{Namespace: p.namespace, Name: profileName}, prof); err != nil {
		return ProfileModel{}, err //TODO: wrap err
	}

	//TODO: convert profile to a ProfileModel && return it.

	return toProfileModel(prof), nil
}

func toListOptions(selectors IpxeSelectors) []client.ListOption {
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

func toProfileModel(input *v1alpha1.Profile) ProfileModel {
	return ProfileModel{}
}
