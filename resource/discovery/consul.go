package discovery

import (
	"context"
	"fmt"

	"github.com/hashicorp/consul/api"

	"github.com/chaihaobo/be-template/resource/config"
)

type consulClient struct {
	consulAPI *api.Client
}

func (c consulClient) DeregisterService(ctx context.Context, serviceID string) error {
	return c.consulAPI.Agent().ServiceDeregister(serviceID)
}

func (c consulClient) RegisterService(ctx context.Context, service *Service) (string, error) {
	var (
		ip          = service.IP
		port        = service.Port
		serviceName = service.Name
	)
	addr := fmt.Sprintf("%s:%d", ip, port)
	serviceID := fmt.Sprintf("%s-%s", serviceName, addr)

	serviceCheck := &api.AgentServiceCheck{
		Interval:                       "10s",
		DeregisterCriticalServiceAfter: "10s",
		Timeout:                        "10s",
	}
	if service.Type == ServiceTypeGRPC {
		serviceCheck.GRPC = addr
	}
	if service.Type == ServiceTypeHTTP {
		serviceCheck.HTTP = fmt.Sprintf("http://%s%s", addr, service.HealthCheckPath)
	}

	return serviceID, c.consulAPI.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Port:    port,
		Address: ip,
		Check:   serviceCheck,
	})
}

func NewConsulClient(config *config.Configuration) (Client, error) {
	apiConfig := api.DefaultConfig()
	apiConfig.Address = config.Service.DiscoveryServerURL
	if apiConfig.Address == "" {
		return noopClient{}, nil
	}
	client, err := api.NewClient(apiConfig)
	if err != nil {
		return nil, err
	}
	return &consulClient{
		consulAPI: client,
	}, nil
}
