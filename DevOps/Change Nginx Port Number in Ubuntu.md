# Change Nginx Port Number in Ubuntu

Sometimes, you may need to run NGINX web server on a different port, other than port 80. Here’s how to change NGINX port number in Ubuntu. You can use the same steps to change NGINX port number in Windows, Linux and Mac.

## How To Change NGINX Port Number in Ubuntu

Let’s see how to change NGINX port 80 to 8080. Before you change NGINX port it is advisable to check which port NGINX is running on.

After you change NGINX port number, you may want to use a [reporting software](http://ubiq.co/dashboard-reporting-software-tool) to monitor the key metrics about your website/application such as signups, traffic, sales, revenue, etc. using dashboards & charts, to ensure everything is working well.

 

## 1. Open NGINX configuration file

Open terminal and run the following command

```
# vi /etc/nginx/sites-enabled/default  [On Debian/Ubuntu]
# vi /etc/nginx/nginx.conf             [On CentOS/RHEL]
```

 

Bonus Read : [How to Rewrite URL Parameters in NGINX](http://ubiq.co/tech-blog/rewrite-url-parameters-nginx/)

 

## 2. Change NGINX port number

Look for the line that begins with [*listen*](http://nginx.org/en/docs/http/ngx_http_core_module.html#listen) inside *server* block. It will look something like

```
server {
        listen 80 default_server;
        listen [::]:80 default_server;
        ...
```

Change port number 80 to 8080 in above lines, to look like

```
server {
        listen 8080 default_server;
        listen [::]:8080 default_server;
        ...
```

 

Bonus Read : [How to Move NGINX Web Root to New Location](http://ubiq.co/tech-blog/how-to-move-nginx-web-root-to-new-location-on-ubuntu-18-04/)

 

## 3. Restart NGINX

Run the following command to check syntax of your updated config file.

```
$ sudo nginx -t
```

 

If there are no errors, run the following command to restart NGINX server.

```
$ sudo service nginx reload #debian/ubuntu
$ systemctl restart nginx #redhat/centos
```