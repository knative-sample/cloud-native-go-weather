package options

import (
	"github.com/spf13/cobra"
)

type Options struct {
	ResourceRoot  string
	CityService   string
	DetailService string
	Version       bool
	Port          string

	ZipKinEndpoint string
	ServiceName    string
}

func (s *Options) SetOps(ac *cobra.Command) {
	ac.Flags().StringVar(&s.ResourceRoot, "resource", "/var/html", "html resource root")
	ac.Flags().StringVar(&s.CityService, "city-service", "city:8080", "city service, default is city:8080")
	ac.Flags().StringVar(&s.DetailService, "detail-service", "detail:8080", "detail service, default is detail:8080")
	ac.Flags().StringVar(&s.Port, "port", "8080", "http listen port")
	ac.Flags().BoolVar(&s.Version, "version", false, "Print version information")
	ac.Flags().StringVar(&s.ZipKinEndpoint, "zipkin-endpoint", "", "zipkin endpoint")
	ac.Flags().StringVar(&s.ServiceName, "service-name", "weather", "service name default is weather")
}
