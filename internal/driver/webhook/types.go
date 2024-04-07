package webhook

import (
	"context"
	"errors"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
)

type validatingFunc = func(ctx context.Context, obj runtime.Object) error

func NewUnsupportedResource(obj runtime.Object, errs ...error) error {
	return errors.Join(
		errors.Join(errs...),
		fmt.Errorf("webhook does not support resource with GroupVersionKind=\"%#v\"",
			obj.GetObjectKind().GroupVersionKind()),
	)
}
