package manifold

import (
	manifold "github.com/manifoldco/go-manifold"
	"github.com/manifoldco/kubernetes-credentials/helpers/client"
)

func newWrapper(team string, cfgs ...manifold.ConfigFunc) (*client.Client, error) {
	cl := manifold.New(cfgs...)
	kube, err := client.New(cl, func(s string) *string { return &s }(team))
	if err != nil {
		return nil, err
	}

	return kube, nil
}
