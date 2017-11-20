package manifold

import (
	manifold "github.com/manifoldco/go-manifold"
	"github.com/manifoldco/go-manifold/integrations"
)

func newWrapper(team string, cfgs ...manifold.ConfigFunc) (*integrations.Client, error) {
	cl := manifold.New(cfgs...)
	kube, err := integrations.NewClient(cl, func(s string) *string { return &s }(team))
	if err != nil {
		return nil, err
	}

	return kube, nil
}
