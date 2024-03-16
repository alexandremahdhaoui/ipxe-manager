package adapters

import (
	"context"
	"errors"
	"github.com/alexandremahdhaoui/ipxer/internal/types"
	"github.com/alexandremahdhaoui/ipxer/pkg/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	ErrAssignmentNotFound = errors.New("assignment cannot be found")

	errAssignmentFindDefault     = errors.New("finding default assignment")
	errAssignmentFindBySelectors = errors.New("error finding assignment by selectors")
	errAssignmentList            = errors.New("listing assignment")
)

// --------------------------------------------------- INTERFACES --------------------------------------------------- //

type Assignment interface {
	FindDefault(ctx context.Context, buildarch string) (string, error)
	FindBySelectors(ctx context.Context, selectors types.IpxeSelectors) (string, error)
}

// --------------------------------------------------- CONSTRUCTORS ------------------------------------------------- //

func NewAssignment(c client.Client, namespace string) Assignment {
	return &assignment{
		client:    c,
		namespace: namespace,
	}
}

// --------------------------------------------- CONCRETE IMPLEMENTATION -------------------------------------------- //

type assignment struct {
	client    client.Client
	namespace string
}

// --------------------------------------------- FindDefault -------------------------------------------------------- //

func (a *assignment) FindDefault(ctx context.Context, buildarch string) (string, error) {
	// list assignment
	list := new(v1alpha1.AssignmentList)

	// Get the list of default matching the buildarch
	if err := a.client.List(ctx, list, toDefaultListOptions(buildarch)...); err != nil {
		return "", errors.Join(err, errAssignmentList, errAssignmentFindDefault)
	}

	if list == nil || len(list.Items) == 0 {
		return "", errors.Join(ErrAssignmentNotFound, errAssignmentFindDefault)
	}

	return list.Items[0].Spec.ProfileName, nil
}

// --------------------------------------------- FindBySelectors ---------------------------------------------------- //

func (a *assignment) FindBySelectors(ctx context.Context, selectors types.IpxeSelectors) (string, error) {
	// list assignment
	list := new(v1alpha1.AssignmentList)
	if err := a.client.List(ctx, list, toListOptions(selectors)...); err != nil {
		return "", errors.Join(err, errAssignmentList, errAssignmentFindBySelectors)
	}

	if list == nil || len(list.Items) == 0 {
		return "", errors.Join(ErrAssignmentNotFound, errAssignmentFindBySelectors)
	}

	return list.Items[0].Spec.ProfileName, nil
}

// --------------------------------------------- UTILS -------------------------------------------------------------- //

func toListOptions(selectors types.IpxeSelectors) []client.ListOption {
	return []client.ListOption{
		client.MatchingLabels{
			v1alpha1.BuildarchAssignmentLabel: selectors.Buildarch,
		},
		client.HasLabels{
			v1alpha1.UUIDLabelSelector(selectors.UUID),
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
