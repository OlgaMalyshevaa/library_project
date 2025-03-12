# Library project

## Запуск проекта
1) Cклонируйте репозиторий
```bash
git clone https://github.com/username/library_project.git
```
2) Установите зависимости 
```bash
go mod tidy
```

3) Проверьте файл .env

```
DATABASE_URL=host=localhost port=5432 user=postgres dbname=postgres password=newpassword sslmode=disable
DB_NAME=music_db
PORT=8080
```
Замените `newpassword` на пароль, который вы используете для пользователя `postgres`

4) Запустите проект
```bash
go run cmd/apiserver/main.go
```

## Доступ к Swagger:
Swagger UI доступен по адресу:
```
http://localhost:8080/swagger/index.html
```
