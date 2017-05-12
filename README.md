# ngxpkg

A cli tool for nginx

## working on it, please don't use it at production.


Create a new site with domain ngxpkg.com:

```bash
ngx new ngxpkg.com
```

it will create a conf file at: /etc/nginx/conf.d/ngxpkg.com.conf

and create site at: /opt/local/www/ngxpkg.com/

can change it by set NGX_SITE_CONFIG and NGX_SITE_ROOT