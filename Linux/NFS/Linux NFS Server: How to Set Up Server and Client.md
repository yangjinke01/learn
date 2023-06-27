# Linux NFS Server: How to Set Up Server and Client

## What is Linux NFS Server?

Network File Sharing (NFS) is a protocol that allows you to share directories and files with other Linux clients over a network. Shared directories are typically created on a file server, running the NFS server component. Users add files to them, which are then shared with other users who have access to the folder.

An NFS file share is mounted on a client machine, making it available just like folders the user created locally. NFS is particularly useful when disk space is limited and you need to exchange public data between client computers.

## Quick Tutorial : Setting Up an NFS Server with an NFS Share

Let’s see how to set up an NFS server and create an NFS file share, which client machines can mount and access.

### Installing NFS Server

Here is how to install the NFS Kernel—this is the server component that enables a machine to expose directories as NFS shares.

```shell
yum -y install nfs-utils
```

### Create Root NFS Directory

We’ll now create the root directory of the NFS shares, this is also known as an export folder.

```shell
sudo mkdir /mnt/myshareddir
```

Set permissions so that any user on the client machine can access the folder (in the real world you need to consider if the folder needs more restrictive settings).

```shell
sudo chown nobody:nobody /mnt/myshareddir #no-one is owner

sudo chmod 777 /mnt/myshareddir #everyone can modify files
```

### Define Access for NFS Clients in Export File

To grant access to NFS clients, we’ll need to define an export file. The file is typically located at /etc/exports

Edit the /etc/exports file in a text editor, and add one of the following three directives.

All the directives below use the options rw, which enables both read and write, sync, which writes changes to disk before allowing users to access the modified file, and no_subtree_check, which means NFS doesn’t check if each subdirectory is accessible to the user.

| **To enable access to a single client**  | /mnt/myshareddir {clientIP}(rw,sync,no_subtree_check)        |
| ---------------------------------------- | ------------------------------------------------------------ |
| **To enable access to several clients**  | /mnt/myshareddir {clientIP-1}(rw,sync,no_subtree_check){clientIP-2}(...){clientIP-3}(...) |
| **To enable access to an entire subnet** | /mnt/myshareddir {subnetIP}/{subnetMask}(rw,sync,no_subtree_check) |

### Make the NFS Share Available to Clients

You can now make the shared directory available to clients using the exportfs command. After running this command, the NFS Kernel should be restarted.

```shell
sudo exportfs -a #making the file share available
```

If you have a firewall enabled, you’ll also need to open up firewall access using the sudo ufw allow command.

### Installing NFS Client Packages

```shell
sudo yum install nfs-utils
```

### Mounting the NFS File Share Temporarily

You can mount the NFS folder to a specific location on the local machine, known as a mount point, using the following commands.

Create a local directory—this will be the mount point for the NFS share. In our example we’ll call the folder /var/locally-mounted.

```shell
sudo mkdir /var/locally-mounted
```

Mount the file share by running the mount command, as follows. There is no output if the command is successful.

```shell
sudo mount -t nfs {IP of NFS server}:{folder path on server} /var/locally-mounted
```

For example:

```shell
sudo mount -t nfs 192.168.20.100:/mnt/myshareddir /var/locally-mounted
```

The mount point now becomes the root of the mounted file share, and under it you should find all the subdirectories stored in the NFS file share on the server.

To verify that the NFS share is mounted successfully, run the **mount** command or **df -h**.

### Mounting NFS File Shares Permanently

Remote NFS directories can be automatically mounted when the local system is started. You can define this in the /etc/fstab file. In order to ensure an NFS file share is mounted locally on startup, you need to add a line to this file with the relevant file share details.

To automatically mount NFS shares on Linux, do the following:

Create a local directory that will be used to mount the file share.

```shell
sudo mkdir /var/locally-mounted
```

Edit the /etc/fstab file using the nano command or any text editor.

Add a line defining the NFS share. Insert a tab character between each parameter. It should appear as one line with no line breaks.

The last three parameters indicate NFS options (which we set to default), dumping of file system and filesystem check (these are typically not used so we set them to 0).

```shell
{IP of NFS server}:{folder path on server} /var/locally-mounted nfs defaults 0 0
```

Now mount the file share using the following command. The next time the system starts, the folder will be mounted automatically.

```shell
mount /var/locally-mounted

mount {IP of NFS server}:{folder path on server}
```

<https://kuboard.cn/learning/k8s-intermediate/persistent/nfs.html#%E8%83%8C%E6%99%AF>