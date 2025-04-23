package discovery

import "context"

type Client interface {
	RegisterService(ctx context.Context, service *Service) (string, error)
	DeregisterService(ctx context.Context, serviceID string) error
}
