package adapter

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
	FindDefaultProfile(ctx context.Context, buildarch string) (types.Assignment, error)
	FindProfileBySelectors(ctx context.Context, selectors types.IpxeSelectors) (types.Assignment, error)
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

// --------------------------------------------- FindDefaultProfile ------------------------------------------------- //

func (a *assignment) FindDefaultProfile(ctx context.Context, buildarch string) (types.Assignment, error) {
	// list assignment
	list := new(v1alpha1.AssignmentList)

	// Get the list of default matching the buildarch
	if err := a.client.List(ctx, list, toDefaultListOptions(buildarch)...); err != nil {
		return types.Assignment{}, errors.Join(err, errAssignmentList, errAssignmentFindDefault)
	}

	if list == nil || len(list.Items) == 0 {
		return types.Assignment{}, errors.Join(ErrAssignmentNotFound, errAssignmentFindDefault)
	}

	return types.Assignment{
		Name:        list.Items[0].Name,
		ProfileName: list.Items[0].Spec.ProfileName,
	}, nil
}

// --------------------------------------------- FindProfileBySelectors --------------------------------------------- //

func (a *assignment) FindProfileBySelectors(ctx context.Context, selectors types.IpxeSelectors) (types.Assignment, error) {
	// list assignment
	list := new(v1alpha1.AssignmentList)
	if err := a.client.List(ctx, list, toListOptions(selectors)...); err != nil {
		return types.Assignment{}, errors.Join(err, errAssignmentList, errAssignmentFindBySelectors)
	}

	if list == nil || len(list.Items) == 0 {
		return types.Assignment{}, errors.Join(ErrAssignmentNotFound, errAssignmentFindBySelectors)
	}

	return types.Assignment{
		Name:        list.Items[0].Name,
		ProfileName: list.Items[0].Spec.ProfileName,
	}, nil
}

// --------------------------------------------- UTILS -------------------------------------------------------------- //

func toListOptions(selectors types.IpxeSelectors) []client.ListOption {
	return setBuildarchLabelSelector(selectors.Buildarch, []client.ListOption{
		client.HasLabels{v1alpha1.NewUUIDLabelSelector(selectors.UUID)},
	})
}

func toDefaultListOptions(buildarch string) []client.ListOption {
	return setBuildarchLabelSelector(buildarch, []client.ListOption{
		client.HasLabels{v1alpha1.DefaultAssignmentLabel},
	})
}

func setBuildarchLabelSelector(buildarch string, opts []client.ListOption) []client.ListOption {
	switch v1alpha1.Buildarch(buildarch) {
	case v1alpha1.Arm32:
		return append(opts, client.HasLabels{v1alpha1.Arm32BuildarchLabelSelector})
	case v1alpha1.Arm64:
		return append(opts, client.HasLabels{v1alpha1.Arm64BuildarchLabelSelector})
	case v1alpha1.I386:
		return append(opts, client.HasLabels{v1alpha1.I386BuildarchLabelSelector})
	case v1alpha1.X8664:
		return append(opts, client.HasLabels{v1alpha1.X8664BuildarchLabelSelector})
	default:
		// not specifying anything implies any buildarch
		return opts
	}
}
