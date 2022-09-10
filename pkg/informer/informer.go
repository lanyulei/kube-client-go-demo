package informer

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/informers"
	"kube-client-go-demo/projects/demo4/pkg/client"
	"time"
)

/*
  @Author : lanyulei
  @Desc :
*/

var sharedInformerFactory informers.SharedInformerFactory

func NewSharedInformerFactory(stopCh <-chan struct{}) (err error) {
	var (
		clients client.Clients
	)

	// 1. 加载客户端
	clients = client.NewClients()

	// 2. 实例化 sharedInformerFactory
	sharedInformerFactory = informers.NewSharedInformerFactory(clients.ClientSet(), time.Second*60)

	// 3. 启动 informer
	gvrs := []schema.GroupVersionResource{
		{Group: "", Version: "v1", Resource: "pods"},
		{Group: "", Version: "v1", Resource: "services"},
		{Group: "", Version: "v1", Resource: "namespaces"},

		{Group: "apps", Version: "v1", Resource: "deployments"},
		{Group: "apps", Version: "v1", Resource: "statefulsets"},
		{Group: "apps", Version: "v1", Resource: "daemonsets"},
	}

	for _, v := range gvrs {
		// 创建 informer
		_, err = sharedInformerFactory.ForResource(v)
		if err != nil {
			return
		}
	}

	// 启动所有创建的 informer
	sharedInformerFactory.Start(stopCh)

	// 等待所有 informer 全量同步数据完成
	sharedInformerFactory.WaitForCacheSync(stopCh)

	return
}

func Get() informers.SharedInformerFactory {
	return sharedInformerFactory
}

func Setup(stopCh <-chan struct{}) (err error) {
	err = NewSharedInformerFactory(stopCh)
	if err != nil {
		return
	}
	return
}
