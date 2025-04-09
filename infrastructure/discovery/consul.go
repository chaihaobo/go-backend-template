package discovery

import (
	"context"
	"fmt"

	"github.com/hashicorp/consul/api"

	"github.com/chaihaobo/be-template/resource"
)

type consulClient struct {
	res       resource.Resource
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

func NewConsulClient(res resource.Resource) (Client, error) {
	config := api.DefaultConfig()
	config.Address = res.Configuration().Service.DiscoveryServerURL
	if config.Address == "" {
		return noopClient{}, nil
	}
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &consulClient{
		res:       res,
		consulAPI: client,
	}, nil
}
