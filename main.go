package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var kubeconfig *string
var nodeName string = "unknown"
var nodeStatusReady string = "unknown"
var nodeStatusNetwork string = "unknown"
var nodeStatusMemory string = "unknown"
var nodeStatusDisk string = "unknown"
var nodeStatusPid string = "unknown"
var nodeUnschedulable string = "unknown"

func startWebServer() {
	r := mux.NewRouter()
	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	})
	r.HandleFunc("/node/conditions", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Node: %s\n", nodeName)
		fmt.Fprintf(w, "NetworkUnavailable: %s\n", nodeStatusNetwork)
		fmt.Fprintf(w, "MemoryPressure: %s\n", nodeStatusMemory)
		fmt.Fprintf(w, "DiskPressure: %s\n", nodeStatusDisk)
		fmt.Fprintf(w, "PIDPressure: %s\n", nodeStatusPid)
		fmt.Fprintf(w, "Ready: %s\n", nodeStatusReady)
	})
	r.HandleFunc("/node/conditions/network", func(w http.ResponseWriter, r *http.Request) {
		if nodeStatusNetwork != "Healthy" {
			w.WriteHeader(500)
		}
		fmt.Fprintf(w, "Node: %s\n", nodeName)
		fmt.Fprintf(w, "NetworkUnavailable: %s\n", nodeStatusNetwork)
	})
	r.HandleFunc("/node/conditions/memory", func(w http.ResponseWriter, r *http.Request) {
		if nodeStatusMemory != "Healthy" {
			w.WriteHeader(500)
		}
		fmt.Fprintf(w, "Node: %s\n", nodeName)
		fmt.Fprintf(w, "MemoryPressure: %s\n", nodeStatusMemory)
	})
	r.HandleFunc("/node/conditions/disk", func(w http.ResponseWriter, r *http.Request) {
		if nodeStatusDisk != "Healthy" {
			w.WriteHeader(500)
		}
		fmt.Fprintf(w, "Node: %s\n", nodeName)
		fmt.Fprintf(w, "DiskPressure: %s\n", nodeStatusDisk)
	})
	r.HandleFunc("/node/conditions/pid", func(w http.ResponseWriter, r *http.Request) {
		if nodeStatusPid != "Healthy" {
			w.WriteHeader(500)
		}
		fmt.Fprintf(w, "Node: %s\n", nodeName)
		fmt.Fprintf(w, "PIDPressure: %s\n", nodeStatusPid)
	})
	r.HandleFunc("/node/conditions/ready", func(w http.ResponseWriter, r *http.Request) {
		if nodeStatusReady != "Healthy" {
			w.WriteHeader(500)
		}
		fmt.Fprintf(w, "Node: %s\n", nodeName)
		fmt.Fprintf(w, "Ready: %s\n", nodeStatusReady)
	})
	r.HandleFunc("/node/unschedulable", func(w http.ResponseWriter, r *http.Request) {
		if nodeUnschedulable != "Healthy" {
			w.WriteHeader(500)
		}
		fmt.Fprintf(w, "Node: %s\n", nodeName)
		fmt.Fprintf(w, "Schedulable: %s\n", nodeUnschedulable)
	})
	http.ListenAndServe(":8888", r)
}

func main() {

	go startWebServer()

	var present bool
	nodeName, present = os.LookupEnv("NODE_NAME")
	if !present {
		panic("Missing NODE_NAME variable")
	}
	fmt.Printf("Node: %s\n", nodeName)
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for {
		node, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			fmt.Printf("Node not found\n")
			nodeStatusReady = "NotFound"
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting node %v\n", statusError.ErrStatus.Message)
			nodeStatusReady = "Error"
		} else if err != nil {
			panic(err.Error())
		}
		for _, condition := range node.Status.Conditions {
			fmt.Printf("%s: %s\n", condition.Type, condition.Status)
			if condition.Type == "Ready" {
				if condition.Status == "True" {
					nodeStatusReady = "Healthy"
				} else {
					nodeStatusReady = "Problem"
				}
			} else if condition.Type == "NetworkUnavailable" {
				if condition.Status == "False" {
					nodeStatusNetwork = "Healthy"
				} else {
					nodeStatusNetwork = "Problem"
				}
			} else if condition.Type == "MemoryPressure" {
				if condition.Status == "False" {
					nodeStatusMemory = "Healthy"
				} else {
					nodeStatusMemory = "Problem"
				}
			} else if condition.Type == "DiskPressure" {
				if condition.Status == "False" {
					nodeStatusDisk = "Healthy"
				} else {
					nodeStatusDisk = "Problem"
				}
			} else if condition.Type == "PIDPressure" {
				if condition.Status == "False" {
					nodeStatusPid = "Healthy"
				} else {
					nodeStatusPid = "Problem"
				}
			}
		}
		if node.Spec.Unschedulable {
			nodeUnschedulable = "Problem"
		} else {
			nodeUnschedulable = "Healthy"
		}
		time.Sleep(60 * time.Second)
	}
}
