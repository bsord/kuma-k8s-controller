package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	networkv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

// TODO: Define type struct for post events here

func main() {

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
		AddFunc: handleNewIngress,
		UpdateFunc: func(old, new interface{}) {
			fmt.Println("ingress was updated")
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println("ingress was deleted")
		},
	})

	// Run the informer
	ingressInformer.Run(stopper)

	// Handle cache sync failure.
	if !cache.WaitForCacheSync(stopper, ingressInformer.HasSynced) {
		runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
		return
	}

}

func handleNewIngress(obj interface{}) {
	fmt.Println("Ingress created")
	// Cast the obj as an ingress
	ingress := obj.(*networkv1.Ingress)

	// Get name
	name := ingress.Name
	fmt.Printf("Name: %s\n", name)

	// Get resource version
	resourceVersion := ingress.GetResourceVersion()
	fmt.Printf("Resource Version: %s\n", resourceVersion)

	ingressMon := &ingressMonitor{
		Name:            name,
		ResourceVersion: resourceVersion,
		Annotations:     []string{},
		Paths:           []string{},
	}

	// Get Annotations
	annotations := ingress.GetAnnotations()
	fmt.Println("Annotations:")

	// Iterate annotations
	for _, annotation := range annotations {

		// add annotation to monitor
		ingressMon.Annotations = append(ingressMon.Annotations, annotation)

		// write console
		fmt.Printf("-%s\n", annotation)
	}

	// get rules
	rules := ingress.Spec.Rules

	// iterate each rule
	for _, rule := range rules {

		// get host
		host := rule.Host

		// get paths
		paths := rule.HTTP.Paths
		fmt.Println("Paths:")

		// iterate each path
		for _, path := range paths {

			// add path to monitor
			ingressMon.Paths = append(ingressMon.Paths, "https://"+host+path.Path)

			// write console
			fmt.Printf("-https://%s%s\n", host, path.Path)
		}

	}

	// Define structs for slack message
	type slackField struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}

	type slackBlock struct {
		Type   string       `json:"type"`
		Text   *slackField  `json:"text,omitempty"`
		Fields []slackField `json:"fields,omitempty"`
	}

	type slackMessage struct {
		Blocks []slackBlock `json:"blocks"`
	}

	// Create slack message
	slackMsg := &slackMessage{
		Blocks: []slackBlock{
			{
				Type: "section",
				Text: &slackField{
					Type: "mrkdwn",
					Text: "New ingress has been created:\n*" + name + "*",
				},
			},
			{
				Type: "section",
				Fields: []slackField{
					{
						Type: "mrkdwn",
						Text: "*Name:*\n" + name,
					},
					{
						Type: "mrkdwn",
						Text: "*Resource Version:*\n" + resourceVersion,
					},
					{
						Type: "mrkdwn",
						Text: "*Paths:*\n" + string(strings.Join(ingressMon.Paths[:], ",")),
					},
					{
						Type: "mrkdwn",
						Text: "*Annotations:*\n" + string(strings.Join(ingressMon.Annotations[:], ",")),
					},
				},
			},
		},
	}

	slackMessageJson, _ := json.Marshal(slackMsg)

	slackUrl := os.Getenv("SLACK_URL")

	sendSlackNotif(slackMessageJson, slackUrl)

	fmt.Println(string(slackMessageJson))

}

func sendSlackNotif(jsonMessage json.RawMessage, url string) {

	//var jsonStr = []byte(`{"text":"Buy cheese and bread for breakfast."}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonMessage))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

type ingressMonitor struct {
	Name            string   `json:"name"`
	ResourceVersion string   `json:"resourceVersion"`
	Annotations     []string `json:"annotations"`
	Paths           []string `json:"paths"`
}
