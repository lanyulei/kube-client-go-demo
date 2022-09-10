package main

import (
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/proxy"
	"k8s.io/client-go/rest"
	"kube-client-go-demo/projects/demo4/pkg/client"
	"kube-client-go-demo/projects/demo4/pkg/informer"
	"net/http"
)

/*
  @Author : lanyulei
  @Desc :
*/

/*
   Gin + client-go 开发示例
*/

func main() {
	stopCh := make(chan struct{})

	// 进行了 client-go 的初始化操作
	err := informer.Setup(stopCh)
	if err != nil {
		panic(err.Error())
	}

	// 启动一个 web 服务
	// 1. 实例化 Gin
	r := gin.Default()

	// 2. 写路由
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    20000, // 20000 的返回值表示正常，其他表示错误
			"message": "pong",
		})
	})

	// 直接查询某一种资源数据的
	r.GET("/pod/list", func(c *gin.Context) {

		items, err := informer.Get().Core().V1().Pods().Lister().List(labels.Everything())
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    40000, // 20000 的返回值表示正常，其他表示错误
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    20000, // 20000 的返回值表示正常，其他表示错误
			"message": "success",
			"data":    items,
		})
	})

	// 动态获取资源数据
	r.GET("/:resource/:group/:version", func(c *gin.Context) {

		resource := c.Param("resource")
		group := c.Param("group")
		version := c.Param("version")

		// 组合成 gvr
		gvr := schema.GroupVersionResource{
			Group:    group,
			Version:  version,
			Resource: resource,
		}

		// 通过 gvr 获取到对应资源对象的 informer
		i, err := informer.Get().ForResource(gvr)
		if err != nil {
			panic(err.Error())
		}

		items, err := i.Lister().List(labels.Everything())
		if err != nil {
			panic(err.Error())
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    20000, // 20000 的返回值表示正常，其他表示错误
			"message": "success",
			"data":    items,
		})
	})

	// 转发请求到 k8s apiserver
	r.Any("/apis/*action", func(c *gin.Context) {

		// transport http.client 的一个封装，也是 client-go 封装的。
		t, err := rest.TransportFor(client.GetConfig())
		if err != nil {
			panic(err)
		}

		s := *c.Request.URL
		s.Host = "192.168.31.141:6443" // k8s 的 apiserver 的地址
		s.Scheme = "https"

		httpProxy := proxy.NewUpgradeAwareHandler(&s, t, true, false, nil)
		httpProxy.UpgradeTransport = proxy.NewUpgradeRequestRoundTripper(t, t)

		c.JSON(http.StatusOK, gin.H{
			"code":    20000, // 20000 的返回值表示正常，其他表示错误
			"message": "success",
			"data":    "",
		})
	})

	// 指定指定方法
	_ = r.Run(":8888")
}
