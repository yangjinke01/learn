# DEPLOY A KUBERNETES CLUSTER USING ANSIBLE

In this article we will take a look at how to deploy a [Kubernetes](https://kubernetes.io/) cluster on CentOS 7.9 using [Ansible](https://www.ansible.com/) Playbooks. I have found Ansible to be a fantastic tool for getting a Kubernetes cluster up and running quickly in my development environment, and now use the Ansible playbooks detailed in this article when I need to stand up a Kubernetes cluster quickly and easily.

For the purposes of this article, we will use Ansible to deploy a small Kubernetes cluster – with one master node, used to manage the cluster, and two worker nodes, which will be used to run our container applications. To achieve this, we will use four Ansible playbooks. These will do the following:

- Create a new User Account for use with Kubernetes on each node
- Install Kubernetes and containerd on each node
- Configure the Master node
- Join the Worker nodes to the new cluster

If you are considering using Ansible to deploy Kubernetes already, I will assume you’re already somewhat familiar with both technologies. So, with that said, let’s get straight into the detail.

### Before we Deploy Kubernetes using Ansible

Before we can get started, we need a few prerequisites to be in place. This is what we are going to need:

- A host with Ansible installed. I’ve written previously about [how to install Ansible](https://buildvirtual.net/installing-and-configuring-ansible-on-centos/) – also, check out the online documentation! You should also set up an SSH key pair, which will be used to authenticate to the Kubernetes nodes without using a password, allowing Ansible to do it’s thing.
- Three servers/hosts to which we will use as our targets to deploy Kubernetes. I am using CentOS 7.9 and my servers each have 8GB ram and 4vCPUs. This is fine for my lab purposes, which I use to try out new things using Kubernetes. You need to be able to SSH into each of these nodes as root using the SSH key pair I mentioned above.

With that lot all in place we should be ready to go!

### Setting up Ansible to Deploy Kubernetes

Before we start to look at the Ansible Playbooks, we need to set up Ansible to communicate with the Kubernetes nodes. First of all, on our Ansible host, lets set up a new directory from which we we run our playbooks.

```shell
$ mkdir kubernetes
$ cd kubernetes
```

With that done, we now need to create a hosts file, to tell Ansible how to communicate with the Kubernetes master and worker nodes.

```shell
$ vim hosts
```

The content of the hosts file should look something like the following:

```ini
[masters]
master ansible_host=master ansible_user=root

[workers]
worker1 ansible_host=worker1 ansible_user=root
worker2 ansible_host=worker2 ansible_user=root
```

Listing the master node and the worker nodes in different sections in the hosts file will allow us to target the playbooks at the specfic node type later on.

Finally, with that done, we can test it’s working by doing a Ansible ping:

```shell
$ ansible -i hosts all -m ping
master | SUCCESS => {
    "ansible_facts": {
        "discovered_interpreter_python": "/usr/bin/python"
    },
    "changed": false,
    "ping": "pong"
}
worker1 | SUCCESS => {
    "ansible_facts": {
        "discovered_interpreter_python": "/usr/bin/python"
    },
    "changed": false,
    "ping": "pong"
}
worker2 | SUCCESS => {
    "ansible_facts": {
        "discovered_interpreter_python": "/usr/bin/python"
    },
    "changed": false,
    "ping": "pong"
}
```

All good! Lets move onto the first playbook.

### Creating a Kubernetes user with Ansible Playbook

Our first task in setting up the Kubernetes cluster is to create a new user on each node. This will be a non-root user, that has sudo privileges. It’s a good idea not to use the root account for day to day operations, of course. We can use Ansible to set the account up on all two nodes, quickly and easily. First, create a file in the working directory:

```shell
$ vim users.yml
```

Then add the following to the playbook:

```yaml
- hosts: 'workers, masters'
  become: yes

  tasks:
    - name: create the kube user account
      user: name=kube append=yes state=present createhome=yes shell=/bin/bash

    - name: allow 'kube' to use sudo without needing a password
      lineinfile:
        dest: /etc/sudoers
        line: 'kube ALL=(ALL) NOPASSWD: ALL'
        validate: 'visudo -cf %s'

    - name: set up authorized keys for the kube user
      authorized_key: 
        user: kube
        key: "{{item}}"
      with_file:
        - ~/.ssh/id_rsa.pub
```

We’re now ready to run our first playbook. To do so:

```shell
$ ansible-playbook -i hosts users.yml
```
Once done you should see:
```
PLAY RECAP *****************************************************************************************************
master                     : ok=4    changed=3    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0
worker1                    : ok=4    changed=3    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0
worker2                    : ok=4    changed=3    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0
```

### Install Kubernetes with Ansible Playbook

Now we’re getting to the fun part! With our user now created, we can move on to installing Kubernetes. Lets dive straight in and have a look at the playbook, which I have named install-k8s.yml:

```yaml
---
- hosts: "masters, workers"
  remote_user: kube
  become: yes
  become_method: sudo
  become_user: root
  gather_facts: yes
  connection: ssh

  tasks:
    - name: Disable swap
      shell: |
        swapoff -a
        sed -i '/ swap / s/^\(.*\)$/# \1/g' /etc/fstab

    - name: Add Docker YUM repository
      shell: |
        yum install -y yum-utils
        yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo

    - name: Install Docker
      yum:
        name: docker-ce
        state: present

    - name: Copy Docker daemon.json to masters and workers
      ansible.builtin.copy:
        src: ./daemon.json
        dest: /etc/docker/
        owner: root
        group: root
        mode: '0644'

    - name: Start Docker
      shell: |
        systemctl daemon-reload
        systemctl restart docker

    - name: Add K8s YUM repository
      ansible.builtin.copy:
        src: ./kubernetes.repo
        dest: /etc/yum.repos.d/
        owner: root
        group: root
        mode: '0644'

    - name: install K8s
      yum:
        name:
          - kubeadm
          - kubelet
          - kubectl
```

Run the playbook using the following syntax:

```shell
$ ansible-playbook -i hosts install-k8s.yml
```

There’s quite a lot going on here, so this one will take a little while to run whist the necessary packages are installed on each node. Once done you should see:

```
PLAY RECAP ******************************************************************************************************
master                     : ok=8    changed=5    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0   
worker1                    : ok=8    changed=5    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0   
worker2                    : ok=8    changed=5    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0 
```

