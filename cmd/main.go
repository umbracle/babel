package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	ethereumcl "github.com/umbracle/babel/plugin/ethereum_cl"
	ethereumel "github.com/umbracle/babel/plugin/ethereum_el"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mitchellh/mapstructure"
	babelSDK "github.com/umbracle/babel/sdk"
	"google.golang.org/grpc"
)

func main() {
	var pluginName string
	var port uint64

	flag.StringVar(&pluginName, "plugin", "", "name of the plugin to run")
	flag.Uint64Var(&port, "port", 2020, "port to listen the grpc server")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		log.Fatalf("no command found")
		os.Exit(1)
	}

	cmd := args[0]

	config := map[string]interface{}{}
	for _, raw := range flag.Args() {
		parts := strings.SplitN(raw, "=", 2)
		if len(parts) == 2 {
			config[parts[0]] = parts[1]
		}
	}

	plugin, ok := pluginList[pluginName]
	if !ok {
		log.Fatalf("plugin not found: %s", pluginName)
		os.Exit(1)
	}

	// decode the input on the plugin
	if err := mapstructure.Decode(config, plugin); err != nil {
		log.Fatalf("failed to decode input: %v", err)
		os.Exit(1)
	}

	if cmd == "call" {
		resp, err := plugin.Query()
		if err != nil {
			log.Printf(err.Error())
		} else {
			log.Printf(resp.String())
		}
	} else if cmd == "server" {
		lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
			os.Exit(1)
		}

		log.Printf("Grpc server running in 0.0.0.0:%d\n", port)

		srv := &service{
			plugin: plugin,
		}
		var opts []grpc.ServerOption
		grpcServer := grpc.NewServer(opts...)
		babelSDK.RegisterBabelServiceServer(grpcServer, srv)

		grpcServer.Serve(lis)
	}
}

type service struct {
	babelSDK.UnimplementedBabelServiceServer

	plugin Plugin
}

func (s *service) GetSyncStatus(ctx context.Context, req *empty.Empty) (*babelSDK.SyncStatus, error) {
	return s.plugin.Query()
}

type Plugin interface {
	Query() (*babelSDK.SyncStatus, error)
}

var pluginList = map[string]Plugin{
	"ethereum_el": &ethereumel.EthereumEL{},
	"ethereum_cl": &ethereumcl.EthereumCL{},
}
