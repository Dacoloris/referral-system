Тестовое задание для Stakewolle
--

---
Запуск:
-- 
1. Поднимаем docker с Postgres
2. Накатываем миграции
3. Запускаем сервис
```
docker compose up -d --build
go run migrations/auto.go
go run cmd/main.go
```

---
Хэндлеры:
--
- POST
```
/login (логин)
/register (регистрация)
/referral-code (создать реферальный код)
```
- GET
```
/referral-code/{email} (получить реферальный код по email юзера)
/referrals/{user_id"} (получить всех реферов юзера)
```
- DELETE
```
/referral-code (удалить реферальный код)
```

---
Тестовое задание реализовано не в полном мере т.к. слишком ресурсоемко по времени \
Для простоты запуска оставил .env в репозитории