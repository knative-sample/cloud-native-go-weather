package options

import (
	"github.com/spf13/cobra"
)

type Options struct {
	ZipKinEndpoint string
	ServiceName    string
	Port           string
	Version        bool
}

func (s *Options) SetOps(ac *cobra.Command) {
	ac.Flags().StringVar(&s.ZipKinEndpoint, "zipkin-endpoint", "", "zipkin endpoint")
	ac.Flags().StringVar(&s.ServiceName, "service-name", "detail", "service name default is detail")
	ac.Flags().StringVar(&s.Port, "port", "8080", "http listen port")
	ac.Flags().BoolVar(&s.Version, "version", false, "Print version information")
}
