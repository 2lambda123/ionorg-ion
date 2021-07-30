package sfu

import (
	"fmt"
	"net/http"
	"os"

	"github.com/cloudwebrtc/nats-discovery/pkg/discovery"
	nrpc "github.com/cloudwebrtc/nats-grpc/pkg/rpc"
	"github.com/cloudwebrtc/nats-grpc/pkg/rpc/reflection"
	log "github.com/pion/ion-log"
	isfu "github.com/pion/ion-sfu/pkg/sfu"
	"github.com/pion/ion/pkg/ion"
	"github.com/pion/ion/pkg/proto"
	pb "github.com/pion/ion/proto/rtc"
	"github.com/spf13/viper"
)

const (
	portRangeLimit = 100
)

type global struct {
	Pprof string `mapstructure:"pprof"`
	Dc    string `mapstructure:"dc"`
}

type natsConf struct {
	URL string `mapstructure:"url"`
}

type nodeConf struct {
	NID string `mapstructure:"nid"`
}

// Config defines parameters for the logger
type logConf struct {
	Level string `mapstructure:"level"`
}

// Config for sfu node
type Config struct {
	Global global   `mapstructure:"global"`
	Log    logConf  `mapstructure:"log"`
	Nats   natsConf `mapstructure:"nats"`
	Node   nodeConf `mapstructure:"node"`
	isfu.Config
}

func unmarshal(rawVal interface{}) error {
	if err := viper.Unmarshal(rawVal); err != nil {
		return err
	}
	return nil
}

func (c *Config) Load(file string) error {
	_, err := os.Stat(file)
	if err != nil {
		return err
	}

	viper.SetConfigFile(file)
	viper.SetConfigType("toml")

	err = viper.ReadInConfig()
	if err != nil {
		log.Errorf("config file %s read failed. %v\n", file, err)
		return err
	}

	err = unmarshal(c)
	if err != nil {
		return err
	}
	err = unmarshal(&c.Config)
	if err != nil {
		return err
	}
	if err != nil {
		log.Errorf("config file %s loaded failed. %v\n", file, err)
		return err
	}

	if len(c.WebRTC.ICEPortRange) > 2 {
		err = fmt.Errorf("config file %s loaded failed. range port must be [min,max]", file)
		log.Errorf("err=%v", err)
		return err
	}

	if len(c.WebRTC.ICEPortRange) != 0 && c.WebRTC.ICEPortRange[1]-c.WebRTC.ICEPortRange[0] < portRangeLimit {
		err = fmt.Errorf("config file %s loaded failed. range port must be [min, max] and max - min >= %d", file, portRangeLimit)
		log.Errorf("err=%v", err)
		return err
	}

	log.Infof("config %s load ok!", file)
	return nil
}

// SFU represents a sfu node
type SFU struct {
	ion.Node
	s *SFUService
}

// NewSFU create a sfu node instance
func NewSFU(nid string) *SFU {
	s := &SFU{
		Node: ion.NewNode(nid),
	}
	return s
}

// Start sfu node
func (s *SFU) Start(conf Config) error {
	var err error

	if conf.Global.Pprof != "" {
		go func() {
			log.Infof("start pprof on %s", conf.Global.Pprof)
			err := http.ListenAndServe(conf.Global.Pprof, nil)
			if err != nil {
				log.Errorf("http.ListenAndServe err=%v", err)
			}
		}()
	}

	err = s.Node.Start(conf.Nats.URL)
	if err != nil {
		s.Close()
		return err
	}

	s.s = NewSFUService(conf.Config)
	//grpc service
	pb.RegisterRTCServer(s.Node.ServiceRegistrar(), s.s)

	// Register reflection service on nats-rpc server.
	reflection.Register(s.Node.ServiceRegistrar().(*nrpc.Server))

	node := discovery.Node{
		DC:      conf.Global.Dc,
		Service: proto.ServiceSFU,
		NID:     s.Node.NID,
		RPC: discovery.RPC{
			Protocol: discovery.NGRPC,
			Addr:     conf.Nats.URL,
			//Params:   map[string]string{"username": "foo", "password": "bar"},
		},
	}

	go func() {
		err := s.Node.KeepAlive(node)
		if err != nil {
			log.Errorf("sfu.Node.KeepAlive(%v) error %v", s.Node.NID, err)
		}
	}()

	//Watch ALL nodes.
	go func() {
		err := s.Node.Watch(proto.ServiceALL)
		if err != nil {
			log.Errorf("Node.Watch(proto.ServiceALL) error %v", err)
		}
	}()

	return nil
}

// Close all
func (s *SFU) Close() {
	s.Node.Close()
}
