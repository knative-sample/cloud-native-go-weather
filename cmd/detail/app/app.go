package app

import (
	"strings"

	"fmt"
	"os"

	"github.com/golang/glog"
	"github.com/knative-sample/cloud-native-go-weather/cmd/detail/app/options"
	"github.com/knative-sample/cloud-native-go-weather/pkg/db"
	"github.com/knative-sample/cloud-native-go-weather/pkg/detail"
	"github.com/knative-sample/cloud-native-go-weather/pkg/version"
	"github.com/spf13/cobra"
)

// start api
func NewCommandStartServer(stopCh <-chan struct{}) *cobra.Command {
	ops := &options.Options{}
	mainCmd := &cobra.Command{
		Short: "AppOS",
		Long:  "Application Operating System",
		RunE: func(c *cobra.Command, args []string) error {
			glog.V(2).Infof("NewCommandStartServer main:%s", strings.Join(args, " "))
			run(stopCh, ops)
			return nil
		},
	}

	ops.SetOps(mainCmd)
	return mainCmd
}

func run(stopCh <-chan struct{}, ops *options.Options) {
	vs := version.Version().Info("Application Operating System")
	if ops.Version {
		fmt.Println(vs)
		os.Exit(0)
	}
	if ops.ZipKinEndpoint == "" {
		glog.Fatalf("zipkin --zipkin-endpoint is empty")
	}
	instanceIp := os.Getenv("INSTANCE_IP")
	if instanceIp == "" {
		instanceIp = "127.0.0.1"
	}
	endpoint := os.Getenv("OTS_ENDPOINT")
	tableName := os.Getenv("TABLE_NAME")
	instanceName := os.Getenv("OTS_INSTANCENAME")
	accessKeyId := os.Getenv("OTS_KEYID")
	accessKeySecret := os.Getenv("OTS_SECRET")
	beta := os.Getenv("beta")
	cm := &detail.Server{
		Port:           ops.Port,
		InstanceIp:     instanceIp,
		ServiceName:    ops.ServiceName,
		ZipKinEndpoint: ops.ZipKinEndpoint,
		Beta:           beta,
		TableStoreConfig: &db.TableStoreConfig{
			Endpoint:        endpoint,
			TableName:       tableName,
			InstanceName:    instanceName,
			AccessKeyId:     accessKeyId,
			AccessKeySecret: accessKeySecret,
		},
	}

	go func() {
		if err := cm.Start(); err != nil {
			glog.Fatalf("start detail server error:%s", err.Error())
		}
	}()

	<-stopCh
}
