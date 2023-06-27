# Resize a logical volume

Have you ever wondered how to extend your root or home directory file system partition using LVM? You might have low storage space and you need to increase the capacity of your partitions. This article looks at how to extend storage in Linux using Logical Volume Manager (LVM).

## Process summary

The process is straightforward. Attach the new storage to the system. Next, create a new Physical Volume (PV) from that storage. Add the PV to the Volume Group (VG) and then extend the Logical Volume (LV).

```shell
lsblk -f
pvcreate /dev/sdb

vgs
vgdisplay
vgextend centos /dev/sdb

lvs
lvdisplay
# lvextend -L +7G /dev/centos/root
lvextend -l +100%FREE /dev/centos/root
xfs_growfs /dev/centos/root
```

