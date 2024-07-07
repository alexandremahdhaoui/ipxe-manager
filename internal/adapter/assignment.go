package adapter

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/alexandremahdhaoui/ipxer/internal/types"
	"github.com/alexandremahdhaoui/ipxer/pkg/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	ErrAssignmentNotFound = errors.New("assignment not found")

	errAssignmentFindDefault     = errors.New("finding default assignment")
	errAssignmentFindBySelectors = errors.New("error finding assignment by selectors")
	errAssignmentList            = errors.New("listing assignment")
)

// --------------------------------------------------- INTERFACES --------------------------------------------------- //

type Assignment interface {
	FindDefaultByBuildarch(ctx context.Context, buildarch string) (types.Assignment, error)
	FindBySelectors(ctx context.Context, selectors types.IpxeSelectors) (types.Assignment, error)
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

// --------------------------------------------- FindDefaultByBuildarch ------------------------------------------------- //

func (a *assignment) FindDefaultByBuildarch(ctx context.Context, buildarch string) (types.Assignment, error) {
	// list assignment
	list := new(v1alpha1.AssignmentList)

	// Get the list of default matching the buildarch
	if err := a.client.List(ctx, list,
		buildarchLabelSelector(buildarch),
		defaultAssignmentLabelSelector(),
	); err != nil {
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

// --------------------------------------------- FindBySelectors --------------------------------------------- //

func (a *assignment) FindBySelectors(ctx context.Context, selectors types.IpxeSelectors) (types.Assignment, error) {
	// list assignment
	list := new(v1alpha1.AssignmentList)
	if err := a.client.List(ctx, list,
		buildarchLabelSelector(selectors.Buildarch),
		uuidLabelSelector(selectors.UUID),
	); err != nil {
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

func buildarchLabelSelector(buildarch string) client.ListOption {
	switch v1alpha1.Buildarch(buildarch) {
	case v1alpha1.Arm32:
		return client.HasLabels{v1alpha1.Arm32BuildarchLabelSelector}
	case v1alpha1.Arm64:
		return client.HasLabels{v1alpha1.Arm64BuildarchLabelSelector}
	case v1alpha1.I386:
		return client.HasLabels{v1alpha1.I386BuildarchLabelSelector}
	case v1alpha1.X8664:
		return client.HasLabels{v1alpha1.X8664BuildarchLabelSelector}
	default:
		// not specifying anything implies any buildarch
		return nil
	}
}

func defaultAssignmentLabelSelector() client.ListOption {
	return client.HasLabels{v1alpha1.DefaultAssignmentLabel}
}

func uuidLabelSelector(id uuid.UUID) client.ListOption {
	return client.HasLabels{v1alpha1.NewUUIDLabelSelector(id)}
}
