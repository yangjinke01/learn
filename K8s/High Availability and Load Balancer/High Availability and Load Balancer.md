# High Availability and Load Balancer

### Overview

Fault tolerance and high availability are two architectural characteristics that people often confuse with each other. High availability focuses on minimizing downtime. It guarantees uptime, but not performance in the event of component failures. Fault tolerance, on the other hand, focuses on stable capacity event in the event of component failures. Fault tolerance is stricter and therefore more expensive. High availability can be achieved either by clustering, or load balancing. A cluster involves several nodes, all able to perform the same function, but may take different roles at different times (e.g. primary, standby) in order for the cluster to perform its function as a single system. In Linux, clustering is implemented by pacemaker or corosync. With a high load system, it is common to set up load balancing system to achieve high availability (and fault tolerance).

### Load balancing

The idea of load balancing is simple: load goes high and we want to scale horizontally instead of simply upgrading server hardware. At a high level, there has been three approaches to load balancing:

- **DNS rotating:** (aka. DNS round robin) DNS record resolves to multiple IPs, very simple and cheap to implement. Since DNS is cached, the load distribution will come imbalanced and it’s hard to re-balance, making this a very limited approach;

- **Hardware Load Balancer** using dedicated hardware device to configure load balancing. This option is expensive and only enterprises can afford it

- **Software Load Balancer** using software to achieve load balancing. These solutions are affordable, and usually open-source. They can be loaded on commodity hardware (including NIC). Some (e.g.Nginx) refers to themselves as software-based ADC. Major players are:
  - HA Proxy
  - Nginx
  - Linux Virtual Server (LVS, L4 only)

The hardware ADCs are usually supported commercially and there are plenty of resources from their white papers. There is an ongoing debate about whether one is better than the other. However, there is no doubt that a software-based load balancer is more approachable as open-source tools. The line between software and hardware load balancers becomes blurred today as hardware vendors try to adapt their software appliance to commodity hardware. Check out [this](https://www.nginx.com/blog/not-all-software-load-balancers-are-created-equal/) article. The rest of this post, will focus on software-based load balancer.

### Software-based load balancer

We explained that ADC (application delivery controller) is an expanded set of features from load balancer, and will only cover the load balancer part of the feature set in this article.

[HAProxy](https://www.haproxy.org/) supports both layer 4 and layer 7 load balancing. It supports load balancing based on cookie and session, as well as health check. Since it is layer 4 load balancing, it supports any TCP protocol such as read traffic for MySQL.  

[Nginx](https://www.nginx.com/) is a high-performance, event-driven, cross-platform layer 7 load balancing application. It works as a reverse proxy where it receives request for the Internet and forwards it to (upstream) internal servers. It consumes less memory than many of its alternatives for layer 7 load balancing. There are many strategies for load balancing such as round robin, by weight, by hash of requesting IP, by upstream response time, or by URL hash. It supports 20-30 k concurrent connections, and support compression and health check. It is known to be very stable and common for small and medium volume. Nginx has a commercial counterpart Nginx Plus with advanced features.

Nginx and HA proxy are commonly used in front end load balancing. For backend traffic such as database (e.g. separating read write traffic), LVS can be used.

### Linux Virtual Server

[LVS](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/4/html/virtual_server_administration/ch-lvs-overview-vsa) (Linux Virtual Server) is part of standard Linux kernel. It performs layer 4 load balancing based on TCP or UDP and therefore consumes less memory and CPU. Compared to layer 7 load balancing, the performance is generally higher, and the configuration is less complex (with simpler routing rules). [LVS](http://www.linuxvirtualserver.org/) is usually configured in a [common cluster architecture](http://www.linuxvirtualserver.org/architecture.html) involving these components:

- Load balancer: the front-end machine of the whole cluster systems, and balances requests from clients among a set of servers, so that the clients consider that all the services is from a single IP address.
- Server cluster: set of servers running actual business workload
- Shared storage: a shared storage space for the servers, such as NFS

![](/Users/jack/Documents/learn/k8s/High Availability and Load Balancer/images/ysG0q.png)

Load balancer is the single entry-point of server cluster systems, it can run IPVS that implements IP load balancing techniques inside the Linux kernel, or KTCPVS that implements application-level load balancing inside the Linux kernel. When IPVS is used, all the servers are required to provide the same services and contents, the load balancer forward a new client request to a server according to the specified scheduling algorithms and the load of each server. No matter which server is selected, the client should get the same result. When KTCPVS is used, servers can have different contents, the load balancer can forward a request to a different server according to the content of request. Since KTCPVS is implemented inside the Linux kernel, the overhead of relaying data is minimal, so that it can still have high throughput.

IPVS is also called layer-4 switching, it directs TCP/UDP requests to the real servers behind load balancer. It works in three modes:

- Network Address Translation (NAT)
- Direct Routing (DR)
- Tunnel mode (TUN)

These are three packet-forwarding methods in IPVS. The IPVS is implemented as a module over the netfilter framework, similar to [iptables](https://www.digihunch.com/2018/10/redhat-firewall-configuration-firewalld-vs-iptables/), which is also built on top of netfilter, based on chain and rules.

### Summary

We had an overview of high availability, and then expanded on load balancing, an important mechanism to implement high availability. We touched on both hardware-based and software-based load balancing technologies, and dived a little more into Linux Virtual Server. It is worth-noting that LVS is also the foundation of kube-proxy, the load balancing mechanism used in Kubernetes.

https://www.digihunch.com/2020/01/several-ways-to-ensure-high-availability/