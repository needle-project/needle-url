# NeedleUrl

The current application offers a local "url shortener" application written in GoLang and using Redis as a storage.

## 1. Installation and usage
The packages is prepared to be used either using Docker or as a stand-alone application.

### 1.1. Install using Docker

To run without any orchestration, you can run from the project root:
```
docker-compose build
docker-compose up
```

To build images, you can run from the specific path (docker/http or docker/redis):
```
docker build .
```

### 1.2. Install as a standalone
To use the application as a standalone, you can use one of the builds from `/build/` directory.

## 2. Usage

The application needs a `config.json` file that should be located in the same path as the binary.

```
{
  "port": 80,
  "redis_hostname": "redis",
  "redis_port": 6379,
  "redis_db": 0,
  "default_redirect_path": "https://www.example.com/",
  "admin_path": "/bin/needle-url/admin/",
  "basic_auth": {
    "username": "test",
    "password": "new"
  }
}
```
Description:
| Parameter | Description
| ------------- |:-------------:| -----:|
| port | The port on which the server will bind. Usually is 80 |
| redis_hostname | Redis connection hostname/IP |
| redis_port | Redis connection PORT |
| redis_db | Redis database used for storage |
| default_redirect_path | Where the client will be redirected if a short url does not exists or an error occurs. |
| admin_path | Fully qualified path to an HTML that will be used as an admin interface |
| basic_auth | Credentials for BasicAuth to application |

## 3. API

For documentation of the API, see docs/swagger/api.yml


## 4. Develop - Setup environment

- Install Go from [here](https://golang.org/)
- Clone this project
- Run `initial_provision.sh`
- Modify code
- Run `containers_up.sh`
- Run `smoke_tests.sh` to validate that all flows are working as it should