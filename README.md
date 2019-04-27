## Service update docker containers.
The service keeps track of container updates in a register. 
And based on the config updates the containers on your workstation. 
The service has been tested in the Ubuntu 18.04 system.

### Dependencies:
[go-bin-deb](https://github.com/mh-cbon/go-bin-deb)

### Installation:
From the project directory, run the command:
```
sh build.sh
```
After a successful build, the **deb** package will be created in the **dist** directory.

### Service API
The service starts on port *23001* and receives 
[notification messages from the registry](https://docs.docker.com/registry/notifications/) at the following address:
```
127.0.0.1:23001/api/v1/buid
```

### Setup registry example:
```
docker run -d -t \
--name registry \
--network host \
-p 5000:5000 \
-e REGISTRY_NOTIFICATIONS_ENDPOINTS="
- name: builder
  url: https://builder.your.site/api/v1/build
  timeout: 5s
  threshold: 5
  backoff: 30s
  " \
registry:2
```

### Server Tuning:
On the machine on which the service will be launched you need to create a configuration file:
```
/srv/auto-builder/config.json
```
Add rights for this file to the user **auto-builder** and group **auto-builder**:
```
sudo chmod 0755 /srv/auto-builder/config.json
sudo chown -R auto-builder:auto-builder /srv/auto-builder/
```

An example of filling the configuration file:
```
{
	"auth": {
		"login": "your_registry_login",
		"password": "your_registry_password"
	},
	"config_list": [
		{
			"image": "your.registry.addr/your-container-1",
		},
		{
			"image": "your.registry.addr/your-container-2",
			"environments": [
				"PRIVKEY_PEM=/keys/priv1.pem",
				"ADDRESS=127.0.0.1",
				"PORT=8872",
				"WEBROOT=/webroot/",
				"TMP=/tmp/",
				"NATS_URL=nats://127.0.0.1:4222",
				"MONGO_HOST=127.0.0.1:27017",
				"GIN_MODE=release"
			],
			"cmd": [
				"pwd"
			],
			"volumes": [
				"VolumeTMP001:/tmp/",
				"VolumeTMP002:/tmp/"
			]
		}
	]
}
```
