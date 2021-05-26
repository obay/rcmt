package rcmtssh

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/lkbhargav/go-scp"
	"github.com/lkbhargav/go-scp/auth"
	"github.com/mitchellh/go-homedir"
	"github.com/obay/rcmt/rcmthost"
	"golang.org/x/crypto/ssh"
)

// Prase this root@192.168.100.133:22 into a struct
func ParseConnectionString(cs string) (sshConnectionParameters rcmthost.HostDetails) {
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

func getRSAPublicKeyPath() (rsaPublicKeyPath string, err error) {
	d, err := homedir.Dir()
	if err != nil {
		log.Fatalf("unable to retrieve home directory: %v", err)
	}
	rsaPublicKeyPath = d + "/.ssh/id_rsa"
	return
}

func getRSAPublicKey() (rsaPublicKey []byte, err error) {
	rsaPublicKeyPath, err := getRSAPublicKeyPath()
	rsaPublicKey, err = ioutil.ReadFile(rsaPublicKeyPath)
	return
}

func ConnectToHost(user, host, port string) (client *ssh.Client, session *ssh.Session, err error) {
	var hostKey ssh.PublicKey
	// A public key may be used to authenticate against the remote
	// server by using an unencrypted PEM-encoded private key file.
	//
	// If you have an encrypted private key, the crypto/x509 package
	// can be used to decrypt it.

	rsaPublicKey, err := getRSAPublicKey()
	if err != nil {
		return nil, nil, err
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(rsaPublicKey)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	client, err = ssh.Dial("tcp", host+":"+port, sshConfig)
	if err != nil {
		return nil, nil, err
	}

	session, err = client.NewSession()
	if err != nil {
		client.Close()
		return nil, nil, err
	}
	/***********************************************************************************/
	return
}

func RunRemoteCommandOnRemoteHost(session *ssh.Session, commandline string) (stdout string, stderr string, exitcode int) {
	out, err := session.CombinedOutput(commandline)
	if err != nil {
		if v, ok := err.(*ssh.ExitError); ok {
			exitcode = v.ExitStatus()
			stderr = v.Msg()
		}
	}
	stdout = string(out)
	/***********************************************************************************/
	return
}

func RunSimpleCommandOnRemoteHost(session *ssh.Session, commandline string) {
	stdout, stderr, exitcode := RunRemoteCommandOnRemoteHost(session, commandline)
	if exitcode != 0 {
		os.Stdout.WriteString(stdout)
		os.Stderr.WriteString(stderr)
		os.Exit(exitcode)
	}
	fmt.Print(stdout)
}

func SCP(hostDetails rcmthost.HostDetails, sourceFile string, destinationFile string, destinationFileMode string) (err error) {
	// Use SSH key authentication from the auth package
	// we ignore the host key in this example, please change this if you use this library
	rsaPublicKeyPath, err := getRSAPublicKeyPath()
	if err != nil {
		return
	}
	clientConfig, err := auth.PrivateKey(hostDetails.Username, rsaPublicKeyPath, ssh.InsecureIgnoreHostKey())
	if err != nil {
		return
	}

	// Create a new SCP client
	client := scp.NewClient(hostDetails.Hostname+":"+hostDetails.Port, &clientConfig)

	// Connect to the remote server
	err = client.Connect()
	if err != nil {
		fmt.Println("Couldn't establish a connection to the remote server ", err)
		return
	}

	// Open a file
	f, err := os.Open(sourceFile)
	if err != nil {
		return
	}

	// Close client connection after the file has been copied
	defer client.Close()

	// Close the file after it has been copied
	defer f.Close()

	// Finaly, copy the file over
	// Usage: CopyFile(fileReader, remotePath, permission)

	err = client.CopyFile(f, destinationFile, destinationFileMode)

	if err != nil {
		fmt.Println("Error while copying file ", err)
	}
	return
}
