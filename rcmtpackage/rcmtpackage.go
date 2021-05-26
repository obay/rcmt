package rcmtpackage

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/obay/rcmt/helpers"
	"github.com/obay/rcmt/rcmthost"
	"github.com/obay/rcmt/rcmtservice"
	"github.com/obay/rcmt/rcmtssh"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
)

type PackageResource struct {
	Type         string
	Name         string
	currentState packageState
	DesiredState packageState
}

type packageState struct {
	PackageName     string
	IsInstalled     bool
	RelatedServices []string
}

func GetPackageCurrentState(hostDetails rcmthost.HostDetails, packageName string) (packageCurrentState packageState, err error) {
	packageCurrentState.PackageName = packageName
	/***********************************************************************************/
	_, session, err := rcmtssh.ConnectToHost(hostDetails.Username, hostDetails.Hostname, hostDetails.Port)
	helpers.Check(err)
	commandLine := "apt -qq list " + packageName
	stdout, _, _ := rcmtssh.RunRemoteCommandOnRemoteHost(session, commandLine)
	packageCurrentState.IsInstalled = strings.Contains(stdout, "[installed]")
	/***********************************************************************************/
	session.Close()
	return
}

func restartServices(hostDetails rcmthost.HostDetails, fileDesiredState packageState) (err error) {
	if len(fileDesiredState.RelatedServices) != 0 {
		helpers.PrintWarningf("restarting related services... ")
		for _, service := range fileDesiredState.RelatedServices {
			helpers.PrintWarningf("[" + service + "] ")
			rcmtservice.RestartServiceOnRemoteHost(hostDetails, service)
		}
	}
	return
}

func ConvergePackageState(hostDetails rcmthost.HostDetails, packageCurrentState, packageDesiredState packageState) (err error) {
	/*********************************************************************************************/
	_, session, err := rcmtssh.ConnectToHost(hostDetails.Username, hostDetails.Hostname, hostDetails.Port)
	helpers.Check(err)
	/*********************************************************************************************/
	if packageCurrentState.PackageName != packageDesiredState.PackageName {
		helpers.PrintError("While you should not be reading this as an end user, if you are, please note that you can't converge one package to be another. Install or uninstall each separately. And please send us an email that you came across this error!")
	}
	if packageCurrentState.IsInstalled != packageDesiredState.IsInstalled {
		helpers.PrintWarningf("Current state and desired state for package \"" + packageDesiredState.PackageName + "\" on " + hostDetails.Hostname + " are different. Converging... ")
		if packageDesiredState.IsInstalled {
			helpers.PrintWarningf("Installing " + packageDesiredState.PackageName + "... ")
			installPackageOnRemoteHost(session, packageDesiredState.PackageName)
			err = restartServices(hostDetails, packageDesiredState)
			if err != nil {
				return
			}
			helpers.PrintWarning("Done!")
		} else {
			helpers.PrintWarningf("Uninstalling " + packageDesiredState.PackageName + "... ")
			uninstallPackageFromRemoteHost(session, packageDesiredState.PackageName)
			err = restartServices(hostDetails, packageDesiredState)
			if err != nil {
				return
			}
			helpers.PrintWarning("Done!")
		}
	} else {
		helpers.PrintSuccess("Current state and desired state for package \"" + packageDesiredState.PackageName + "\" on " + hostDetails.Hostname + " are the same. Nothing to do.")
	}
	/*********************************************************************************************/
	session.Close()
	return
}

func installPackageOnRemoteHost(session *ssh.Session, packageName string) {
	commandLine := "apt update && apt install -y " + packageName
	rcmtssh.RunSimpleCommandOnRemoteHost(session, commandLine)
}

func uninstallPackageFromRemoteHost(session *ssh.Session, packageName string) {
	commandLine := "apt remove -y " + packageName
	rcmtssh.RunSimpleCommandOnRemoteHost(session, commandLine)
}

func UnmarshalPackageResource(yamlBlock string) (packageResource PackageResource, err error) {
	bytes := []byte(yamlBlock)
	err = yaml.Unmarshal(bytes, &packageResource)
	return
}

func MrshalPackageResource(packageResource PackageResource) (yamlBlock string, err error) {
	bytes, err := yaml.Marshal(&packageResource)
	yamlBlock = string(bytes)
	return
}

func AddPackageUsingPackageResource(newPackageResource PackageResource) (err error) {
	yamlBlock, err := MrshalPackageResource(newPackageResource)
	if err == nil {
		rcmtFilename := "resource_package_" + newPackageResource.DesiredState.PackageName + ".rcmt"
		if _, oserr := os.Stat(rcmtFilename); !os.IsNotExist(oserr) {
			return errors.New("a file for this resource already exist")
		}
		err = ioutil.WriteFile(rcmtFilename, []byte(yamlBlock), 0644)
	}
	return
}

func AddPackage(packageName string) (err error) {
	var newPackageResource PackageResource
	newPackageResource.Type = "package"
	newPackageResource.Name = packageName
	newPackageResource.DesiredState.PackageName = packageName
	newPackageResource.DesiredState.IsInstalled = true
	err = AddPackageUsingPackageResource(newPackageResource)
	return
}

func RemovePackage(packageName string) (err error) {
	rcmtFilename := "./resource_package_" + packageName + ".rcmt"
	if _, oserr := os.Stat(rcmtFilename); os.IsNotExist(oserr) {
		return errors.New("this resource file doesn't exist")
	}
	err = os.Remove(rcmtFilename)
	return
}

/** Interface Implementation *********************************************************************/

func (r PackageResource) ResourceType() string {
	return r.Type
}

func (r PackageResource) ResourceName() string {
	return r.Name
}

func (r PackageResource) ResourceCurrentState() string {
	return helpers.ConvertStructToString(r.currentState)
}

func (r PackageResource) ResourceDesiredState() string {
	return helpers.ConvertStructToString(r.DesiredState)
}

func (r PackageResource) Converge(hosts []rcmthost.HostDetails) (err error) {
	for _, host := range hosts {
		fileCurrentState, err := GetPackageCurrentState(host, r.DesiredState.PackageName)
		if err != nil {
			break
		}
		err = ConvergePackageState(host, fileCurrentState, r.DesiredState)
		if err != nil {
			break
		}
	}
	return
}

/*************************************************************************************************/
