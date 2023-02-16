SSH сервер с контролем доступа по заранее заготовленному списку пользователей и их SSH-ключей позволяющий индивидуально проксировать их запросы на просмотр логов, рестарт контейнеров либо запуск интерактивной консоли внутри контейнера (exec -it sh/bash/etc). 

Перед запуском скопируйте `config.yaml.example` в `config.yaml` и настройте его под себя.

Для работы необходимо наличие в `PATH` `docker` и (для windows/git bash - `winpty`).

### Примеры:

Получение списка доступных пользователю `us` контейнеров и команд для них:

```
$ ssh us@127.0.0.1 -p 23229 ps

Non-accessible containers:
- php74
- centrifugo
- php8
- nginx
- sphinx
- postgres
- mysql
- elasticsearch
- redis
- rabbitmq
- kibana

Accessible containers:
- php81 restart|logs|sh|bash
```

Подключение внутрь контейнера:
```
$ ssh us@127.0.0.1 -p 23229 exec php81 bash
```

Отслеживание логов:
```
$ ssh us@127.0.0.1 -p 23229 logs php81 
```

Рестарт контейнера:
```
$ ssh us@127.0.0.1 -p 23229 logs php81 
```

Список доступных команд:
```
$ ssh us@127.0.0.1 -p 23229

Usage:
   [command]

Available Commands:
  exec        Runs 'docker exec -it [container_name] [command=sh]'
  help        Help about any command
  logs        Runs 'docker logs [container] --follow'
  ps          List available containers
  restart     Runs 'docker restart [container]'

Flags:
  -h, --help   help for this command

Use " [command] --help" for more information about a command.
```

