package rcmtresource

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/obay/rcmt/helpers"
	"github.com/obay/rcmt/rcmtfile"
	"github.com/obay/rcmt/rcmthost"
	"github.com/obay/rcmt/rcmtpackage"
	"github.com/obay/rcmt/rcmtservice"
	"gopkg.in/yaml.v2"
)

type ResourceInterface interface {
	ResourceType() string
	ResourceName() string
	ResourceCurrentState() string
	ResourceDesiredState() string
	Converge(hosts []rcmthost.HostDetails) error
}

type ResourceStruct struct {
	Type string
	Name string
}

func LoadResources() (resources []ResourceInterface) {
	rcmtFiles := getAllRCMTFilesInCurrentDirectory()
	for _, rcmtFile := range rcmtFiles {
		if rcmtFile != "hosts.rcmt" {
			resources = append(resources, LoadResource(rcmtFile))
		}
	}
	return
}

func getAllRCMTFilesInCurrentDirectory() (yamlFiles []string) {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".rcmt") {
			yamlFiles = append(yamlFiles, f.Name())
		}
	}
	return
}

// cast the resource type to the correct struct
func LoadResource(rcmtFile string) ResourceInterface {
	resourceType, yamlBlock := discoverRCMTResourceType(rcmtFile)
	var res ResourceInterface
	switch resourceType {
	case "package":
		r, err := rcmtpackage.UnmarshalPackageResource(yamlBlock)
		res = ResourceInterface(r)
		helpers.Check(err)
		// fmt.Printf("%v\n", res)
	case "file":
		r, err := rcmtfile.UnmarshalFileResource(yamlBlock)
		res = ResourceInterface(r)
		helpers.Check(err)
		// fmt.Printf("%v\n", res)
	case "service":
		r, err := rcmtservice.UnmarshalServiceResource(yamlBlock)
		res = ResourceInterface(r)
		helpers.Check(err)
		// fmt.Printf("%v\n", res)
	}
	return res
}

// detect resource type (file / service / package)
// read the rcmt file looking for the resource type field
func discoverRCMTResourceType(f string) (resourceType string, yamlBlock string) {
	bytes, err := ioutil.ReadFile(f)
	helpers.Check(err)
	var resource ResourceStruct
	err = yaml.Unmarshal(bytes, &resource)
	helpers.Check(err)
	return resource.Type, string(bytes)
}

func Do(hosts []rcmthost.HostDetails) (err error) {
	resources := LoadResources()
	for _, resource := range resources {
		err = resource.Converge(hosts)
		if err != nil {
			break
		}
	}
	return
}
