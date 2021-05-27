# rcmt (Rudimentary Configuration Management Tool)

## Building & Releasing

rcmt is using [GoReleaser](https://goreleaser.com) In order to build and release this code, simply run the following commands to build and release.

rcmt uses [Homebrew](https://brew.sh) to publish the application. You will need your own homebrew-tap repo in your own GitHub account in order for [GoReleaser](https://goreleaser.com) to create the necessary files there as well.

```bash
export RELEASE="v0.15.0"
export GITHUB_TOKEN="YOUR_OWN_GITHUB_TOKEN"
git tag -a $RELEASE -m "Release $RELEASE"
git push origin $RELEASE
goreleaser release --rm-dist

```

## Installation
You can install rcmt on any MacOS or Linux machine with [Homebrew](https://brew.sh) installed.

You can install [Homebrew](https://brew.sh) using the following command:

```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

Once [Homebrew](https://brew.sh) is installed, run the following commands to install the tool:

```bash
brew tap obay/tap
brew install rcmt
```

Confirm the tool is installed by running:

```bash
rcmt version
```

## Usage

## SSH Access
Lets say you have 2 sparkling new Ubuntu 18.04.5 servers with SSH root access enabled. To make it easy to follow up, I'll define them here:

```bash
HOST1=192.168.100.133
HOST2=192.168.100.134
```

Now, let's also assume that you have run `ssh-copy-id root@$HOST1` & `ssh-copy-id root@$HOST2` on both hosts and copied your public key so you now can login as root with no need for password prompts.

## Folder & File Structure

Now that you have your hosts ready, let's create the resources to install an Apache webserver with PHP and create a simple Hello, World PHP page.

rcmt requires that you create a YAML file for each resource. Execution of those resources will happen in a lexicographic order, so let's keep that in mind.

rcmt also does not require you to look for the documentation of how to write the resources. You can easily create stubs for them which you can then edit using your favorite editor (Vi of course).

But first, we need to create the list of hosts that we will target with our tool. rcmt allows you to easily create the hosts file as well. Simply run the following commands:

```bash
rcmt host add -n web1 $HOST1
rcmt host add -n web1 $HOST2
```

To confirm those hosts are added correctly, you can run the following command:
```bash
rcmt host list
```

You should see something like this:
```bash
NAME	HOSTNAME     	USERNAME	PORT
web1	192.168.100.133	root    	22
web1	192.168.100.134	root    	22
```

Now, let's add the resource file stubs.

### Adding a Package Resource File
Run the following command to create a new resource file for a package
```bash
rcmt resource package add php
```
You should always check the details of the generated resource file to make sure they match what you want to do.

### Create Template File
rcmt looks at a folder called "templates" in the current working directory. If you don't have a templates directory created, make sure you create one and place the following file into it.
```bash
mkdir templates
cat > ./templates/index.php << EOF
<?php
header("Content-Type: text/plain");
echo "Hello, world!\n";
EOF
```

### Adding a File Resource File
We will need to manage two files in order to see the right index page in those new Apache servers. We need to manage the default index.html (have it deleted/created) and we need to manage the index.php (have it copied).

We will create 2 resource files and we will edit their settings to make sure they do what we want.

```bash
# Create a resource for the default html file
rcmt resource file add index.html
# edit resource_file_index.html.rcmt and:
# * Change filename value from "index.html" to "/var/www/html/index.html"
# * Change exists value from "true" to "false"
# If you don't delete this file, the PHP file will not show as HTML will take precedence over PHP in default Apache settings
```

```bash
# Create a resource for the new PHP file
rcmt resource file add index.php
# edit resource_file_index.php.rcmt and:
# Change filename value from "index.php" to "/var/www/html/index.php"
# Change filetemplatename value from "" to "index.php"
```

### Run Things in Order
Now I mentioned at the begining that rcmt will execute those resources in lexicographic order, so we need to make sure the file names are ordered accordingly. Not pretty, I know. But the "r" in rcmt stands for rudimentary 😬.

```bash
mv resource_package_php.rcmt 1.resource_package_php.rcmt
mv resource_file_index.html.rcmt 2.resource_file_index.html.rcmt
mv resource_file_index.php.rcmt 3.resource_file_index.php.rcmt
```

### Action Time!
Just Do It!

```bash
rcmt do
```

### Verify Your Work
If all went well, running `curl -sv http://$HOST1` should give you the following reply:

```bash
*   Trying 192.168.100.133...
* TCP_NODELAY set
* Connected to 192.168.100.133 (192.168.100.133) port 80 (#0)
> GET / HTTP/1.1
> Host: 192.168.100.133
> User-Agent: curl/7.64.1
> Accept: */*
>
< HTTP/1.1 200 OK
< Date: Wed, 26 May 2021 17:40:54 GMT
< Server: Apache/2.4.29 (Ubuntu)
< Content-Length: 14
< Content-Type: text/plain;charset=UTF-8
<
Hello, world!
* Connection #0 to host 192.168.100.133 left intact
* Closing connection 0
```

## How it Looks Like?
[![asciicast](https://asciinema.org/a/s9VJKkOFj4CpdlhOmjer2ldBb.svg)](https://asciinema.org/a/s9VJKkOFj4CpdlhOmjer2ldBb)
