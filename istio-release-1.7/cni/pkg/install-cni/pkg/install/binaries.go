package install

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/coreos/etcd/pkg/fileutil"

	"istio.io/istio/cni/pkg/install-cni/pkg/constants"
	"istio.io/istio/cni/pkg/install-cni/pkg/util"
	"istio.io/pkg/log"
)

func copyBinaries(updateBinaries bool, skipBinaries []string) error {
	srcDir := constants.CNIBinDir
	targetDirs := []string{constants.HostCNIBinDir, constants.SecondaryBinDir}

	skipBinariesSet := arrToSet(skipBinaries)

	for _, targetDir := range targetDirs {
		if fileutil.IsDirWriteable(targetDir) != nil {
			log.Infof("Directory %s is not writable, skipping.", targetDir)
			continue
		}

		files, err := ioutil.ReadDir(srcDir)
		if err != nil {
			return err
		}

		for _, file := range files {
			filename := file.Name()
			if skipBinariesSet[filename] {
				log.Infof("%s is in SKIP_CNI_BINARIES, skipping.", filename)
				continue
			}

			targetFilepath := filepath.Join(targetDir, filename)
			if _, err := os.Stat(targetFilepath); err == nil && !updateBinaries {
				log.Infof("%s is already here and UPDATE_CNI_BINARIES isn't true, skipping", targetFilepath)
				continue
			}

			srcFilepath := filepath.Join(srcDir, filename)
			err := util.AtomicCopy(srcFilepath, targetDir, filename)
			if err != nil {
				return err
			}
			log.Infof("Copied %s to %s.", filename, targetDir)
		}
	}

	return nil
}

func arrToSet(array []string) map[string]bool {
	set := make(map[string]bool)
	for _, v := range array {
		set[v] = true
	}
	return set
}
