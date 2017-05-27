# ngxpkg [![Build Status](https://travis-ci.org/webpkg/ngxpkg.svg?branch=master)](https://travis-ci.org/webpkg/ngxpkg)

A cli tool for nginx

## Default Env

* NGX_CONFIG = 〜／.config／ngxpkg
* NGX_DIRECTORY_URL = https://acme-v01.api.letsencrypt.org/directory
* NGX_ALLOW_RENEW_DAYS = 30
* NGX_RESOURCE = https://rc.ngxpkg.com/
* NGX_SITE_CONFIG = /etc/nginx/conf.d
* NGX_SITE_ROOT = /opt/local/www

## Create new sites with domains test1.com, test2.com ...

```bash
ngx new test1.com test2.com
```

it will create conf files:

* /etc/nginx/conf.d/test1.com.conf
* /etc/nginx/conf.d/test2.com.conf

and create sites:

* /opt/local/www/test1.com/
* /opt/local/www/test2.com/

and apply ssl certificates for sites.

## Renew ssl with domains test1.com, test2.com ...

```bash
ngx renew test1.com test2.com
```

it will renew all ssl certificates list on sites conf file when it's valid days less then NGX_ALLOW_RENEW_DAYS

## Notice.

please replace the files with yours after create new site.

* dhparam.pem
* ticket.pem

## Reference Links

https://github.com/google/acme

https://gethttpsforfree.com/

https://mozilla.github.io/server-side-tls/ssl-config-generator/


## License

Copyright (C) 2017  ZhiQiang Huang (email: ngxpkg@gmail.com)

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