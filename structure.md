1.  🚪 **API Gateway:**
    *   Принимает запросы.
    *   *Новая фича:* Валидирует JWT-токен (чтобы не пускать неавторизованных к платежам).
    *   Перенаправляет запросы.
2.  🔐 **Auth Service (PostgreSQL):**
    *   `POST /register` -> создает юзера в БД.
    *   `POST /login` -> возвращает JWT.
    *   *(Опционально)*: при регистрации кидает в Kafka событие `user.registered`.
3.  💰 **Billing Service (PostgreSQL):**
    *   Слушает `user.registered` и создает кошелек с нулевым балансом для нового юзера.
    *   Слушает `payment.init`, списывает баланс (если хватает), кидает `payment.result`.
4.  💳 **Payment Service (PostgreSQL):**
    *   `POST /pay` -> создает платеж `pending`, кидает `payment.init`.
    *   Слушает `payment.result` -> обновляет статус на `success/failed`.
5.  ✉️ **Notification Service (без БД):**
    *   Слушает `payment.result` и логирует чек.
ан? Если согласен, давай начнем с **Этапа 2 (Auth Service)**. Создай папку `auth` (или `services/auth`), и я скину тебе начальный код для поднятия HTTP сервера и структуры User!