# IPVS, iptables and kube-proxy

This is an overview of the underlying technologies that drives load balancing. It covers LVS, Netfilter, iptables, IPVS and eventually kube-proxy.

### LVS (Linux Virtual Server)

One of the ways to implement software load balancing is via LVS (Linux Virtual Server), as [previously discussed](https://www.digihunch.com/2020/01/several-ways-to-ensure-high-availability/). The diagram below shows the LVS [framework](http://www.linuxvirtualserver.org/about.html), with IPVS as the fundamental technology:

![](/Users/jack/Documents/learn/k8s/Net/images/lvs.jpeg)

The major work of the LVS project is to develop advanced IP load balancing software (IPVS), application-level load balancing software (KTCPVS), cluster management components. [KTCPVS ](http://www.linuxvirtualserver.org/software/ktcpvs/ktcpvs.html)implements application-level load balancing inside the Linux kernel (still under development). [IPVS ](http://www.linuxvirtualserver.org/software/ipvs.html)is an advanced IP load balancing software implemented inside the Linux kernel. The IPVS code was already included into the standard Linux kernel 2.4 and 2.6.







https://www.digihunch.com/2020/11/ipvs-iptables-and-kube-proxy/