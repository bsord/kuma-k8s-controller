package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

var (
	configPath string
	runCmd     = &cobra.Command{
		Use:   "run",
		Short: "Start kuma-k8s-operator",
		Run:   run,
	}

	signalChannel = make(chan os.Signal, 1) // for trapping SIGHUP and friends
)

func init() {

	runCmd.PersistentFlags().StringVarP(&configPath, "config", "c",
		"/configs/dmarcd.conf.json", "Path toe the configuration file")

	rootCmd.AddCommand(runCmd)
}

func sigHandler() {
	// handle SIGHUP for reloading the configuration while running
	signal.Notify(signalChannel,
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGUSR1,
	)
	// Keep the daemon busy by waiting for signals to come
	for sig := range signalChannel {
		if sig == syscall.SIGHUP {
			//do something on sig up
			//d.ReloadConfigFile(configPath)
		} else if sig == syscall.SIGUSR1 {
			//do something on reload?
			//d.ReopenLogs()
		} else if sig == syscall.SIGTERM || sig == syscall.SIGQUIT || sig == syscall.SIGINT {
			//handle shutdown gracefully
			//mainlog.Infof("Shutdown signal caught")
			//d.Shutdown()
			//mainlog.Infof("Shutdown completed, exiting.")
			return
		} else {
			//handle unknown signal
			//mainlog.Infof("Shutdown, unknown signal caught")
			return
		}
	}
}

func run(cmd *cobra.Command, args []string) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for {

		ingressList, err := clientset.NetworkingV1().Ingresses("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		ingresses := ingressList.Items
		if len(ingresses) > 0 {
			for _, ingress := range ingresses {
				fmt.Printf("ingress %s exists in namespace %s on host %s at path %s\n", ingress.Name, ingress.Namespace, ingress.Spec.Rules[0].Host, ingress.Spec.Rules[0].HTTP.Paths[0].String())
			}
		} else {
			fmt.Println("no ingress found")
		}

		time.Sleep(300 * time.Second)
	}

}
