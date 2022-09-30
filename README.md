# Updtr

Drop-in replacement of [Watchtower](https://github.com/containrrr/watchtower) for manual updating.


https://user-images.githubusercontent.com/1058537/193275951-29149586-00b2-45cc-a6f2-64ea4424e1ab.mp4


## Quick Start

Run the updtr container with the following command:

```bash
docker run -d --name updtr -v /var/run/docker.sock:/var/run/docker.sock shemanaev/updtr
```

or with docker-compose:

```yaml
services:
  updtr:
    image: shemanaev/updtr
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - /etc/localtime:/etc/localtime:ro
    ports:
      - 8080:8080
    restart: unless-stopped
```

## Configuration

Updtr is built on top of Watchtower and can be configured via environment variables.

### Supported Env Variables

- [TZ](https://containrrr.dev/watchtower/arguments/#time_zone)

- [DOCKER_HOST](https://containrrr.dev/watchtower/arguments/#docker_host)
- [DOCKER_API_VERSION](https://containrrr.dev/watchtower/arguments/#docker_api_version)
- [DOCKER_TLS_VERIFY](https://containrrr.dev/watchtower/arguments/#tls_verification)

- [WATCHTOWER_DEBUG](https://containrrr.dev/watchtower/arguments/#debug)
- [WATCHTOWER_TRACE](https://containrrr.dev/watchtower/arguments/#trace)

- [WATCHTOWER_REMOVE_VOLUMES](https://containrrr.dev/watchtower/arguments/#remove_attached_volumes)
- [WATCHTOWER_INCLUDE_RESTARTING](https://containrrr.dev/watchtower/arguments/#include_restarting)
- [WATCHTOWER_INCLUDE_STOPPED](https://containrrr.dev/watchtower/arguments/#include_stopped)
- [WATCHTOWER_REVIVE_STOPPED](https://containrrr.dev/watchtower/arguments/#revive_stopped)
- [WATCHTOWER_NO_PULL](https://containrrr.dev/watchtower/arguments/#without_pulling_new_images)
- [WATCHTOWER_WARN_ON_HEAD_FAILURE](https://containrrr.dev/watchtower/arguments/#head_failure_warnings)

- [WATCHTOWER_POLL_INTERVAL](https://containrrr.dev/watchtower/arguments/#poll_interval)
- [WATCHTOWER_SCHEDULE](https://containrrr.dev/watchtower/arguments/#scheduling)
- [WATCHTOWER_LABEL_ENABLE](https://containrrr.dev/watchtower/arguments/#filter_by_enable_label)
- [WATCHTOWER_SCOPE](https://containrrr.dev/watchtower/arguments/#filter_by_scope)
- [WATCHTOWER_TIMEOUT](https://containrrr.dev/watchtower/arguments/#wait_until_timeout)

## Mappings

```yaml
# version of mappings file
version: 1
images:
  -
    # array of image names
    # can be with tag or without
    # tagger names has higher priority
    names: [shemanaev/updtr]
    # url to changelog
    url: https://github.com/shemanaev/updtr
    # possilble values:
    # Plaintext
    # Markdown
    # Asciidoc
    # Html
    # Github - notes of latest release, url must point to github repo
    type: Github
```

Additional mappings can be loaded from file:
`/data/mapping.yml`

```bash
docker run -d --name updtr -v ${PWD}/mapping.yml:/data/mapping.yml -v /var/run/docker.sock:/var/run/docker.sock shemanaev/updtr
```

It takes higher priority than built-in.
