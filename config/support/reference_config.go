package support

import (
	"context"
	"github.com/dubbo/dubbo-go/config"
	"github.com/dubbo/dubbo-go/protocol"
)

import (
	log "github.com/AlexStocks/log4go"
)

import (
	"github.com/dubbo/dubbo-go/common/extension"
)

var refprotocol = extension.GetProtocolExtension("registry")

type ReferenceConfig struct {
	context    context.Context
	Interface  string                    `required:"true"  yaml:"interface"  json:"interface,omitempty"`
	Registries []referenceConfigRegistry `required:"true"  yaml:"registries"  json:"registries,omitempty"`
	Cluster    string                    `default:"failover" yaml:"cluster"  json:"cluster,omitempty"`
	Methods    []method                  `yaml:"methods"  json:"methods,omitempty"`
	URLs       []config.URL              `yaml:"-"`
	invoker    protocol.Invoker
}
type referenceConfigRegistry struct {
	string
}

type method struct {
	name    string `yaml:"name"  json:"name,omitempty"`
	retries int    `yaml:"retries"  json:"retries,omitempty"`
}

func NewReferenceConfig(ctx context.Context) *ReferenceConfig {
	return &ReferenceConfig{context: ctx}
}

func (refconfig *ReferenceConfig) CreateProxy() {
	//首先是user specified URL, could be peer-to-peer address, or register center's address.

	//其次是assemble URL from register center's configuration模式
	urls := refconfig.loadRegistries()

	if len(urls) == 1 {
		refconfig.invoker = refprotocol.Refer(urls[0])
	} else {
		//TODO:multi registries
	}
	//TODO:invoker yincheng 's proxy
}

func (refconfig *ReferenceConfig) loadRegistries() []*config.RegistryURL {
	var urls []*config.RegistryURL
	for _, registry := range refconfig.Registries {
		for _, registryConf := range consumerConfig.Registries {
			if registry.string == registryConf.Id {
				url, err := config.NewRegistryURL(refconfig.context, registryConf.Address)
				if err != nil {
					log.Error("The registry id:%s url is invalid ,and will skip the registry", registryConf.Id)
				} else {
					urls = append(urls, url)
				}

			}
		}

	}
	return urls
}
