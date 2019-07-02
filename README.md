# ngxpkg [![Build Status](https://travis-ci.org/webpkg/ngxpkg.svg?branch=master)](https://travis-ci.org/webpkg/ngxpkg)

A cli tool for nginx

## Install

Install go before use ngxpkg, or you can download [release](https://github.com/webpkg/ngxpkg/releases) version.
How to install go please look at https://golang.org/doc/install

```bash
go get github.com/webpkg/ngxpkg

mv $GOPATH/bin/ngxpkg $GOPATH/bin/ngc
```

or you can build it for other system by use

```bash
go get github.com/webpkg/ngxpkg

cd $GOPATH/src/github.com/webpkg/ngxpkg/

make
```
## Default Env

* NGX_CONFIG = 〜／.config／ngxpkg
* NGX_DIRECTORY_URL = https://acme-v01.api.letsencrypt.org/directory
* NGX_ALLOW_RENEW_DAYS = 30
* NGX_RESOURCE = https://ngxpkg.com/rc/
* NGX_SITE_CONFIG = /etc/nginx/conf.d
* NGX_SITE_ROOT = /opt/local/www

## Create new sites with domains ngxpkg.com, dbpkg.com ...

```bash
ngc new ngxpkg.com dbpkg.com
```

it will create conf files:

* /etc/nginx/conf.d/ngxpkg.com.conf
* /etc/nginx/conf.d/dbpkg.com.conf

and create sites:

* /opt/local/www/ngxpkg.com/
    - /opt/local/www/ngxpkg.com/public
    - /opt/local/www/ngxpkg.com/conf/fullchain.ecdsa.pem
    - /opt/local/www/ngxpkg.com/conf/privkey.ecdsa.pem
    - /opt/local/www/ngxpkg.com/conf/fullchain.rsa.pem
    - /opt/local/www/ngxpkg.com/conf/privkey.rsa.pem
    - /opt/local/www/ngxpkg.com/conf/ocsp.pem

* /opt/local/www/dbpkg.com/
    - /opt/local/www/dbpkg.com/public
    - /opt/local/www/dbpkg.com/conf/fullchain.ecdsa.pem
    - /opt/local/www/dbpkg.com/conf/privkey.ecdsa.pem
    - /opt/local/www/dbpkg.com/conf/fullchain.rsa.pem
    - /opt/local/www/dbpkg.com/conf/privkey.rsa.pem
    - /opt/local/www/dbpkg.com/conf/ocsp.pem

* /opt/local/www/conf
    - /opt/local/www/conf/dhparam.pem
    - /opt/local/www/conf/ticket.pem

## Notice.

please replace the files with yours after create new site.

* /opt/local/www/conf/dhparam.pem
* /opt/local/www/conf/ticket.pem

## Renew ssl with domains ngxpkg.com, dbpkg.com ...

```bash
ngc renew ngxpkg.com dbpkg.com
```

it will renew all ssl certificates list on sites conf file when it's valid days less then NGX_ALLOW_RENEW_DAYS

## Reference Links

https://github.com/google/acme

https://gethttpsforfree.com/

https://mozilla.github.io/server-side-tls/ssl-config-generator/

## License

Copyright (C) 2017-2018 ZhiQiang Huang, All Rights Reserved.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.