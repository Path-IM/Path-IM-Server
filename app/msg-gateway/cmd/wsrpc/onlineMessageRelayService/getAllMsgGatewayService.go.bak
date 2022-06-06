package onlinemessagerelayservice

import (
	"context"
	"fmt"
	"github.com/Path-IM/Path-IM-Server/common/xetcd"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/core/threading"
	"github.com/zeromicro/go-zero/zrpc"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"strconv"
	"strings"
	"sync"
	"time"
)

// GetAllByEtcd 获取所有 service
// @param msggatewayRpcEtcdKey msggateway-rpc 的 etcd key
func GetAllByEtcd(
	ctx context.Context,
	etcdConf discov.EtcdConf,
	msggatewayRpcEtcdKey string,
) (services []OnlineMessageRelayService, err error) {
	var zrpcConns []zrpc.Client
	zrpcConns, err = xetcd.GetGoZeroZrpcConns(ctx, xetcd.NewClient(etcdConf), msggatewayRpcEtcdKey)
	if err != nil {
		return
	}
	for _, conn := range zrpcConns {
		service := NewOnlineMessageRelayService(conn)
		services = append(services, service)
	}
	return
}

var (
	kube *kubernetes.Clientset
	ks   *k8sSvc
)

type k8sSvc struct {
	msggatewayRpcName string
	namespace         string
	endpoints         []string
	lock              sync.Mutex
	port              int
}

func (s *k8sSvc) OnAdd(obj interface{}) {
	s.Update()
}

func (s *k8sSvc) OnUpdate(oldObj, newObj interface{}) {
	s.Update()
}

func (s *k8sSvc) OnDelete(obj interface{}) {
	s.Update()
}

func (s *k8sSvc) Update() {
	s.lock.Lock()
	defer s.lock.Unlock()
	endpoints, e := kube.CoreV1().Endpoints(ks.namespace).Get(context.Background(), ks.msggatewayRpcName, v1.GetOptions{})
	if e != nil {
		logx.Error("get endpoints error: %s", e)
		return
	}
	s.endpoints = nil
	for _, sub := range endpoints.Subsets {
		for _, point := range sub.Addresses {
			s.endpoints = append(s.endpoints, fmt.Sprintf("%s:%d", point.IP, s.port))
		}
	}
}

func parseTarget(target string) (*k8sSvc, error) {
	err := fmt.Errorf("invalid target: %s", target)
	if !strings.HasPrefix(target, "k8s://") {
		return nil, err
	}
	target = strings.TrimPrefix(target, "k8s://")
	parts := strings.Split(target, "/")
	if len(parts) != 2 {
		return nil, err
	}
	namespace, name := parts[0], parts[1]
	var port = 80
	if strings.Contains(name, ":") {
		parts := strings.Split(name, ":")
		name = parts[0]
		port, err = strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}
	}
	return &k8sSvc{
		namespace:         namespace,
		msggatewayRpcName: name,
		port:              port,
		endpoints:         nil,
	}, nil
}

// GetAllByK8s k8s里获取所有 service
// @param msggatewayRpcEtcdKey msggateway-rpc 的 etcd key
func GetAllByK8s(
	msggatewayRpcTarget string,
) (services []OnlineMessageRelayService, err error) {
	if kube == nil {
		err = initK8sSvc(msggatewayRpcTarget)
		if err != nil {
			panic(err)
		}
	}
	for _, endpoint := range ks.endpoints {
		client, e := zrpc.NewClient(zrpc.RpcClientConf{
			Endpoints: []string{endpoint},
			NonBlock:  true,
			Timeout:   5000,
		})
		if e != nil {
			logx.Error("create zrpc client error: %s", e)
		} else {
			services = append(services, NewOnlineMessageRelayService(client))
		}
	}
	return
}

func initK8sSvc(msggatewayRpcTarget string) (err error) {
	var config *rest.Config
	config, err = rest.InClusterConfig()
	if err != nil {
		return err
	}
	var kb *kubernetes.Clientset
	kb, err = kubernetes.NewForConfig(config)
	if err != nil {
		return err
	} else {
		kube = kb
		var svc *k8sSvc
		svc, err = parseTarget(msggatewayRpcTarget)
		if err != nil {
			return
		}
		ks = svc
		inf := informers.NewSharedInformerFactoryWithOptions(kb, 3*time.Minute,
			informers.WithNamespace(ks.namespace),
			informers.WithTweakListOptions(func(options *v1.ListOptions) {
				options.FieldSelector = "metadata.name=" + svc.msggatewayRpcName
			}))
		in := inf.Core().V1().Endpoints()
		in.Informer().AddEventHandler(ks)
		threading.GoSafe(func() {
			inf.Start(proc.Done())
		})
		ks.Update()
	}
	return nil
}
