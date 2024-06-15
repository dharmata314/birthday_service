# Backend сервис для уведомлений о днях рождениях, написанный на Go
## Развертывание
### Общие настройки
Общие настройки приложения содержатся в [конфиге](https://github.com/dharmata314/birthday_service/blob/main/config/config.yaml). В зависимости от способа развертывания какие-либо параметры могут меняться. 
В конфиге содержатся основные данные, необходимые для работы приложения.

### Docker 
Для развертывания в Docker Compose создан файл [docker-compose.yml](https://github.com/dharmata314/birthday_service/blob/main/docker-compose.yml)
Необходимо запустить команду
```
docker-compose up --build app
```
### Нативно
Для нативного запуска достаточно запустить приложение из папки [cmd](https://github.com/dharmata314/birthday_service/tree/main/cmd). 
Предварительно, необходимо установить зависимости из [go.mod](https://github.com/dharmata314/birthday_service/blob/main/go.mod) и изменить в [конфиге](https://github.com/dharmata314/birthday_service/blob/main/config/config.yaml) ```host: postgres``` на ```host: localhost```
## Общее
Приложение представляет из себя сервис уведомлений о днях рождениях. 

Сервис устроен следующим образом:

Пользователи региструется в сервисе, авторизуются, затем добавляются сотрудники с именами и датами рождений, после чего создаются подписки пользователей на уведомления о днях рождениях сотрудников.

Доступны следующие операции:
 - Регистрация пользователя
 - Авторизация пользователя
 - Удаление пользователя
 - Изменение данных пользователя
 - Добавление сотрудника
 - Получение списка всех сотрудников
 - Удаление сотрудника
 - Добавление подписки на уведомление о дне рождении
 - Удаление подписки на уведомление о дне рождении

Для доступа к большинству функционала (кроме регистрации и авторизации) необходим доступ по токену.
Токен выдается пользователю после авторизации.
В дальнейшем токен должен передаваться вместе с заголовком запроса:
```
Authorization: Bearer <token>
```
Уведомления о днях рождениях присылаются на электронную почту, которая указывается при регистрации. Чтобы функция отправки писем работала, необходимо в [main файле](https://github.com/dharmata314/birthday_service/blob/main/cmd/main.go) указать данные конфигурации SMTP профиля для Вашей почты. Пример:
```
cfgSMTP := &config.ConfigSMTP{
		SMTPHost:     "smtp.yandex.ru",
		SMTPPort:     587,
		SMTPUsername: "test@yandex.ru",
		SMTPPassword: "mzvsllelcirlsfpr",
	}
 ```
Письмо может оказаться в папке спама.
Письма присылаются раз в минуту. Если нужно изменить этот параметр, то необходимо поменять константу notificationFrequency в main файле на нужное количество минут. 
## Примеры запросов

Запросы при нативном запуске делаются без команды ```docker-compose exec app```

Регистрация пользователя:
```
 docker-compose exec app curl -X POST \
    -H "Content-Type: application/json" \
    -d '{"email": "test@email.com", "password": "testPassword"}' \
    http://localhost:8080/users/new
```
Авторизация пользователя:
```
docker-compose exec app curl -X POST \
    -H "Content-Type: application/json" \
    -d '{"email": "test@email.com", "password": "testPassword"}' \
    http://localhost:8080/login
```
Удаление пользователя:
```
docker-compose exec app curl -X DELETE \
-H "Authorization: Bearer <token>" \
http://localhost:8080/users/{id}
```
Изменение данных пользователя:
```
docker-compose exec curl -X PATCH \
-H "Authorization: Bearer <token>" \
-H "Content-Type: application/json" \
-d '{"email": "newEmail@email.com", "password": "NewPassword", "id": 1}' \
http://localhost:8080/users/{id}
```
Добавление сотрудника:
```
docker-compose exec app curl -X POST \
-H "Authorization: Bearer <token>" \
-H "Content-Type: application/json" \
-d '{"name": "John", "birthday": "14.06.1995"}' \
http://localhost:8080/emp
```
Получение списка всех сотрудников
```
docker-compose exec curl -X GET \
-H "Authorization: Bearer <token>" \
http://localhost:8080/employees
```
Добавление подписки на уведомление о дне рождении:
```
docker-compose exec app curl -X POST \
-H "Authorization: Bearer <token>" \
-H "Content-Type: application/json" \
-d '{"emp_id": 1, "user_id": 1}' \
http://localhost:8080/subs
```
Удаление подписки на уведомление о дне рождении:
```
docker-compose exec curl -X DELETE \
-H "Authorization: Bearer <token>" \
http://localhost:8080/subs/{id}
```
