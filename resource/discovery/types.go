package discovery

const (
	ServiceTypeHTTP ServiceType = iota + 1
	ServiceTypeGRPC
)

type (
	Service struct {
		Name string
		IP   string
		Port int
		Type ServiceType
		// HealthCheckPath is the path to the health check endpoint, grpc keep empty
		HealthCheckPath string
	}
)

type ServiceType int
