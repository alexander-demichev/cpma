package transform

import (
	"github.com/fusor/cpma/pkg/api"
	"github.com/fusor/cpma/pkg/transform/clusterreport"
	k8sapicore "k8s.io/api/core/v1"
	k8sapistorage "k8s.io/api/storage/v1"
)

// ClusterReportExtraction is an API specific transform
type ClusterReportExtraction struct {
	NamespaceList        *k8sapicore.NamespaceList
	PersistentVolumeList *k8sapicore.PersistentVolumeList
	StorageClassList     *k8sapistorage.StorageClassList
}

// ClusterReportTransform is an API specific transform
type ClusterReportTransform struct {
}

// Transform transform
func (e ClusterReportExtraction) Transform() ([]Output, error) {
	output, err := clusterreport.Report(api.APIResources{
		NamespaceList:        e.NamespaceList,
		PersistentVolumeList: e.PersistentVolumeList,
		StorageClassList:     e.StorageClassList,
	})

	if err != nil {
		return nil, err
	}

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
