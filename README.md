Build in deb package:
```
sh build.sh
```

After a successful build, the package will be created in the dist directory.

Registry must send hooks to the address:
```
your.server.com/api/vi/build
```

Setup registry example:
```
docker run -d -t \
--name $NODE_NAME \
--network host \
-p 5000:5000 \
--restart=always \
-v "$(pwd)"/auth:/auth \
-e REGISTRY_AUTH=htpasswd \
-e REGISTRY_AUTH_HTPASSWD_REALM="Registry Realm" \
-e REGISTRY_AUTH_HTPASSWD_PATH=/auth/htpasswd \
-e REGISTRY_NOTIFICATIONS_ENDPOINTS="
- name: builder
  url: https://builder.your.site/api/v1/build
  timeout: 5s
  threshold: 5
  backoff: 30s
  " \
registry:2

```

Port used by default
- 23001 - HTTP API

Create a file in the directory:
```
/srv/auto-builder/config.json
```

Add the following data to it:
```
{
	"auth": { // Autenticate block
		"login": "your_login", // Your registry login
		"password": "your_password" // Your registry password
	},
	"config_list": [ // Containers configs
		{
			"image": "registry.host/repository", // Image name
			"environments": [
				"YOUR_ENV_KEY=YOUR_ENV_VALUE", // Startup container variables
				"YOUR_ENV_KEY=YOUR_ENV_VALUE"  // Startup container variables
			]
		},
		{
			"image": "registry.host/repository", // Image name
			"environments": [
				"YOUR_ENV_KEY=YOUR_ENV_VALUE", // Startup container variables
				"YOUR_ENV_KEY=YOUR_ENV_VALUE"  // Startup container variables
			]
		}
	]
}
```

Add rights to this config:
```
sudo chown -R auto-builder /srv/auto-builder
sudo chmod 755 /srv/auto-builder/config.json
```
