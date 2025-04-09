package discovery

import "context"

type noopClient struct {
}

func (n noopClient) DeregisterService(ctx context.Context, serviceID string) error {
	return nil
}

func (n noopClient) RegisterService(ctx context.Context, service *Service) (string, error) {
	return "", nil

}
