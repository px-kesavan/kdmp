package kopia

import (
	"fmt"
	"time"

	kdmpapi "github.com/portworx/kdmp/pkg/apis/kdmp/v1alpha1"
	"github.com/portworx/kdmp/pkg/executor"
	"github.com/portworx/kdmp/pkg/kopia"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newDeleteCommand() *cobra.Command {
	var (
		snapshotID                  string
		credSecretName              string
		credSecretNamespace         string
		volumeBackupDeleteName      string
		volumeBackupDeleteNamespace string
	)
	deleteCommand := &cobra.Command{
		Use:   "delete",
		Short: "delete a backup snapshot",
		Run: func(c *cobra.Command, args []string) {
			executor.HandleErr(runDelete(snapshotID, volumeBackupDeleteName, volumeBackupDeleteNamespace))
		},
	}
	deleteCommand.Flags().StringVar(&snapshotID, "snapshot-id", "", "snapshot ID for kopia backup snapshot that need to be deleted")
	deleteCommand.Flags().StringVar(&credSecretName, "cred-secret-name", "", " cred secret name for kopia backup snapshot that need to be deleted")
	deleteCommand.Flags().StringVar(&credSecretNamespace, "cred-secret-namespace", "", "cred secret namespace for kopia backup snapshot that need to be deleted")
	deleteCommand.Flags().StringVar(&volumeBackupDeleteName, "volume-backup-delete-name", "", "volumeBackupdelete CR name for kopia backup snapshot that need to be deleted")
	deleteCommand.Flags().StringVar(&volumeBackupDeleteNamespace, "volume-backup-delete-namespace", "", "volumeBackupdelete CR namespace for kopia backup snapshot that need to be deleted")
	return deleteCommand
}

func runDelete(snapshotID, volumeBackupDeleteName, volumeBackupDeleteNamespace string) error {
	// Parse using the mounted secrets
	fn := "runDelete:"
	repo, rErr := executor.ParseCloudCred()
	if rErr != nil {
		errMsg := fmt.Sprintf("failed in parsing backuplocation: %s", rErr)
		logrus.Errorf("%s %v", fn, errMsg)
		if err := executor.WriteVolumeBackupDeleteStatus(kdmpapi.VolumeBackupDeleteStatusFailed, errMsg, volumeBackupDeleteName, volumeBackupDeleteNamespace); err != nil {
			errMsg := fmt.Sprintf("failed in updating VolumeBackupDelete CR [%s:%s]: %v", volumeBackupDeleteName, volumeBackupDeleteNamespace, err)
			logrus.Errorf("%v", errMsg)
			return fmt.Errorf(errMsg)
		}
		return fmt.Errorf(errMsg)
	}

	repo.Name = frameBackupPath()

	if err := runKopiaRepositoryConnect(repo); err != nil {
		errMsg := fmt.Sprintf("repository [%v] connect failed: %v", repo.Name, err)
		logrus.Errorf("%s: %v", fn, errMsg)
		if err := executor.WriteVolumeBackupDeleteStatus(kdmpapi.VolumeBackupDeleteStatusFailed, errMsg, volumeBackupDeleteName, volumeBackupDeleteNamespace); err != nil {
			errMsg := fmt.Sprintf("failed in updating VolumeBackupDelete CR [%s:%s]: %v", volumeBackupDeleteName, volumeBackupDeleteNamespace, err)
			logrus.Errorf("%v", errMsg)
			return fmt.Errorf(errMsg)
		}
		return fmt.Errorf(errMsg)
	}

	if err := runKopiaDelete(repo, snapshotID); err != nil {
		errMsg := fmt.Sprintf("snapshot [%v] delete failed: %v", snapshotID, err)
		logrus.Errorf("%s: %v", fn, errMsg)
		if err := executor.WriteVolumeBackupDeleteStatus(kdmpapi.VolumeBackupDeleteStatusFailed, errMsg, volumeBackupDeleteName, volumeBackupDeleteNamespace); err != nil {
			errMsg := fmt.Sprintf("failed in updating VolumeBackupDelete CR [%s:%s]: %v", volumeBackupDeleteName, volumeBackupDeleteNamespace, err)
			logrus.Errorf("%v", errMsg)
			return fmt.Errorf(errMsg)
		}
		return fmt.Errorf(errMsg)
	}

	return nil
}

func runKopiaDelete(repository *executor.Repository, snapshotID string) error {
	fn := "runKopiaDelete:"
	deleteCmd, err := kopia.GetDeleteCommand(
		snapshotID,
	)
	if err != nil {
		errMsg := fmt.Sprintf("getting delete backup snapshot command for snapshot ID [%v] failed: %v", snapshotID, err)
		logrus.Errorf("%s %v", fn, errMsg)
		return fmt.Errorf(errMsg)
	}
	initExecutor := kopia.NewDeleteExecutor(deleteCmd)
	if err := initExecutor.Run(); err != nil {
		errMsg := fmt.Sprintf("running delete backup snapshot command for snapshotID [%v] failed: %v", snapshotID, err)
		logrus.Errorf("%s %v", fn, errMsg)
		return fmt.Errorf(errMsg)
	}

	for {
		time.Sleep(progressCheckInterval)
		status, err := initExecutor.Status()
		if err != nil {
			return err
		}
		if status.LastKnownError != nil {
			return status.LastKnownError
		}

		if status.Done {
			break
		}
	}
	return nil
}
