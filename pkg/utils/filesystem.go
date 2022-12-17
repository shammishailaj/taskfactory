package utils

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	//"syscall"
	"gopkg.in/djherbis/times.v1"
	"time"
)

func (u *Utils) FileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func (u *Utils) DeleteFileByAge(path string, minAgeForDeletion int64) (bool,error) {
	u.Log.Infof("Deleting file: %s", path)
	/*fileStat*/_, fileStatErr := os.Stat(path)
	if os.IsNotExist(fileStatErr) {
		u.Log.Errorf("Could not find file %s. %s", path, fileStatErr.Error())
		return false, fileStatErr
	} else {
		t, err := times.Stat(path)
		if err != nil {
			log.Errorf("Error stating file %s using times.Stat(). %s", path, err.Error())
			return false, err
		}
		//tFileCreationTime := fileStat.Sys().(*syscall.Stat_t).Ctimespec
		//fileCreationTime := time.Unix(tFileCreationTime.Sec, tFileCreationTime.Nsec)
		var fileCreationTime time.Time
		if t.HasBirthTime() {
			fileCreationTime = t.BirthTime()
		}
		if t.HasChangeTime() {
			fileCreationTime = t.ChangeTime()
		}
		fileModTime := t.ModTime()
		fileAccessTime := t.AccessTime()

		tCurrentTime := time.Now()
		tCurrentTimeUnix := tCurrentTime.Unix()
		tFileAgeForDeletion := minAgeForDeletion // 10 secs X 60 = 600 secs  OR 10 mins

		if (tCurrentTimeUnix - fileCreationTime.Unix()) >= tFileAgeForDeletion {
			delFileErr := os.Remove(path)
			if delFileErr != nil {
				u.Log.Errorf("FAILED to remove file %s. %s", path, delFileErr.Error())
				return false, delFileErr
			} else {
				u.Log.Infof("Successfully Removed file %s", path)
				return true, nil
			}
		} else {
			u.Log.Errorf("Specified File: %s is newer than specified deletion age of - %d second(s). WONT DELETE!", path, minAgeForDeletion)
			u.Log.Errorf("File: %s, Current Time (Unix): %d, Creation Time (Unix): %d, Age for Deletion: %d", path, tCurrentTimeUnix, fileCreationTime.Unix(), tFileAgeForDeletion)
			u.Log.Errorf("Creation Time: %s, Modification Time: %s, Access Time: %s", fileCreationTime.String(), fileModTime.String(), fileAccessTime.String())
			return false, errors.New(fmt.Sprintf("Specified File: %s is newer than specified deletion age of - %d second(s). WONT DELETE!", path, minAgeForDeletion))
		}
	}
}


func (u *Utils) GetFileList( directoryPath string ) map[int]string {
	filesList := make(map[int]string)
	filesListIterator := 0

	u.Log.Infof("directoryPath = %s", directoryPath)
	files, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		u.Log.Errorf("Error reading directories: %s", err.Error())
		return filesList
	}

	for fileKey, file := range files {
		fileStat, fileStatErr := os.Stat(file.Name())
		if fileStatErr != nil {
			u.Log.Errorf("Error stating file #%d: %s. %s", fileKey, file.Name(), fileStatErr.Error())
			continue
		}

		if fileStat.IsDir() {
			tFilesList := u.GetFileList(file.Name())
			tFilesListLen := len(tFilesList)
			if tFilesListLen > 0 {
				for _, value := range tFilesList {
					filesList[filesListIterator] = value
					filesListIterator++
				}
			}
			continue
		}

		filesList[filesListIterator] = file.Name()
		filesListIterator++
	}
	return filesList
}

func (u *Utils) GetFileListGlob( pattern string ) map[int]string {
	filesList := make(map[int]string)
	filesListIterator := 0

	u.Log.Infof("Globbing pattern: %s", pattern)
	files, filesErr := filepath.Glob(pattern)
	if filesErr != nil {
		u.Log.Errorf("Error reading directories: %s", filesErr.Error())
		return filesList
	}

	for fileKey, file := range files {
		fileStat, fileStatErr := os.Stat(file)
		if fileStatErr != nil {
			u.Log.Errorf("Error stating file #%d: %s. %s", fileKey, file, fileStatErr.Error())
			continue
		}

		if fileStat.IsDir() {
			log.Printf("Pattern = %s", pattern)
			patternBase := filepath.Base(pattern)
			log.Printf("Pattern Base = %s", patternBase)
			patternPath := strings.TrimSuffix(pattern, patternBase)
			log.Printf("pattern path = %s", patternPath)
			newDirWithPattern := fmt.Sprintf("%s/%s", file, patternBase)
			log.Printf("New Dir Pattern = %s", newDirWithPattern)

			tFilesList := u.GetFileListGlob(newDirWithPattern)
			tFilesListLen := len(tFilesList)
			if tFilesListLen > 0 {
				for _, value := range tFilesList {
					filesList[filesListIterator] = value
					filesListIterator++
				}
			}
			continue
		}

		filesList[filesListIterator] = file
		filesListIterator++
	}
	return filesList
}

func (u *Utils) DeleteFileList( fileExtToClean, directoryPath string ) int {
	deletionCount := 0
	if fileExtToClean == "" {
		fileExtToClean = "pdf"
	}

	fileExtToClean = fmt.Sprintf("*%s", fileExtToClean)

	if directoryPath == "" {
		u.Log.Errorf("Must specify Directory Path where files need to deleted")
	} else {
		if u.FileExists( directoryPath ) {
			filesList := u.GetFileList(fmt.Sprintf("%s/%s", directoryPath, fileExtToClean) )
			filesListLen := len(filesList)
			if filesListLen > 0 {
				for _, filePath := range filesList {
					fileDel, fileDelErr := u.DeleteFileByAge(filePath, 600)
					if fileDelErr == nil {
						if fileDel {
							deletionCount++
						}
					} else {
						u.Log.Errorf("Error deleting file %s. %s", filePath, fileDelErr.Error())
					}
				}
			} else {
				u.Log.Infof("No files found in directory %s", directoryPath)
			}
		} else {
			u.Log.Infof("Path %s does not exist", directoryPath)
		}
	}

	u.Log.Infof("Deleted %d %s files at path %s", deletionCount, fileExtToClean, directoryPath)
	return deletionCount
}


func (u *Utils) DeleteFileListE(prefix, fileExtToClean, directoryPath string, maxAge int64, forceNoExt bool) int {
	deletionCount := 0
	cpuCores := u.CPUCores()

	u.Log.Infof("LoadAvgCheck():: Found %d CPU Cores", cpuCores)
	if !forceNoExt && fileExtToClean == "" {
		u.Log.Infof("Extensionless option is not being force with an empty/nil file extension. Defaulting to .PDF")
		fileExtToClean = "pdf"
	}
	globPatternFormat := "%s*.%s"
	if forceNoExt {
		u.Log.Infof("--force-no-ext is true, removing . (dot) from globPatternFormat")
		globPatternFormat = "%s*%s"
	}

	fileExtToClean = fmt.Sprintf(globPatternFormat, prefix, fileExtToClean)

	if directoryPath == "" {
		u.Log.Errorf("Must specify Directory Path where files need to deleted")
	} else {
		if u.FileExists( directoryPath ) {
			filesList := u.GetFileListGlob(fmt.Sprintf("%s/%s", directoryPath, fileExtToClean))
			filesListLen := len(filesList)
			u.Log.Infof("Found %d files", filesListLen)
			if filesListLen > 0 {
				for _, filePath := range filesList {
					if u.LoadAvgCheckCPUCores(cpuCores) == LAVG_TREND_NORMAL {
						fileDel, fileDelErr := u.DeleteFileByAge(filePath, maxAge)
						if fileDelErr == nil {
							if fileDel {
								deletionCount++
							}
						} else {
							u.Log.Errorf("Error deleting file %s. %s", filePath, fileDelErr.Error())
						}
					} else {
						loadAvg, loadAvgErr := u.LoadAvg()
						if loadAvgErr != nil {
							u.Log.Errorf("Unable to read System Load Average. %s", loadAvgErr.Error())
						} else {
							u.Log.Infof("Load Average (1), (5), (15) = (%f), (%f), (%f)", loadAvg.Load1, loadAvg.Load5, loadAvg.Load15)
						}
						u.Log.Infof("Deleted %d file(s) till now", deletionCount)
						u.Log.Infof("Sleeping for %d seconds...", DEL_FILE_LIST_E_SLEEP)
						time.Sleep(DEL_FILE_LIST_E_SLEEP * time.Second)
						u.Log.Infof("Woke-up!!!")
					}
				}
			} else {
				u.Log.Infof("No files found in directory %s", directoryPath)
			}
		} else {
			u.Log.Infof("Path %s does not exist", directoryPath)
		}
	}

	u.Log.Infof("Deleted %d %s files at path %s", deletionCount, fileExtToClean, directoryPath)
	return deletionCount
}

