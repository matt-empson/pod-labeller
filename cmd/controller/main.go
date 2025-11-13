package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/matt-empson/pod-labeller/internal/config"
	"github.com/matt-empson/pod-labeller/internal/controller"
	"github.com/matt-empson/pod-labeller/internal/kube"
	"k8s.io/client-go/informers"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	conf, err := config.NewConfigFromFlags(os.Args[1:])
	if err != nil {
		fmt.Println("conf", err)
	}

	clientBuilder := kube.NewClientBuilder()
	client, err := clientBuilder.NewClient(conf)
	if err != nil {
		fmt.Println("client error", err)
	}

	err = client.CheckConnection(context.Background(), "kube-system")
	if err != nil {
		fmt.Println("connection err", err)
	}

	f := informers.NewSharedInformerFactory(client.ClientSet, 0)

	pods := f.Core().V1().Pods()

	c := controller.NewController(client, pods, conf)

	f.Start(ctx.Done())
	c.Run(ctx)
}
