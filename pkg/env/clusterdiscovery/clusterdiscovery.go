package clusterdiscovery

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"github.com/mitchellh/go-homedir"
	"k8s.io/client-go/tools/clientcmd"
)

// DiscoverCluster Get kubeconfig using $KUBECONFIG, if not try ~/.kube/config
// Parse kubeconfig and select cluster from available contexts, then get server url from context
// query k8s api for nodes, get node urls from api response and select master node
func DiscoverCluster() string {
	parseKubeConfig()
	return ""
}

func parseKubeConfig() error {
	// Get kubeconfig using $KUBECONFIG, if not try ~/.kube/config
	var kubeConfigPath string
	kubeconfigEnv := os.Getenv("NAME")
	if kubeconfigEnv != "" {
		kubeConfigPath = kubeconfigEnv
	} else {
		home, err := homedir.Dir()
		if err != nil {
			return errors.New("Can't detect home user directory")
		}
		kubeConfigPath = fmt.Sprintf("%s/.kube/config", home)
	}

	kubeConfigFile, err := ioutil.ReadFile(kubeConfigPath)
	if err != nil {
		return err
	}

	// kubeConfig := &clientcmdapi.Config{}

	kubeConfig, err := clientcmd.Load(kubeConfigFile)
	if err != nil {
		panic(err)
	}

	fmt.Println(kubeConfig)

	return nil
}

func getContextClusters() {

}
