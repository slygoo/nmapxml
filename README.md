A small NMAP xml paresr that grabs all the open TCP ports of a host and transforms them into CSV useful for post analysis in excel. 

Can be used on one xml file with the -f flag or on the current working directory with no flags.

Example:

```
nmap 192.168.254.186 -oA scan

PORT     STATE SERVICE

53/tcp   open  domain

80/tcp   open  http

88/tcp   open  kerberos-sec

135/tcp  open  msrpc

139/tcp  open  netbios-ssn

389/tcp  open  ldap

443/tcp  open  https

445/tcp  open  microsoft-ds

464/tcp  open  kpasswd5

593/tcp  open  http-rpc-epmap

636/tcp  open  ldapssl

3268/tcp open  globalcatLDAP

3269/tcp open  globalcatLDAPssl

3389/tcp open  ms-wbt-server

9091/tcp open  xmltec-xmlmail=

./nmapxml

2025/01/08 00:59:13 [+] Successfully Outputed To NmapResults.csv

cat NmapResults.csv 

IP,Ports

192.168.254.186,53 80 88 135 139 389 443 445 464 593 636 3268 3269 3389 9091
```

