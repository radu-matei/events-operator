package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang/glog"
	clientset "github.com/radu-matei/events-operator/pkg/client/clientset/versioned"
	informers "github.com/radu-matei/events-operator/pkg/client/informers/externalversions"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfig = getEnvVarOrExit("KUBECONFIG")
)

func main() {

	c := make(chan os.Signal, 2)
	stop := make(chan struct{})

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1)
	}()

	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	exampleClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building example clientset: %s", err.Error())
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)
	exampleInformerFactory := informers.NewSharedInformerFactory(exampleClient, time.Second*30)

	controller := NewController(kubeClient, exampleClient, kubeInformerFactory, exampleInformerFactory)

	go kubeInformerFactory.Start(stop)
	go exampleInformerFactory.Start(stop)

	if err = controller.Run(2, stop); err != nil {
		glog.Fatalf("Error running controller: %s", err.Error())
	}
}

func getEnvVarOrExit(varName string) string {
	value := os.Getenv(varName)
	if value == "" {
		glog.Fatalf("missing environment variable %s\n", varName)
	}

	return value
}
