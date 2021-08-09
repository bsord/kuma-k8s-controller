package main

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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

	// main loop *ideally this would be a watch/listener or configurable via Cobra flag
	for {

		// get List of ingresses from kubeclient
		ingressList, err := clientset.NetworkingV1().Ingresses("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		// iterate ingresses
		ingresses := ingressList.Items
		if len(ingresses) > 0 {
			for _, ingress := range ingresses {

				// TODO: Iterate each of the sub objects such as rules and iterate
				// TODO: define and extract items from ingress results and store to custom type/struct for posting
				fmt.Printf("ingress %s exists in namespace %s on host %s at path %s\n", ingress.Name, ingress.Namespace, ingress.Spec.Rules[0].Host, ingress.Spec.Rules[0].HTTP.Paths[0].String())
			}
		} else {
			fmt.Println("no ingress found")
		}

		// interval
		time.Sleep(60 * time.Second)
	}

}
