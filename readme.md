# Cloudflare GeoBlock plugin for traefik

Cloudflare GeoBlock is a plugin to block traffic based on header **_Cf-Ipcountry_** that issued by Cloudflare.

## Usage

You need to make sure that traefik is live behind Cloudfare, otherwise the plugin will return HTTP Code 403.

This plugin also will automatically copy the IP from **_Cf-Connecting-Ip_** header, so you can chain it with traefik IpWhitelist middleware (https://doc.traefik.io/traefik/middlewares/http/ipwhitelist/)

### Configuration

Static Configuration:

```yaml
pilot:
  token: xxxxx

experimental:
  plugins:
    example:
      moduleName: github.com/moduit-engineering/cloudflare-geoblock
      version: v0.1.4
```

Dynamic Configuration:

```yaml
# Dynamic configuration

http:
  routers:
    my-router:
      rule: host(`demo.localhost`)
      service: service-foo
      entryPoints:
        - web
      middlewares:
        - my-plugin
 
  middlewares:
    my-plugin:
      plugin:
        cfgeoblock:
          whitelistCountry: ["ID", "SG"]
          disabled: false #optional
```
