package rcmthost

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/obay/rcmt/helpers"
	"gopkg.in/yaml.v2"
)

var hostsfilename string = "./hosts.rcmt"

type HostDetails struct {
	Name     string
	Hostname string
	Username string
	Port     string
}

func hostExistAlready(hostname string, currenthosts []HostDetails) bool {
	for _, host := range currenthosts {
		if host.Hostname == hostname {
			return true
		}
	}
	return false
}

func AddHostUsingSSHConnectionParameters(host HostDetails) (err error) {
	err = AddHost(host.Name, host.Hostname, host.Username, host.Port)
	return
}

func AddHost(name, hostname, username, port string) (err error) {
	currenthosts := LoadHosts()
	if hostExistAlready(hostname, currenthosts) {
		return errors.New("Hostname " + hostname + " is already in the hosts list.")
	}
	var newhost = HostDetails{Name: name, Hostname: hostname, Username: username, Port: port}
	currenthosts = append(currenthosts, newhost)
	bytes, err := yaml.Marshal(currenthosts)
	helpers.Check(err)
	err = ioutil.WriteFile("hosts.rcmt", bytes, 0644)
	helpers.Check(err)
	return
}

func RemoveHost(name, hostname, username, port string) (err error) {
	currenthosts := LoadHosts()
	if !hostExistAlready(hostname, currenthosts) {
		return errors.New("Hostname " + hostname + " is not in the hosts list.")
	}

	var newhosts []HostDetails
	for _, host := range currenthosts {
		if host.Hostname != hostname {
			newhosts = append(newhosts, host)
		}
	}
	bytes, err := yaml.Marshal(newhosts)
	helpers.Check(err)
	err = ioutil.WriteFile("hosts.rcmt", bytes, 0644)
	helpers.Check(err)
	return
}

func RemoveHostUsingSSHConnectionParameters(host HostDetails) (err error) {
	err = RemoveHost(host.Name, host.Hostname, host.Username, host.Port)
	return
}

func addLabHosts() {
	AddHost("Ansible1", "192.168.100.133", "", "")
	AddHost("Ansible2", "192.168.100.134", "", "")
}

func addSlackHosts() {
	AddHost("web1", "3.89.118.49", "", "")
	AddHost("web2", "54.198.54.144", "", "")
}

func LoadHosts() (hosts []HostDetails) {
	bytes, err := ioutil.ReadFile(hostsfilename)
	if _, ok := err.(*os.PathError); ok {
		return nil
	}
	helpers.Check(err)
	err = yaml.Unmarshal(bytes, &hosts)
	helpers.Check(err)
	return
}

func SerializeHostsToJSON(hosts []HostDetails) (jsonOutput string, err error) {
	bytes, err := json.Marshal(hosts)
	if err != nil {
		return "", err
	}
	jsonOutput = string(bytes)
	return
}

// Prase this root@192.168.100.133:22 into a struct
func ParseConnectionString(cs string) (sshConnectionParameters HostDetails) {
	sshConnectionParameters.Hostname = cs
	sshConnectionParameters.Username = "root"
	sshConnectionParameters.Port = "22"
	/***********************************************************************************/
	if parts := strings.Split(cs, "@"); len(parts) == 2 {
		sshConnectionParameters.Username = parts[0]
		sshConnectionParameters.Hostname = parts[1]
	}
	if parts := strings.Split(sshConnectionParameters.Hostname, ":"); len(parts) == 2 {
		sshConnectionParameters.Hostname = parts[0]
		sshConnectionParameters.Port = parts[1]
	}
	/***********************************************************************************/
	return
}
