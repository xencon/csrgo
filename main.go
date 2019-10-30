package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/xencon/csrgo/asserts/istio"
	"github.com/xencon/csrgo/asserts/k8s"
	"strings"
	"time"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type ApiV1PodMetricsList struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata   struct {
		SelfLink        string `json:"selfLink"`
		ResourceVersion string `json:"resourceVersion"`
	} `json:"metadata"`
	Items []struct {
		Metadata struct {
			Name              string    `json:"name"`
			Namespace         string    `json:"namespace"`
			SelfLink          string    `json:"selfLink"`
			UID               string    `json:"uid"`
			ResourceVersion   string    `json:"resourceVersion"`
			CreationTimestamp time.Time `json:"creationTimestamp"`
			Labels            struct {
				Run string `json:"run"`
			} `json:"labels"`
		} `json:"metadata"`
		Spec struct {
			Volumes []struct {
				Name   string `json:"name"`
				Secret struct {
					SecretName  string `json:"secretName"`
					DefaultMode string `json:"defaultMode"`
				} `json:"secret"`
			} `json:"volumes"`
			Containers []struct {
				Name      string `json:"name"`
				Image     string `json:"image"`
				Resources struct {
				} `json:"resources"`
				VolumeMounts []struct {
					Name      string `json:"name"`
					ReadOnly  string `json:"readOnly"`
					MountPath string `json:"mountPath"`
				} `json:"volumeMounts"`
				TerminationMessagePath   string `json:"terminationMessagePath"`
				TerminationMessagePolicy string `json:"terminationMessagePolicy"`
				ImagePullPolicy          string `json:"imagePullPolicy"`
				Stdin                    string `json:"stdin"`
				StdinOnce                string `json:"stdinOnce"`
				Tty                      string `json:"tty"`
			} `json:"containers"`
			RestartPolicy                 string `json:"restartPolicy"`
			TerminationGracePeriodSeconds string `json:"terminationGracePeriodSeconds"`
			DnsPolicy                     string `json:"dnsPolicy"`
			ServiceAccountName            string `json:"serviceAccountName"`
			ServiceAccount                string `json:"serviceAccount"`
			NodeName                      string `json:"nodeName"`
			SecurityContext               struct {
			} `json:"securityContext"`
			SchedulerName string `json:"schedulerName"`
			Tolerations   []struct {
				Key               string `json:"key"`
				Operator          string `json:"operator"`
				Effect            string `json:"effect"`
				TolerationSeconds string `json:"tolerationSeconds"`
			} `json:"tolerations"`
			Priority           string `json:"priority"`
			EnableServiceLinks string `json:"enableServiceLinks"`
		} `json:"spec"`
		Status struct {
			Phase      string `json:"phase"`
			Conditions []struct {
				Type               string    `json:"type"`
				Status             string    `json:"status"`
				LastProbeTime      time.Time `json:"lastProbeTime"`
				LastTransitionTime time.Time `json:"lastTransitionTime"`
				Reason             string    `json:"reason"`
			} `json:"conditions"`
			HostIP            string `json:"hostIP"`
			PodIP             string `json:"podIP"`
			StartTime         string `json:"startTime"`
			ContainerStatuses []struct {
				Name  string `json:"name"`
				State struct {
					Terminated struct {
						ExitCode    string    `json:"exitCode"`
						Reason      string    `json:"reason"`
						StartedAt   time.Time `json:"startedAt"`
						FinishedAt  time.Time `json:"finishedAt"`
						ContainerID string    `json:"containerID"`
					} `json:"terminated"`
				} `json:"state"`
				LastState struct {
				} `json:"lastState"`
				Ready        bool   `json:"ready"`
				RestartCount int64  `json:"restartCount"`
				Image        string `json:"image"`
				ImageID      string `json:"imageID"`
				ContainerID  string `json:"containerID"`
			} `json:"containerStatuses"`
			QosClass string `json:"qosClass"`
		} `json:"status"`
	} `json:"items"`
}



func getApiV1Metrics(clientset *kubernetes.Clientset, apiV1Pods *ApiV1PodMetricsList) error {
	data, err := clientset.RESTClient().Get().AbsPath("/api/v1/pods").DoRaw()
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &apiV1Pods)
	return err
}


func main() {

	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println(err)
	}

	clientsetApiV1, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println(err)
	}

	var podsApiV1 ApiV1PodMetricsList

	err = getApiV1Metrics(clientsetApiV1, &podsApiV1)
	fmt.Print("\n\n")

	//Create []string of pod names from api call
	for _, p := range podsApiV1.Items {
		fmt.Println("FOUND POD: ", p.Metadata.Name)
	}
	fmt.Print("\n\n")

	
	//TODO Check all pods exist

	//test All-in-one has correct image version docker.io/jaegertracing/all-in-one:1.9
	for _, p := range podsApiV1.Items {
		if p.Metadata.Namespace == "istio-system" && p.Metadata.Name == "all-in-one" {
			for _, cr := range p.Spec.Containers {
				if strings.Contains(cr.Name, "all-in-one") {
					color.Cyan("POD = " + p.Metadata.Name + "\n")
					fmt.Println("EXP =", env_vars.AllInOneExpected())
					fmt.Println("ACT =", cr.Image)
					fmt.Print("RES = ")
					if cr.Image == env_vars.AllInOneExpected() {
						color.Green("PASS" + "\n\n")
					} else {
						color.Red("FAIL" + "\n\n")
					}
				}
			}
		}

	}

	//test Citadel has correct image version docker.io/istio/citadel:1.2.4
	for _, p := range podsApiV1.Items {
		if p.Metadata.Namespace == "istio-system" && p.Metadata.Name == "citadel" {
			for _, cr := range p.Spec.Containers {
				if strings.Contains(cr.Name, "citadel") {
					color.Cyan("POD = " + p.Metadata.Name + "\n")
					fmt.Println("EXP =", env_vars.CitadelExpected())
					fmt.Println("ACT =", cr.Image)
					fmt.Print("RES = ")
					if cr.Image == env_vars.CitadelExpected() {
						color.Green("PASS" + "\n\n")
					} else {
						color.Red("FAIL" + "\n\n")
					}
				}
			}
		}
	}

	//test Egress-gateway has correct image version docker.io/istio/proxyv2:1.2.4
	for _, p := range podsApiV1.Items {
		if p.Metadata.Namespace == "istio-system" && p.Metadata.Name == "proxyv2" {
			for _, cr := range p.Spec.Containers {
				if strings.Contains(cr.Name, "proxyv2") {
					color.Cyan("POD = " + p.Metadata.Name + "\n")
					fmt.Println("EXP =", env_vars.EgressGatewayExpected())
					fmt.Println("ACT =", cr.Image)
					fmt.Print("RES = ")
					if cr.Image == env_vars.EgressGatewayExpected() {
						color.Green("PASS" + "\n\n")
					} else {
						color.Red("FAIL" + "\n\n")
					}
				}
			}
		}
	}

	//test Galley has correct image version docker.io/istio/galley:1.2.4
	for _, p := range podsApiV1.Items {
		if p.Metadata.Namespace == "istio-system" && p.Metadata.Name == "galley" {
			for _, cr := range p.Spec.Containers {
				if strings.Contains(cr.Name, "galley") {
					color.Cyan("POD = " + p.Metadata.Name + "\n")
					fmt.Println("EXP =", env_vars.GalleyExpected())
					fmt.Println("ACT =", cr.Image)
					fmt.Print("RES = ")
					if cr.Image == env_vars.GalleyExpected() {
						color.Green("PASS" + "\n\n")
					} else {
						color.Red("FAIL" + "\n\n")
					}
				}
			}
		}
	}

	//test Grafana has correct image version grafana/grafana:6.1.6
	for _, p := range podsApiV1.Items {
		if p.Metadata.Namespace == "istio-system" && p.Metadata.Name == "grafana" {
			for _, cr := range p.Spec.Containers {
				if strings.Contains(cr.Name, "grafana") {
					color.Cyan("POD = " + p.Metadata.Name + "\n")
					fmt.Println("EXP =", env_vars.GrafanaExpected())
					fmt.Println("ACT =", cr.Image)
					fmt.Print("RES = ")
					if cr.Image == env_vars.GrafanaExpected() {
						color.Green("PASS" + "\n\n")
					} else {
						color.Red("FAIL" + "\n\n")
					}
				}
			}
		}
	}

	//test Ingress-gateway has correct image version docker.io/istio/proxyv2:1.2.4
	for _, p := range podsApiV1.Items {
		if p.Metadata.Namespace == "istio-system" && p.Metadata.Name == "proxyv2" {
			for _, cr := range p.Spec.Containers {
				if strings.Contains(cr.Name, "proxyv2") {
					color.Cyan("POD = " + p.Metadata.Name + "\n")
					fmt.Println("EXP =", env_vars.IngressGatewayExpected())
					fmt.Println("ACT =", cr.Image)
					fmt.Print("RES = ")
					if cr.Image == env_vars.IngressGatewayExpected() {
						color.Green("PASS" + "\n\n")
					} else {
						color.Red("FAIL" + "\n\n")
					}
				}
			}
		}
	}

	//test Mixer has correct image version docker.io/istio/mixer:1.2.4
	for _, p := range podsApiV1.Items {
		if p.Metadata.Namespace == "istio-system" && p.Metadata.Name == "mixer" {
			for _, cr := range p.Spec.Containers {
				if strings.Contains(cr.Name, "mixer") {
					color.Cyan("POD = " + p.Metadata.Name + "\n")
					fmt.Println("EXP =", env_vars.MixerExpected())
					fmt.Println("ACT =", cr.Image)
					fmt.Print("RES = ")
					if cr.Image == env_vars.MixerExpected() {
						color.Green("PASS" + "\n\n")
					} else {
						color.Red("FAIL" + "\n\n")
					}
				}
			}
		}
	}

	//test Node-agent has correct image version docker.io/istio/node-agent-k8s:1.2.4
	for _, p := range podsApiV1.Items {
		if p.Metadata.Namespace == "istio-system" && p.Metadata.Name == "node-agent" {
			for _, cr := range p.Spec.Containers {
				if strings.Contains(cr.Name, "node-agent") {
					color.Cyan("POD = " + p.Metadata.Name + "\n")
					fmt.Println("EXP =", env_vars.NodeAgentExpected())
					fmt.Println("ACT =", cr.Image)
					fmt.Print("RES = ")
					if cr.Image == env_vars.NodeAgentExpected() {
						color.Green("PASS" + "\n\n")
					} else {
						color.Red("FAIL" + "\n\n")
					}
				}
			}
		}
	}

	//test Pilot has correct image version docker.io/istio/pilot:1.2.4
	for _, p := range podsApiV1.Items {
		if p.Metadata.Namespace == "istio-system" && p.Metadata.Name == "pilot" {
			for _, cr := range p.Spec.Containers {
				if strings.Contains(cr.Name, "pilot") {
					color.Cyan("POD = " + p.Metadata.Name + "\n")
					fmt.Println("EXP =", env_vars.PilotExpected())
					fmt.Println("ACT =", cr.Image)
					fmt.Print("RES = ")
					if cr.Image == env_vars.PilotExpected() {
						color.Green("PASS" + "\n\n")
					} else {
						color.Red("FAIL" + "\n\n")
					}
				}
			}
		}
	}

	//test Prometheus has correct image version docker.io/prom/prometheus:v2.8.0
	for _, p := range podsApiV1.Items {
		if p.Metadata.Namespace == "istio-system" && p.Metadata.Name == "prometheus" {
			for _, cr := range p.Spec.Containers {
				if strings.Contains(cr.Name, "prometheus") {
					color.Cyan("POD = " + p.Metadata.Name + "\n")
					fmt.Println("EXP =", env_vars.PrometheusExpected())
					fmt.Println("ACT =", cr.Image)
					fmt.Print("RES = ")
					if cr.Image == env_vars.PrometheusExpected() {
						color.Green("PASS" + "\n\n")
					} else {
						color.Red("FAIL" + "\n\n")
					}
				}
			}
		}
	}

	//test Quay has correct image version quay.io/kiali/kiali:v0.20
	for _, p := range podsApiV1.Items {
		if p.Metadata.Namespace == "istio-system" && p.Metadata.Name == "quay" {
			for _, cr := range p.Spec.Containers {
				if strings.Contains(cr.Name, "quay") {
					color.Cyan("POD = " + p.Metadata.Name + "\n")
					fmt.Println("EXP =", env_vars.QuayExpected())
					fmt.Println("ACT =", cr.Image)
					fmt.Print("RES = ")
					if cr.Image == env_vars.QuayExpected() {
						color.Green("PASS" + "\n\n")
					} else {
						color.Red("FAIL" + "\n\n")
					}
				}
			}
		}
	}

	//test Sidecar has correct image version docker.io/istio/sidecar_injector:1.2.4
	for _, p := range podsApiV1.Items {
		if p.Metadata.Namespace == "istio-system" && p.Metadata.Name == "sidecar_injector" {
			for _, cr := range p.Spec.Containers {
				if strings.Contains(cr.Name, "sidecar_injector") {
					color.Cyan("POD = " + p.Metadata.Name + "\n")
					fmt.Println("EXP =", env_vars.SidecarExpected())
					fmt.Println("ACT =", cr.Image)
					fmt.Print("RES = ")
					if cr.Image == env_vars.SidecarExpected() {
						color.Green("PASS" + "\n\n")
					} else {
						color.Red("FAIL" + "\n\n")
					}
				}
			}
		}
	}



	//Dummy test kube-addon-manager-minikube has correct image version
	for _, p := range podsApiV1.Items {
		if p.Metadata.Namespace == "kube-system" && p.Metadata.Name == "kube-addon-manager-minikube" {
			for _, cr := range p.Spec.Containers {
				if strings.Contains(cr.Name, "kube-addon-manager") {
					color.Cyan("POD = " + p.Metadata.Name + "\n")
					fmt.Println("EXP =", env_vars.KubeAddonExpected())
					fmt.Println("ACT =", cr.Image)
					fmt.Print("RES = ")
					if cr.Image == env_vars.KubeAddonExpected() {
						color.Green("PASS" + "\n\n")
					} else {
						color.Red("FAIL" + "\n\n")
					}
				}
			}
		}
	}

	//Dummy test kube-api-server-minikube has correct image version
	for _, p := range podsApiV1.Items {
		if p.Metadata.Namespace == "kube-system" && p.Metadata.Name == "kube-apiserver-minikube" {
			for _, cr := range p.Spec.Containers {
				if strings.Contains(cr.Name, "kube-apiserver") {
					color.Cyan("POD = " + p.Metadata.Name + "\n")
					fmt.Println("EXP =", env_vars.KubeApiExpected())
					fmt.Println("ACT =", cr.Image)
					fmt.Print("RES = ")
					//Erroneous evaluation
					if cr.Image == env_vars.KubeApiExpected() {
						color.Green("PASS" + "\n\n")
					} else {
						color.Red("FAIL" + "\n\n")
					}
				}
			}
		}
	}

	//Dummy test kube-addon-manager-minikube has correct image version |||**INJECT ERROR**|||
	for _, p := range podsApiV1.Items {
		if p.Metadata.Namespace == "kube-system" && p.Metadata.Name == "kube-addon-manager-minikube" {
			for _, cr := range p.Spec.Containers {
				if strings.Contains(cr.Name, "kube-addon-manager") {
					color.Cyan("POD = " + p.Metadata.Name + "\n")
					fmt.Println("EXP =", env_vars.KubeAddonExpected())
					fmt.Println("ACT =", "k8s.gcr.io/kube-addon-manager:v9.2")
					fmt.Print("RES = ")
					//Erroneous evaluation
					if "k8s.gcr.io/kube-addon-manager:v9.2" == env_vars.KubeAddonExpected() {
						color.Green("PASS" + "\n\n")
					} else {
						color.Red("FAIL" + "\n\n")
					}
				}
			}
		}
	}
}
