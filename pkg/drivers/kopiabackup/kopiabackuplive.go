package kopiabackup

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/portworx/kdmp/pkg/drivers"
	"github.com/portworx/kdmp/pkg/drivers/utils"
	coreops "github.com/portworx/sched-ops/k8s/core"
	"github.com/sirupsen/logrus"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	defaultPodsMountPath = "/var/lib/kubelet/pods"
)

func jobForLiveBackup(
	jobName,
	namespace,
	pvcName,
	backuplocationName string,
	backuplocationNamespace string,
	mountPod corev1.Pod,
	resources corev1.ResourceRequirements,
	labels map[string]string) (*batchv1.Job, error) {
	volDir, err := getVolumeDirectory(pvcName, namespace)
	if err != nil {
		logrus.Infof("line 32 jobForLiveBackup")
		return nil, err
	}
	logrus.Infof("line 33 jobForLiveBackup")
	// pod volumes reside under /var/lib/kubelet/pods/<podUID>/volumes/<volumePlugin>/<volumeName> directory.
	// mount /var/lib/kubelet/pods/<podUID>/volumes as a /data directory to a resticexecutor job and
	// use /data/*/<volumeName> as a backup directory and determine volume plugin by resticexecutor.
	podVolumesPath := fmt.Sprintf("%s/%s/volumes", defaultPodsMountPath, mountPod.UID)
	backupPath := fmt.Sprintf("/data/*/%s", volDir)

	backupName := jobName

	labels = addJobLabels(labels)

	cmd := strings.Join([]string{
		"/kopiaexecutor",
		"backup",
		"--volume-backup-name",
		backupName,
		"--namespace",
		backuplocationNamespace,
		"--repository",
		toRepoName(pvcName, namespace),
		"--credentials",
		backuplocationName,
		"--secret-file-path",
		filepath.Join(drivers.KopiaSecretMount, drivers.KopiaSecretKey),
		"--source-path-glob",
		backupPath,
	}, " ")

	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					RestartPolicy:      corev1.RestartPolicyOnFailure,
					ImagePullSecrets:   utils.ToImagePullSecret(utils.KopiaExecutorImageSecret()),
					ServiceAccountName: jobName,
					NodeName:           mountPod.Spec.NodeName,
					Containers: []corev1.Container{
						{
							Name:            "kopiaexecutor",
							Image:           utils.KopiaExecutorImage(),
							ImagePullPolicy: corev1.PullIfNotPresent,
							Command: []string{
								"/bin/sh",
								"-x",
								"-c",
								cmd,
							},
							Resources: resources,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "vol",
									MountPath: "/data",
								},
								{
									Name:      "cred-secret",
									MountPath: drivers.KopiaCredSecretMount,
									ReadOnly:  true,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "vol",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: podVolumesPath,
								},
							},
						},
						{
							Name: "cred-secret",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: backuplocationName,
								},
							},
						},
					},
				},
			},
		},
	}, nil
}

// getVolumeDirectory gets the name of the directory on the host, under /var/lib/kubelet/pods/<podUID>/volumes/,
// where the specified volume lives. For volumes with a CSIVolumeSource, append "/mount" to the directory name.
func getVolumeDirectory(pvcName, pvcNamespace string) (string, error) {
	pvc, err := coreops.Instance().GetPersistentVolumeClaim(pvcName, pvcNamespace)
	if err != nil {
		return "", err
	}

	pv, err := coreops.Instance().GetPersistentVolume(pvc.Spec.VolumeName)
	if err != nil {
		return "", err
	}

	// PV's been created with a CSI source.
	if pv.Spec.CSI != nil {
		return pvc.Spec.VolumeName + "/mount", nil
	}

	return pvc.Spec.VolumeName, nil
}
