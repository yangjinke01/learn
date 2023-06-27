# Introduction to OpenSSL tool

SSL certificates are now in high demand. Since Google’s “HTTPS Everywhere” campaign, the encryption landscape is changing considerably. They initially gave digital certificates an SEO boost as an inducement to install them, and then Chrome made HTTPS nearly necessary for everyone. Popular browsers like Firefox and Chrome will label that website as Not Secure if you don’t use an SSL certificate.SSL deployment is critical to the success and security of a website. And, because so many website owners are learning about SSL for the first time, it’s critical to provide them with all of the necessary tools and services. OpenSSL is one such utility. In this topic, we are going to learn about the OpenSSL tool.

**What is OpenSSL?**

Open SSL is a general-purpose cryptography package that implements the TLS protocol in an open-source manner. It is available for Windows, Linux, macOS, as well as BSD computers, and was first released in 1998. Users can use OpenSSL to execute a variety of SSL-related operations, such as generating CSRs and private keys, as well as installing SSL certificates.

**What is the purpose of OpenSSL?**

User applies for a digital certificate and installs SSL files on the server using OpenSSL Generate the Certificate Signing Request. You can also convert your certificate to different SSL formats and do additional verifications.

## How to use OpenSSL?

It’s all about the command lines in OpenSSL. We’ve included a list of typical OpenSSL commands for individual users below.

Make sure you’re using the most recent version of OpenSSL.

Knowing what version of OpenSSL you have is critical since it dictates the cryptographic algorithms and protocols you may use.

the most recent OpenSSL release was 1.1.1. It’s the first version to provide support for TLS 1.3. 1.0.2 and 1.1.0 are two previous releases that are still supported.

Run the following command to see what version of OpenSSL you have:

```shell
openssl version -a
```

**CSR Generation:**

Users can generate their own CSR code with OpenSSL. A CSR is a block of encoded text that contains information about the website and business. Users should submit the CSR for approval to the Certificate Authority. A private key is required for the certificate request, from which the public key is generated. While you can utilise a current private key, it’s best to produce a fresh one every time users create a CSR.

## How to Generate private keys separately?

Users must specify the key algorithm, key size, and an alternate passphrase to produce your private key. RSA is The typical key algorithm, however, ECDSA can be used in some cases. Make sure users won’t have any compatibility difficulties while selecting a key algorithm.

When utilising the RSA key algorithm, choose 2048 bits for key size, and 256 bits for the ECDSA algorithm. Any key size less than 2048 is insecure, while a greater value may cause performance to suffer.

Then, users must determine whether or not a passphrase is required for the private key.  some servers will refuse to accept private keys with passwords.

Run the commands below whenever ready to produce a private key (using the RSA algorithm):

```shell
openssl genrsa -out domain.key 2048
```

The domain.key file will be created in your current directory using this command. The PEM format will be used to store your private key.

Run the following command to decode the private key:

```shell
openssl rsa -text -in domain.key -noout
```

## How to Extract public key?

To extract the public key from the private key, the following command is used:

```shell
openssl rsa -in domain.key -pubout -out domain_public.key
```

**Generate a Certificate Signing Request:**

It’s time to build CSR after you’re successfully generating the private key. It will be in PEM format and will contain information about the business and the public key generated from the private key. To generate a CSR, use the following command:

```shell
# Generate a Certificate Signing Request
openssl req -new -key domain.key -out domain.csr

# x509 Generate a self-signed certificate
openssl req -x509 -new -nodes -key myCA.key -sha256 -days 1825 -out myCA.pem
```

*user will be asked a few questions by OpenSSL. Consider the following situations:*

* Country Name: Enter country’s two-letter code. Make sure the country user submit his the organization’s official residence if the user has a Business Validation or Extended Validation certificate.
* Name of State/Province: enter the complete name of the state in which the user’s business is registered.
* Name of Locality: Enter the name of the city or town where the company is located.
* Organization Name: enter your company’s official registered name. For Domain Validation certificates, for example, users can use NA
* Organization Unit Name: it’s commonly Web Administration.
* Common Name: Enter the Fully Qualified Domain Name (FQDN) to which your SSL certificate will be assigned. Consider the domain educba.com. Add an asterisk in front of the domain name (e.g. *.educba.com) to activate a wildcard certificate.
* Email Address: Give a valid email address.
* A challenging password: It is an out-of-date characteristic that the Certificate Authorities no longer require. If there is any confusion, leave this box blank.

```text
Country Name (2 letter code) [XX]:CN
State or Province Name (full name) []:JiangSu
Locality Name (eg, city) [Default City]:NanJing
Organization Name (eg, company) [Default Company Ltd]:ecloud tech
Organizational Unit Name (eg, section) []:IT Dept
Common Name (eg, your name or your server's hostname) []:*.ecloudtech.com
Email Address []:yangjk@ecloudtech.com
```

**Verify the certificate’s information:**

After CA sends an SSL certificate, execute the command below to make that the certificate’s information matches the private key.

```shell
openssl x509 -text -in domain.crt –noout
```
