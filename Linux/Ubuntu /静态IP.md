# Static IP

```shell
 vim /etc/netplan/00-installer-config.yaml
```

```yaml
network:
  ethernets:
    enp0s3:
      dhcp4: true
  version: 2
```

Before changing the configuration, let’s explain the code in a short.

Each Netplan Yaml file starts with the `network` key that has at least two required elements. The first required element is the version of the network configuration format, and the second one is the device type. The device type can be `ethernets`, `bonds`, `bridges`, or `vlans`.

Under the device’s type (`ethernets`), you can specify one or more network interfaces. In this example, we have only one interface `ens3` that is configured to obtain IP addressing from a DHCP server `dhcp4: yes`.

To assign a static IP address to `ens3` interface, edit the file as follows:

- Set DHCP to `dhcp4: no`.
- Specify the static IP address. Under `addresses:` you can add one or more IPv4 or IPv6 IP addresses that will be assigned to the network interface.
- Specify the gateway.
- Under `nameservers`, set the IP addresses of the nameservers.

/etc/netplan/00-installer-config.yaml

```yaml
network:
  version: 2
  ethernets:
    enp0s3:
      dhcp4: no 
      addresses:
        - 10.0.50.29/24
      routes:
        - to: default
          via: 10.0.50.1
      nameservers:
          addresses: [8.8.8.8, 114.114.114.114]
helm install nfs-subdir-external-provisioner nfs-subdir-external-provisioner/nfs-subdir-external-provisioner \
    --set nfs.server=192.168.2.161 \
    --set nfs.path=/nfs/
```

When editing Yaml files, make sure you follow the YAML code indent standards. If the syntax is not correct, the changes will not be applied.

Once done, save the file and apply the changes by running the following command:

```shell
netplan apply
```

