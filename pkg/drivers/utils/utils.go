package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/portworx/kdmp/pkg/drivers"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

// NamespacedName returns a name in form "<namespace>/<name>".
func NamespacedName(namespace, name string) string {
	v := types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}
	return v.String()
}

// ParseJobID parses input string as namespaced name ("<namespace>/<name>").
func ParseJobID(id string) (namespace, name string, err error) {
	v := strings.SplitN(id, string(types.Separator), 2)
	if len(v) != 2 {
		return "", "", fmt.Errorf("invalid job id")
	}
	return v[0], v[1], nil
}

// IsJobCompleted checks if a kubernetes job is completed.
func IsJobCompleted(j *batchv1.Job) bool {
	for _, c := range j.Status.Conditions {
		if c.Type == batchv1.JobComplete && c.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

// IsJobFailed checks if a kubernetes job is failed.
func IsJobFailed(j *batchv1.Job) bool {
	for _, c := range j.Status.Conditions {
		if c.Type == batchv1.JobFailed && c.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

// ToJobStatus returns a job status for provided parameters.
func ToJobStatus(progress float64, errMsg string) *drivers.JobStatus {
	if len(errMsg) > 0 {
		return &drivers.JobStatus{
			State:  drivers.JobStateFailed,
			Reason: errMsg,
		}
	}

	if drivers.IsTransferCompleted(progress) {
		return &drivers.JobStatus{
			State:            drivers.JobStateCompleted,
			ProgressPercents: progress,
		}
	}

	// TODO: require bounds [0, 100)?
	return &drivers.JobStatus{
		State:            drivers.JobStateInProgress,
		ProgressPercents: progress,
	}
}

// ResticExecutorImage returns a docker image that contains resticexecutor binary.
func ResticExecutorImage() string {
	image := "portworx/resticexecutor"
	if customImage := strings.TrimSpace(os.Getenv("KDMP_RESTICEXECUTOR_IMAGE")); customImage != "" {
		image = customImage
	}
	return image
}

// ResticExecutorImageSecret returns an image pull secret for the resticexecutor image.
func ResticExecutorImageSecret() string {
	return strings.TrimSpace(os.Getenv("KDMP_RESTICEXECUTOR_IMAGE_SECRET"))
}

// RsyncImage returns a docker image that contains rsync binary.
func RsyncImage() string {
	image := "eeacms/rsync"
	if customImage := strings.TrimSpace(os.Getenv("KDMP_RSYNC_IMAGE")); customImage != "" {
		image = customImage
	}
	return image
}

// RsyncCommandFlags allows to change rsync command flags.
func RsyncCommandFlags() string {
	return strings.TrimSpace(os.Getenv("KDMP_RSYNC_FLAGS"))
}

// RsyncOpenshiftSCC is used to set a custom openshift security context constraints for a rsync deployment.
func RsyncOpenshiftSCC() string {
	return strings.TrimSpace(os.Getenv("KDMP_RSYNC_OPENSHIFT_SCC"))
}

// RsyncImageSecret returns an image pull secret for the rsync image.
func RsyncImageSecret() string {
	return strings.TrimSpace(os.Getenv("KDMP_RSYNC_IMAGE_SECRET"))
}

// ToImagePullSecret converts a secret name to the ImagePullSecret struct.
func ToImagePullSecret(name string) []corev1.LocalObjectReference {
	if name == "" {
		return nil
	}
	return []corev1.LocalObjectReference{
		{
			Name: name,
		},
	}

}
