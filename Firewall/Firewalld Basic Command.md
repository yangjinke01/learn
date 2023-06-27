# Firewalld Basic Command

```text
firewall-cmd -h

Usage: firewall-cmd [OPTIONS...]

General Options
  -h, --help           Prints a short help text and exists
  -V, --version        Print the version string of firewalld
  -q, --quiet          Do not print status messages

Status Options
  --state              Return and print firewalld state
  --reload             Reload firewall and keep state information
  --complete-reload    Reload firewall and lose state information
  --runtime-to-permanent
                       Create permanent from runtime configuration
  --check-config       Check permanent configuration for errors

Permanent Options
  --permanent          Set an option permanently

Zone Options
  --get-default-zone   Print default zone for connections and interfaces
  --set-default-zone=<zone>
  --get-active-zones   Print currently active zones
  --get-zones          Print predefined zones [P]
  --get-services       Print predefined services [P]
  --get-icmptypes      Print predefined icmptypes [P]
  --get-zone-of-interface=<interface>
                       Print name of the zone the interface is bound to [P]
  --get-zone-of-source=<source>[/<mask>]|<MAC>|ipset:<ipset>
                       Print name of the zone the source is bound to [P]
  --list-all-zones     List everything added for or enabled in all zones [P]
  --new-zone=<zone>    Add a new zone [P only]
  --new-zone-from-file=<filename> [--name=<zone>]
                       Add a new zone from file with optional name [P only]
  --delete-zone=<zone> Delete an existing zone [P only]
  --load-zone-defaults=<zone>
                       Load zone default settings [P only] [Z]
  --zone=<zone>        Use this zone to set or query options, else default zone
                       Usable for options marked with [Z]
  --get-target         Get the zone target [P only] [Z]
  --set-target=<target>
                       Set the zone target [P only] [Z]
  --info-zone=<zone>   Print information about a zone
  --path-zone=<zone>   Print file path of a zone [P only]

Service Options
  --new-service=<service>
                       Add a new service [P only]
  --new-service-from-file=<filename> [--name=<service>]
                       Add a new service from file with optional name [P only]
  --delete-service=<service>
                       Delete an existing service [P only]
  --load-service-defaults=<service>
                       Load icmptype default settings [P only]
  --info-service=<service>
                       Print information about a service
  --path-service=<service>
                       Print file path of a service [P only]
  --service=<service> --set-description=<description>
                       Set new description to service [P only]
  --service=<service> --get-description
                       Print description for service [P only]
  --service=<service> --set-short=<description>
                       Set new short description to service [P only]
  --service=<service> --get-short
                       Print short description for service [P only]
  --service=<service> --add-port=<portid>[-<portid>]/<protocol>
                       Add a new port to service [P only]
  --service=<service> --remove-port=<portid>[-<portid>]/<protocol>
                       Remove a port from service [P only]
  --service=<service> --query-port=<portid>[-<portid>]/<protocol>
                       Return whether the port has been added for service [P only]
  --service=<service> --get-ports
                       List ports of service [P only]
  --service=<service> --add-protocol=<protocol>
                       Add a new protocol to service [P only]
  --service=<service> --remove-protocol=<protocol>
                       Remove a protocol from service [P only]
  --service=<service> --query-protocol=<protocol>
                       Return whether the protocol has been added for service [P only]
  --service=<service> --get-protocols
                       List protocols of service [P only]
  --service=<service> --add-source-port=<portid>[-<portid>]/<protocol>
                       Add a new source port to service [P only]
  --service=<service> --remove-source-port=<portid>[-<portid>]/<protocol>
                       Remove a source port from service [P only]
  --service=<service> --query-source-port=<portid>[-<portid>]/<protocol>
                       Return whether the source port has been added for service [P only]
  --service=<service> --get-source-ports
                       List source ports of service [P only]
  --service=<service> --add-module=<module>
                       Add a new module to service [P only]
  --service=<service> --remove-module=<module>
                       Remove a module from service [P only]
  --service=<service> --query-module=<module>
                       Return whether the module has been added for service [P only]
  --service=<service> --get-modules
                       List modules of service [P only]
  --service=<service> --set-destination=<ipv>:<address>[/<mask>]
                       Set destination for ipv to address in service [P only]
  --service=<service> --remove-destination=<ipv>
                       Disable destination for ipv i service [P only]
  --service=<service> --query-destination=<ipv>:<address>[/<mask>]
                       Return whether destination ipv is set for service [P only]
  --service=<service> --get-destinations
                       List destinations in service [P only]
```

## Understanding Predefined Services

A service can be a list of local ports and destinations as well as a list of firewall helper modules automatically loaded if a service is enabled. The use of predefined services makes it easier for the user to enable and disable access to a service. Using the predefined services, or custom defined services, as opposed to opening ports or ranges or ports may make administration easier. Service configuration options and generic file information are described in the *firewalld.service(5)* man page. The services are specified by means of individual XML configuration files which are named in the following format: *service-name.xml*.

To list the default predefined services available using the command line, issue the following command as root:

```shell
ls /usr/lib/firewalld/services/
```

Files in **/usr/lib/firewalld/services/** must not be edited. Only the files in **/etc/firewalld/services/** should be edited.
To list the system or user created services, issue the following command as root:

```shell
ls /etc/firewalld/services/
```

Services can be added and removed using the graphical firewall-config tool and by editing the XML files in **/etc/firewalld/services/**. If a service has not be added or changed by the user, then no corresponding XML file will be found in **/etc/firewalld/services/**. The files **/usr/lib/firewalld/services/** can be used as templates if you wish to add or change a service. As root, issue a command in the following format:

```shell
cp /usr/lib/firewalld/services/[service].xml /etc/firewalld/services/[service].xml
```

You may then edit the newly created file. firewalld will prefer files in **/etc/firewalld/services/** but will fall back to **/usr/lib/firewalld/services/** should a file be deleted, but only after a reload.

## Add the rule to both the permanent and runtime sets

```shell
firewall-cmd --zone=public --add-service=http --permanent

# The reload command drops all runtime configurations and applies a permanent configuration. Because firewalld manages the ruleset dynamically, it won’t break an existing connection and session.
firewall-cmd --reload
```

## Understanding The Direct Interface

firewalld has a so called direct interface, which enables directly passing rules to iptables, ip6tables and ebtables. It is intended for use by applications and not users. It is dangerous to use the direct interface if you are not very familiar with iptables as you could inadvertently cause a breach in the firewall. firewalld still tracks what has been added, so it is still possible to query firewalld and see the changes made by an application using the direct interface mode. The direct interface is used by adding the --direct option to firewall-cmd.
The direct interface mode is intended for services or applications to add specific firewall rules during run time. The rules are not permanent and need to be applied every time after receiving the start, restart or reload message from firewalld using D-BUS.
