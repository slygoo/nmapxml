A small NMAP xml paresr that grabs all the open TCP ports of a host and transforms them into CSV useful for post analysis in excel. 

Can be used on one xml file with the -f flag or on the current working directory with no flags.

Example:

nmap 192.168.254.186 -oA scan

Starting Nmap 7.94 ( https://nmap.org ) at 2025-01-08 00:58 AEDT

Nmap scan report for dc01.sly.local (192.168.254.186)

Host is up (0.00020s latency).

Not shown: 985 closed tcp ports (reset)

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
MAC Address: 00:0C:29:64:89:D7 (VMware)

Nmap done: 1 IP address (1 host up) scanned in 10.70 seconds

./nmapxml

2025/01/08 00:59:13 [+] Successfully Outputed To NmapResults.csv

cat NmapResults.csv 

IP,Ports

192.168.254.186,53 80 88 135 139 389 443 445 464 593 636 3268 3269 3389 9091

