package transform

import (
	"github.com/fusor/cpma/pkg/api"
	"github.com/fusor/cpma/pkg/transform/clusterreport"
)

// ClusterReportExtraction is an API specific transform
type ClusterReportExtraction struct {
	api.APIResources
}

// ClusterReportTransform is an API specific transform
type ClusterReportTransform struct {
}

// Transform transform
func (e ClusterReportExtraction) Transform() ([]Output, error) {
	clusterReport, err := clusterreport.Report(api.APIResources{
		NamespaceList:        e.NamespaceList,
		PersistentVolumeList: e.PersistentVolumeList,
		StorageClassList:     e.StorageClassList,
		NamespaceMap:         e.NamespaceMap,
	})

	if err != nil {
		return nil, err
	}

	output := ClusterOutput{clusterReport}

	outputs := []Output{output}
	return outputs, nil
}

// Validate validate
func (e ClusterReportExtraction) Validate() error {
	return nil
}

// Extract collects data for cluster report
func (e ClusterReportTransform) Extract() (Extraction, error) {
	extraction := &ClusterReportExtraction{}

	namespacesList, err := api.ListNamespaces()
	if err != nil {
		return nil, err
	}
	extraction.NamespaceList = namespacesList

	// Map all namespaces to their resources
	extraction.NamespaceMap = make(map[string]*api.NamespaceResources)
	for _, namespace := range namespacesList.Items {
		namespaceResources := &api.NamespaceResources{}

		podsList, err := api.ListPods(namespace.Name)
		if err != nil {
			return nil, err
		}
		namespaceResources.PodList = podsList

		extraction.NamespaceMap[namespace.Name] = namespaceResources
	}

	pvList, err := api.ListPVs()
	if err != nil {
		return nil, err
	}
	extraction.PersistentVolumeList = pvList

	storageClassList, err := api.ListStorageClasses()
	if err != nil {
		return nil, err
	}
	extraction.StorageClassList = storageClassList

	return *extraction, nil
}

// Name returns a human readable name for the transform
func (e ClusterReportTransform) Name() string {
	return SDNComponentName
}
