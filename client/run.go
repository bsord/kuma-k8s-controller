package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

var (
	configPath string
	runCmd     = &cobra.Command{
		Use:   "run",
		Short: "Start kuma-k8s-controller",
		Run:   run,
	}
)

func init() {

	runCmd.PersistentFlags().StringVarP(&configPath, "config", "c",
		"/configs/kuma-k8s-controller.conf.json", "Path to the configuration file")

	rootCmd.AddCommand(runCmd)
}

// TODO: Define type struct for post events here

func run(cmd *cobra.Command, args []string) {

	// create in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// create kubernetes client using in-cluster config
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Initiate informer factory
	informerFactory := informers.NewSharedInformerFactory(clientset, 10*time.Minute)

	// create ingress informer inside the factory
	ingressInformer := informerFactory.Networking().V1().Ingresses().Informer()

	// Create a channel to stop the shared informer gracefully
	stopper := make(chan struct{})
	defer close(stopper)

	// Kubernetes serves an utility to handle API crashes
	defer runtime.HandleCrash()

	// define function handlers for ingress informer
	ingressInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(new interface{}) {
			fmt.Println("ingress was added")
		},
		UpdateFunc: func(old, new interface{}) {
			fmt.Println("ingress was updated")
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println("ingress was deleted")
		},
	})

	fmt.Println("test")

	// Run the informer
	ingressInformer.Run(stopper)

	// Handle cache sync failure.
	if !cache.WaitForCacheSync(stopper, ingressInformer.HasSynced) {
		runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
		return
	}

}
