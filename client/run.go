package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

var (
	configPath string
	runCmd     = &cobra.Command{
		Use:   "run",
		Short: "Start kuma-k8s-operator",
		Run:   run,
	}
)

func init() {

	runCmd.PersistentFlags().StringVarP(&configPath, "config", "c",
		"/configs/kuma-k8s-operator.conf.json", "Path to the configuration file")

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
	ingressInformer := informerFactory.Networking().V1().Ingresses()

	// define function handlers for ingress informer
	ingressInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
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

	// start the informers
	informerFactory.Start(wait.NeverStop)
	informerFactory.WaitForCacheSync(wait.NeverStop)

	// list the ingresses in the informer.
	ingress, err := ingressInformer.Lister().Ingresses("").Get("")
	if err != nil {
		fmt.Println("there was an error listing the ingresses")
	}
	fmt.Println(ingress)

}
