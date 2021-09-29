/*

LICENSE

*/

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/portworx/kdmp/pkg/apis/kdmp/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// VolumeBackupDeleteLister helps list VolumeBackupDeletes.
// All objects returned here must be treated as read-only.
type VolumeBackupDeleteLister interface {
	// List lists all VolumeBackupDeletes in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.VolumeBackupDelete, err error)
	// VolumeBackupDeletes returns an object that can list and get VolumeBackupDeletes.
	VolumeBackupDeletes(namespace string) VolumeBackupDeleteNamespaceLister
	VolumeBackupDeleteListerExpansion
}

// volumeBackupDeleteLister implements the VolumeBackupDeleteLister interface.
type volumeBackupDeleteLister struct {
	indexer cache.Indexer
}

// NewVolumeBackupDeleteLister returns a new VolumeBackupDeleteLister.
func NewVolumeBackupDeleteLister(indexer cache.Indexer) VolumeBackupDeleteLister {
	return &volumeBackupDeleteLister{indexer: indexer}
}

// List lists all VolumeBackupDeletes in the indexer.
func (s *volumeBackupDeleteLister) List(selector labels.Selector) (ret []*v1alpha1.VolumeBackupDelete, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.VolumeBackupDelete))
	})
	return ret, err
}

// VolumeBackupDeletes returns an object that can list and get VolumeBackupDeletes.
func (s *volumeBackupDeleteLister) VolumeBackupDeletes(namespace string) VolumeBackupDeleteNamespaceLister {
	return volumeBackupDeleteNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// VolumeBackupDeleteNamespaceLister helps list and get VolumeBackupDeletes.
// All objects returned here must be treated as read-only.
type VolumeBackupDeleteNamespaceLister interface {
	// List lists all VolumeBackupDeletes in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.VolumeBackupDelete, err error)
	// Get retrieves the VolumeBackupDelete from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.VolumeBackupDelete, error)
	VolumeBackupDeleteNamespaceListerExpansion
}

// volumeBackupDeleteNamespaceLister implements the VolumeBackupDeleteNamespaceLister
// interface.
type volumeBackupDeleteNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all VolumeBackupDeletes in the indexer for a given namespace.
func (s volumeBackupDeleteNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.VolumeBackupDelete, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.VolumeBackupDelete))
	})
	return ret, err
}

// Get retrieves the VolumeBackupDelete from the indexer for a given namespace and name.
func (s volumeBackupDeleteNamespaceLister) Get(name string) (*v1alpha1.VolumeBackupDelete, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("volumebackupdelete"), name)
	}
	return obj.(*v1alpha1.VolumeBackupDelete), nil
}
