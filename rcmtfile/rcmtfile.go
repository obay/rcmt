package rcmtfile

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/obay/rcmt/helpers"
	"github.com/obay/rcmt/rcmthost"
	"github.com/obay/rcmt/rcmtservice"
	"github.com/obay/rcmt/rcmtssh"
	"gopkg.in/yaml.v2"
)

type FileResource struct {
	Type         string
	Name         string
	currentState FileState
	DesiredState FileState
}

type FileState struct {
	FileName         string
	FileTemplateName string
	Exists           bool
	Mode             string
	MD5              string
	IsDirectory      bool
	UsernameOfOwner  string
	GroupNameOfOwner string
	RelatedServices  []string
	// FileBirth time.Time
	// TimeOfLastDataModification time.Time
	// TimeOfLastStatusChange     time.Time
	// FileSizeInByte             int64
	// UserIDOfOwner              int
	// GroupIDOfOwner             int
	// TimeOfLastAccess           time.Time
}

func getFileMD5(hostDetails rcmthost.HostDetails, destinationFile string) (md5hash string, err error) {
	_, session, err := rcmtssh.ConnectToHost(hostDetails.Username, hostDetails.Hostname, hostDetails.Port)
	helpers.Check(err)
	commandLine := "md5sum " + destinationFile + " | cut -f1 -d' '"
	stdout, _, _ := rcmtssh.RunRemoteCommandOnRemoteHost(session, commandLine)
	md5hash = strings.TrimSpace(stdout)
	return
}

func getLocalFileMD5(fileName string) (md5hash string, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer file.Close()

	hash := md5.New()

	//Copy the file in the hash interface and check for any error
	if _, err = io.Copy(hash, file); err != nil {
		return
	}

	//Get the 16 bytes hash
	hashInBytes := hash.Sum(nil)[:16]
	md5hash = fmt.Sprintf("%x", hashInBytes)
	return
}

func SetFileState(hostDetails rcmthost.HostDetails) (err error) {
	// To-do: Implement me
	return
}

// func convertStringtoFileMode(fileModeString string) (fileMode fs.FileMode, err error) {
// 	fileModeUint64, err := strconv.ParseUint(fileModeString, 10, 32)
// 	helpers.Check(err)
// 	fileMode = fs.FileMode(fileModeUint64)
// 	return
// }

func GetFileCurrentState(hostDetails rcmthost.HostDetails, destinationFile string) (fileStat FileState, err error) {
	_, session, err := rcmtssh.ConnectToHost(hostDetails.Username, hostDetails.Hostname, hostDetails.Port)
	if err != nil {
		return
	}

	fileStat.FileName = destinationFile
	commandLine := "stat -c \"FileBirth=%W:FileName=%n:FileSizeInByte=%s:FileType=%F:GroupIDOfOwner=%g:GroupNameOfOwner=%G:Mode=%04a:TimeOfLastAccess=%X:TimeOfLastDataModification=%Y:TimeOfLastStatusChange=%Z:UserIDOfOwner=%u:UsernameOfOwner=%U\" " + destinationFile
	stdout, _, exitcode := rcmtssh.RunRemoteCommandOnRemoteHost(session, commandLine)
	if exitcode == 1 {
		if strings.Contains(stdout, "No such file or directory") {
			fileStat.Exists = false
			return fileStat, nil
		} else {
			log.Fatalln("\"" + stdout + "\"")
		}
		return
	}
	// The standard output of this command will look like this:
	// FileBirth=0:FileName=filex.txt:FileSizeInByte=6:FileType=regular file:GroupIDOfOwner=0:GroupNameOfOwner=root:Mode=0644:TimeOfLastAccess=1621791034:TimeOfLastDataModification=1621791034:TimeOfLastStatusChange=1621792920:UserIDOfOwner=0:UsernameOfOwner=root
	// FileBirth=0:FileName=/:FileSizeInByte=4096:FileType=directory:GroupIDOfOwner=0:GroupNameOfOwner=root:Mode=0755:TimeOfLastAccess=1621793056:TimeOfLastDataModification=1621582124:TimeOfLastStatusChange=1621582124:UserIDOfOwner=0:UsernameOfOwner=root
	// FileBirth=0:FileName=/usr/bin/ls:FileSizeInByte=138856:FileType=regular file:GroupIDOfOwner=0:GroupNameOfOwner=root:Mode=0755:TimeOfLastAccess=1621789518:TimeOfLastDataModification=1551367831:TimeOfLastStatusChange=1621582075:UserIDOfOwner=0:UsernameOfOwner=root

	fileStat.Exists = true
	statLines := strings.Split(stdout, ":")
	for _, statLine := range statLines {
		kv := strings.Split(statLine, "=")
		// FileBirth=0
		// FileName=/root/filex.txt
		// FileSizeInByte=6
		// FileType=regular file
		// GroupIDOfOwner=0
		// GroupNameOfOwner=root
		// Mode=644
		// TimeOfLastAccess=1621791034
		// TimeOfLastDataModification=1621791034
		// TimeOfLastStatusChange=1621792920
		// UserIDOfOwner=0
		// UsernameOfOwner=root
		switch kv[0] {
		// case "FileBirth":
		// 	fileStat.FileBirth = helpers.GetEpochTime(kv[1])
		case "FileName":
			fileStat.FileName = kv[1]
		// case "FileSizeInByte":
		// 	fileStat.FileSizeInByte = helpers.StringToInt64(kv[1])
		case "FileType":
			fileStat.IsDirectory = kv[1] == "directory"
		// case "GroupIDOfOwner":
		// 	fileStat.GroupIDOfOwner = helpers.StringToInt(kv[1])
		case "GroupNameOfOwner":
			fileStat.GroupNameOfOwner = kv[1]
		case "Mode":
			fileStat.Mode = kv[1]
		// case "TimeOfLastAccess":
		// 	fileStat.TimeOfLastAccess = helpers.GetEpochTime(kv[1])
		// case "TimeOfLastDataModification":
		// 	fileStat.TimeOfLastDataModification = helpers.GetEpochTime(kv[1])
		// case "TimeOfLastStatusChange":
		// 	fileStat.TimeOfLastStatusChange = helpers.GetEpochTime(kv[1])
		// case "UserIDOfOwner":
		// 	fileStat.UserIDOfOwner = helpers.StringToInt(kv[1])
		case "UsernameOfOwner":
			fileStat.UsernameOfOwner = strings.TrimSpace(kv[1])
		}
	}
	if fileStat.Exists && !fileStat.IsDirectory {
		fileStat.MD5, err = getFileMD5(hostDetails, destinationFile)
		helpers.Check(err)
	}
	/*********************************************************************************************/
	return
}

func restartServices(hostDetails rcmthost.HostDetails, servicesList []string) (err error) {
	if len(servicesList) > 0 {
		helpers.PrintWarningf("restarting related services... ")
		for _, service := range servicesList {
			helpers.PrintWarningf("[" + service + "] ")
			rcmtservice.RestartServiceOnRemoteHost(hostDetails, service)
		}
	}
	return
}

func ConvergeFileState(hostDetails rcmthost.HostDetails, fileCurrentState, fileDesiredState FileState) (err error) {
	/*********************************************************************************************/
	_, session, err := rcmtssh.ConnectToHost(hostDetails.Username, hostDetails.Hostname, hostDetails.Port)
	helpers.Check(err)
	/*********************************************************************************************/
	// Fill MD5
	if fileDesiredState.MD5 == "" && fileDesiredState.FileTemplateName != "" {
		fileDesiredState.MD5, err = getLocalFileMD5("./templates/" + fileDesiredState.FileTemplateName)
		if err != nil {
			return
		}
	}
	/*********************************************************************************************/
	if fileCurrentState.FileName != fileDesiredState.FileName {
		helpers.PrintError("While you should not be reading this as an end user, if you are, please note that you can't converge one file to be another. Create or delete each separately. And please send us an email that you came across this error!")
	}

	if !fileDesiredState.Exists && fileCurrentState.Exists {
		// file exist but should not. delete file
		helpers.PrintWarningf("File exist but it should not. Deleting \"" + fileDesiredState.FileName + "\"... ")
		err = deleteFileFromRemoteHost(hostDetails, fileDesiredState.FileName)
		if err != nil {
			return
		}
		err = restartServices(hostDetails, fileDesiredState.RelatedServices)
		if err != nil {
			return
		}
		helpers.PrintWarning("Done!")
	} else if !fileCurrentState.Exists && fileDesiredState.Exists {
		// File doesn't exist. create file
		helpers.PrintWarningf("File \"" + fileDesiredState.FileName + "\" doesn't exist on " + hostDetails.Hostname + ". Creating file...")
		err = rcmtssh.SCP(hostDetails, "./templates/"+fileDesiredState.FileTemplateName, fileDesiredState.FileName, fileDesiredState.Mode)
		// err = copyFileToRemoteHost(hostDetails, "./templates/"+fileDesiredState.FileTemplateName, fileDesiredState.FileName, fileDesiredState.Mode)
		if err != nil {
			return
		}
		err = restartServices(hostDetails, fileDesiredState.RelatedServices)
		if err != nil {
			return
		}
		// To-do: apply time & ownershipo for users and groups
		err = setFileOwnershipOnRemoteHost(hostDetails, fileDesiredState.UsernameOfOwner, fileDesiredState.GroupNameOfOwner, fileCurrentState.FileName)
		if err != nil {
			return
		}
		helpers.PrintWarning("Done!")
	} else if filesStatesAreDifferent(fileCurrentState, fileDesiredState) && fileDesiredState.Exists {
		if fileCurrentState.MD5 != fileDesiredState.MD5 {
			// File exist. different content. change file
			helpers.PrintWarningf("File exist but with different content. Copying \"" + fileDesiredState.FileName + "\"...")
			err = rcmtssh.SCP(hostDetails, "./templates/"+fileDesiredState.FileTemplateName, fileDesiredState.FileName, fileDesiredState.Mode)
			if err != nil {
				return
			}
			err = restartServices(hostDetails, fileDesiredState.RelatedServices)
			if err != nil {
				return
			}
			helpers.PrintWarning("Done!")
		}
		if fileCurrentState.GroupNameOfOwner != fileDesiredState.GroupNameOfOwner || fileCurrentState.UsernameOfOwner != fileDesiredState.UsernameOfOwner {
			// File exist. different ownership. set ownerhsip
			helpers.PrintWarningf("File exist but with different ownership. Setting ownership for " + fileDesiredState.FileName + "...")
			err = setFileOwnershipOnRemoteHost(hostDetails, fileDesiredState.UsernameOfOwner, fileDesiredState.GroupNameOfOwner, fileDesiredState.FileName)
			if err != nil {
				return
			}
			helpers.PrintWarning("Done!")
		}
		if fileCurrentState.Mode != fileDesiredState.Mode {
			// File exist. different modes. set file mode
			helpers.PrintWarningf("File exist but with different mode. Setting file mode for \"" + fileDesiredState.FileName + "\"...")
			err = setFileModeOnRemoteHost(hostDetails, fileDesiredState.FileName, fileDesiredState.Mode)
			if err != nil {
				return
			}
			helpers.PrintWarning("Done!")
		}
	} else {
		helpers.PrintSuccess("Current state and desired state for file \"" + fileDesiredState.FileName + "\" on " + hostDetails.Hostname + " are the same. Nothing to do.")
	}
	/*********************************************************************************************/
	session.Close()
	return
}

func filesStatesAreDifferent(fileCurrentState FileState, fileDesiredState FileState) (differentStates bool) {
	differentStates = differentStates || (fileCurrentState.MD5 != fileDesiredState.MD5)
	differentStates = differentStates || (fileCurrentState.GroupNameOfOwner != fileDesiredState.GroupNameOfOwner)
	differentStates = differentStates || (fileCurrentState.UsernameOfOwner != fileDesiredState.UsernameOfOwner)
	differentStates = differentStates || (fileCurrentState.Mode != fileDesiredState.Mode)
	return
}

func deleteFileFromRemoteHost(hostDetails rcmthost.HostDetails, fileName string) (err error) {
	_, session, err := rcmtssh.ConnectToHost(hostDetails.Username, hostDetails.Hostname, hostDetails.Port)
	commandLine := "rm -rf " + fileName
	rcmtssh.RunSimpleCommandOnRemoteHost(session, commandLine)
	session.Close()
	return
}

// func copyFileToRemoteHost(hostDetails rcmthost.HostDetails, sourceFile string, destinationFile string, destinationFileMode string) (err error) {
// 	_, session, err := rcmtssh.ConnectToHost(hostDetails.Username, hostDetails.Hostname, hostDetails.Port)
// 	if err != nil {
// 		return
// 	}
// 	err = scp.CopyPath(sourceFile, destinationFile, session)
// 	if err != nil {
// 		session.Close()
// 		return
// 	}
// 	err = setFileModeOnRemoteHost(hostDetails, destinationFileMode, destinationFile)
// 	return
// }

func setFileModeOnRemoteHost(hostDetails rcmthost.HostDetails, destinationFile string, fileMode string) (err error) {
	_, session, err := rcmtssh.ConnectToHost(hostDetails.Username, hostDetails.Hostname, hostDetails.Port)
	commandLine := "chmod " + fileMode + " " + destinationFile
	rcmtssh.RunSimpleCommandOnRemoteHost(session, commandLine)
	session.Close()
	return
}

func setFileOwnershipOnRemoteHost(hostDetails rcmthost.HostDetails, fileUser string, fileGroup string, destinationFile string) (err error) {
	_, session, err := rcmtssh.ConnectToHost(hostDetails.Username, hostDetails.Hostname, hostDetails.Port)
	commandLine := "chown " + fileUser + ":" + fileGroup + " " + destinationFile
	rcmtssh.RunSimpleCommandOnRemoteHost(session, commandLine)
	session.Close()
	return
}

func UnmarshalFileResource(yamlBlock string) (fileResource FileResource, err error) {
	bytes := []byte(yamlBlock)
	err = yaml.Unmarshal(bytes, &fileResource)
	return
}

func MrshalFileResource(fileResource FileResource) (yamlBlock string, err error) {
	bytes, err := yaml.Marshal(&fileResource)
	yamlBlock = string(bytes)
	return
}

func AddFileUsingFileResource(newFileResource FileResource) (err error) {
	yamlBlock, err := MrshalFileResource(newFileResource)
	if err == nil {
		rcmtFilename := "resource_file_" + newFileResource.DesiredState.FileName + ".rcmt"
		if _, oserr := os.Stat(rcmtFilename); !os.IsNotExist(oserr) {
			return errors.New("a file for this resource already exist")
		}
		err = ioutil.WriteFile(rcmtFilename, []byte(yamlBlock), 0644)
	}
	return
}

func AddFile(fileName string) (err error) {
	var newFileResource FileResource
	newFileResource.Type = "file"
	newFileResource.Name = fileName
	newFileResource.DesiredState.UsernameOfOwner = "root"
	newFileResource.DesiredState.GroupNameOfOwner = "root"
	newFileResource.DesiredState.FileName = fileName
	newFileResource.DesiredState.Mode = "0644"
	newFileResource.DesiredState.Exists = true
	err = AddFileUsingFileResource(newFileResource)
	return
}

func RemoveFile(fileName string) (err error) {
	rcmtFilename := "./resource_file_" + fileName + ".rcmt"
	if _, oserr := os.Stat(rcmtFilename); os.IsNotExist(oserr) {
		return errors.New("this resource file doesn't exist")
	}
	err = os.Remove(rcmtFilename)
	return
}

/** Interface Implementation *********************************************************************/

func (r FileResource) ResourceType() string {
	return r.Type
}

func (r FileResource) ResourceName() string {
	return r.Name
}

func (r FileResource) ResourceCurrentState() string {
	return helpers.ConvertStructToString(r.currentState)
}

func (r FileResource) ResourceDesiredState() string {
	return helpers.ConvertStructToString(r.DesiredState)
}

func (r FileResource) Converge(hosts []rcmthost.HostDetails) (err error) {
	for _, host := range hosts {
		fileCurrentState, err := GetFileCurrentState(host, r.DesiredState.FileName)
		if err != nil {
			return err
		}
		err = ConvergeFileState(host, fileCurrentState, r.DesiredState)
		if err != nil {
			return err
		}
	}
	return
}

/*************************************************************************************************/
