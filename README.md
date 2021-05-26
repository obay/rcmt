# rcmt

The following commands are how you run rcmt

`rcmt help`
Show this help page.

`rcmt host list`
List the hosts in the hosts.rcmt file. This is the list of hosts where the taks will be applied to. You can't have multiple hosts in the hosts.rcmt with the same hostname.

`rcmt host add [username@]<hostname>[:port]`
Adds a host to the hosts.rcmt file

`rcmt host remove [username@]<hostname>[:port]`
Removes a host from the hosts.rcmt file

`rcmt resource packge add <packagename> [--state="IsInstalled=yes"]`
Creates a new package resource file in the current directory. If no parameters are provided, the resource file will be created using the default settings. You are expected to verify the resource settings and amend it as needed before applying it.

`rcmt resource packge remove <packagename>`
Removes the resource file related to packagename from the current directory.

`rcmt resource service add <servicename> [--state="IsRunning=yes"]`
Creates a new service resource file in the current directory. If no parameters are provided, the resource file will be created using the default settings. You are expected to verify the resource settings and amend it as needed before applying it.

`rcmt resource service remove <servicename>`
Removes the resource file related to servicename from the current directory.

`rcmt resource file add <filename> [--state="FileName=/etc/modtd"]`
Creates a new file resource file in the current directory. If no parameters are provided, the resource file will be created using the default settings. You are expected to verify the resource settings and amend it as needed before applying it.

`rcmt resource file remove <filename>`
Removes the resource file related to filename from the current directory.

`rcmt do`
Applies the task in the current folder. You can `--auto-approve` if you don't want to be asked for a confirmation.

`rcmt undo`
Undo the tasks in the current folder. You can `--auto-approve` if you don't want to be asked for a confirmation.
