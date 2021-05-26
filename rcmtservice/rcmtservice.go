package rcmtservice

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/obay/rcmt/helpers"
	"github.com/obay/rcmt/rcmthost"
	"github.com/obay/rcmt/rcmtssh"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
)

type ServiceResource struct {
	Type         string
	Name         string
	currentState serviceState
	DesiredState serviceState
}

type serviceState struct {
	ServiceName string
	IsRunning   bool
}

func ConvergeServiceState(sshConnectionParameters rcmthost.HostDetails, serviceCurrentState, serviceDesiredState serviceState) (err error) {
	/*********************************************************************************************/
	_, session, err := rcmtssh.ConnectToHost(sshConnectionParameters.Username, sshConnectionParameters.Hostname, sshConnectionParameters.Port)
	helpers.Check(err)
	/*********************************************************************************************/
	if serviceCurrentState.ServiceName != serviceDesiredState.ServiceName {
		helpers.PrintError("While you should not be reading this as an end user, if you are, please note that you can't converge one service to be another. Run or stop each separately. And please send us an email that you came across this error!")
	}
	if serviceCurrentState.IsRunning != serviceDesiredState.IsRunning {
		helpers.PrintWarningf("Current state and desired state for service \"" + serviceDesiredState.ServiceName + "\" are different. Converging... ")
		if serviceDesiredState.IsRunning {
			helpers.PrintWarningf("Starting " + serviceDesiredState.ServiceName + "... ")
			startServiceOnRemoteHost(session, serviceDesiredState.ServiceName)
			helpers.PrintWarning("Done!")
		} else {
			helpers.PrintWarningf("Uninstalling " + serviceDesiredState.ServiceName + "... ")
			stopServiceOnRemoteHost(session, serviceDesiredState.ServiceName)
			helpers.PrintWarning("Done!")
		}
	} else {
		helpers.PrintSuccess("Current state and desired state for package \"" + serviceDesiredState.ServiceName + "\" are the same. Nothing to do.")
	}
	/*********************************************************************************************/
	session.Close()
	return
}

func startServiceOnRemoteHost(session *ssh.Session, serviceName string) {
	commandLine := "service " + serviceName + " start"
	rcmtssh.RunSimpleCommandOnRemoteHost(session, commandLine)
}

func stopServiceOnRemoteHost(session *ssh.Session, serviceName string) {
	commandLine := "service " + serviceName + " stop"
	rcmtssh.RunSimpleCommandOnRemoteHost(session, commandLine)
}

func RestartServiceOnRemoteHost(hostDetails rcmthost.HostDetails, serviceName string) {
	_, session, err := rcmtssh.ConnectToHost(hostDetails.Username, hostDetails.Hostname, hostDetails.Port)
	helpers.Check(err)
	/*********************************************************************************************/
	commandLine := "service " + serviceName + " restart"
	rcmtssh.RunSimpleCommandOnRemoteHost(session, commandLine)
}

func GetServiceCurrentState(hostDetails rcmthost.HostDetails, serviceName string) (serviceCurrentState serviceState, err error) {
	serviceCurrentState.ServiceName = serviceName
	/***********************************************************************************/
	_, session, err := rcmtssh.ConnectToHost(hostDetails.Username, hostDetails.Hostname, hostDetails.Port)
	// if the service is running, the following command will return: "active (running)"
	commandLine := "systemctl status " + serviceName + " | grep Active | cut -f2 -d':' | cut -f2-3 -d' '"
	stdout, _, _ := rcmtssh.RunRemoteCommandOnRemoteHost(session, commandLine)
	serviceCurrentState.IsRunning = strings.TrimSpace(stdout) == "active (running)"
	/***********************************************************************************/
	session.Close()
	return
}

func LoadServiceResource(yamlBlock string) (serviceResource ServiceResource) {
	bytes := []byte(yamlBlock)
	err := yaml.Unmarshal(bytes, &serviceResource)
	helpers.Check(err)
	return serviceResource
}

func UnmarshalServiceResource(yamlBlock string) (serviceResource ServiceResource, err error) {
	bytes := []byte(yamlBlock)
	err = yaml.Unmarshal(bytes, &serviceResource)
	return
}

func MrshalServiceResource(serviceResource ServiceResource) (yamlBlock string, err error) {
	bytes, err := yaml.Marshal(&serviceResource)
	yamlBlock = string(bytes)
	return
}

func AddServiceUsingServiceResource(newServiceResource ServiceResource) (err error) {
	yamlBlock, err := MrshalServiceResource(newServiceResource)
	if err == nil {
		rcmtFilename := "resource_service_" + newServiceResource.DesiredState.ServiceName + ".rcmt"
		if _, oserr := os.Stat(rcmtFilename); !os.IsNotExist(oserr) {
			return errors.New("a file for this resource already exist")
		}
		err = ioutil.WriteFile(rcmtFilename, []byte(yamlBlock), 0644)
	}
	return
}

func AddService(serviceName string) (err error) {
	var newServiceResource ServiceResource
	newServiceResource.Type = "service"
	newServiceResource.Name = serviceName
	newServiceResource.DesiredState.ServiceName = serviceName
	newServiceResource.DesiredState.IsRunning = true
	err = AddServiceUsingServiceResource(newServiceResource)
	return
}

func RemoveService(serviceName string) (err error) {
	rcmtFilename := "./resource_service_" + serviceName + ".rcmt"
	if _, oserr := os.Stat(rcmtFilename); os.IsNotExist(oserr) {
		return errors.New("this resource file doesn't exist")
	}
	err = os.Remove(rcmtFilename)
	return
}

/** Interface Implementation *********************************************************************/

func (r ServiceResource) ResourceType() string {
	return r.Type
}

func (r ServiceResource) ResourceName() string {
	return r.Name
}

func (r ServiceResource) ResourceCurrentState() string {
	return helpers.ConvertStructToString(r.currentState)
}

func (r ServiceResource) ResourceDesiredState() string {
	return helpers.ConvertStructToString(r.DesiredState)
}

func (r ServiceResource) Converge(hosts []rcmthost.HostDetails) (err error) {
	for _, host := range hosts {
		serviceCurrentState, err := GetServiceCurrentState(host, r.DesiredState.ServiceName)
		if err != nil {
			break
		}
		err = ConvergeServiceState(host, serviceCurrentState, r.DesiredState)
		if err != nil {
			break
		}
	}
	return
}

/*************************************************************************************************/
