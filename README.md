# RFC. Создание сервиса: Накопительная система лояльности «Гофермарт»

Данный проект носит учебный характер. Создаётся в целях
закрепления пройденного материала.

Нужно создать сервис накопительной системы лояльности по заданной [спецификации](./SPECIFICATION.md)

## Задача

Создать сервис который будет заниматься обслуживанием бонусных счетов пользователей.

Требования к сервису:

* Должна иметься регистрация новых пользователей
* Аутентификация зарегистрированных пользователей
* Взаимодействие с внешней системой расчета баллов за заказ
* Ведение счета зарегистрированного пользователя
  * Пополнение на величину выдаваемую внешним сервисом при регистрации заказа
  * Частичное или полное списание с бонусного счёта в оплату заказа
* Выдача величины текущего счета, статистики пополнений и расходов

## Предложеное решение

Создать сервис на языке Go.

Взаимодействие с клиентами осуществляется по REST API описанному в [спецификации](./SPECIFICATION.md)

В качестве хранилища данных использовать Postgres.

## Реализация



### Схема БД

```mermaid
erDiagram
user ||--o{ payment : "user.id=payment.user_id"
account ||--o{ payment : "account.id=payment.account_id"
user ||--|| account : "user.id = account.user_id"
account ||--o{ transaction : "account.id = transaction.account_id"
user ||--o{ accrual : "user.id = accrual.user_id"
account ||--o{ accrual : "account.id=accrual.account_id"
accrual ||--|o transaction : "accrual.transaction_id=transaction.id"
payment ||--|o transaction : "payment.transaction_id=transaction.id"
user ||--o{ transaction : "user.id = transaction.user_id"


user {
  uuid id PK
  string login UK
  string password "Хэш пароля"
}

account {
  uuid id PK
  uuid user_id FK,UK "Хозяин счёта"
  numeric balance "Текущая сумма баллов на счету"
}

accrual {
  uuid id PK
  uuid user_id FK
  uuid account_id FK
  uuid transaction_id FK
  string order_number UK
  accrual_status status "NEW,PROCESSING,PROCESSED,INVALID"
  timestamp created_at
  timestamp status_changed_at
}

payment {
  uuid id PK
  uuid user_id FK
  uuid account_id FK
  uuid transaction_id FK
  string order_number 
  payment_status status "NEW,PROCESSING,PROCESSED,INVALID"
  timestamp created_at
  timestamp status_changed_at
}

transaction {
  uuid id PK
  uuid user_id FK ""
  uuid account_id FK ""
  transaction_direction direction "DEPOSIT,WITHDRAW"
  numeric amount "Сумма операции"
  transaction_status status "NEW,PROCESSING,PROCESSED,REJECTED"
  timestamp processed_at "Время окончания обработки"
}
```

Все денежные операции сосредоточены в таблицах счёта: `account` и переводов: `transaction`. Денежные единицы хранятся в 
колонках с типом `decimal(20,2)`.

Таблицы `payment` и `accrual` хранят контекст переводов. Платежи и пополнения.

Информация о пользователе, хранится в таблице `user`. Все остальные таблицы ссылаются на `user` как на корень агрегата,
дабы в будущем иметь удобный ключ шардирования, если вдруг понадобится раскидывать схему на шарды.


### Взаимодействие в системе

Взаимодействия клиента с компонентами системы можно видеть на следующей диаграмме.

```mermaid
flowchart TD
    client([client])
    accrualExternal{{accrual external service}}
    subgraph Application
       auth[[auth]]
       user[[user]]
       order[[order]]
       account[[account]]
       accrual[[accrual]]
    end
    client --аутентификация\nавторизация--> auth
    auth --> user
    client --регистрация--> user
    client --приём номеров заказов\nучёт и ведение списка переданных номеров--> order
    client --учёт и ведение накопительного счёта--> account
    order --проверка принятых номеров заказов\nчерез систему расчёта баллов лояльности--> accrual
    order --списание--> account
    accrual --начисление вознаграждения--> account
    accrual --> accrualExternal
```


### Взаимодействие со счетом

Операции пополнения и снятия со счета похожи, и изображены на следующих диаграмах.

```mermaid
sequenceDiagram
    title Снятие со счёта
    actor *
    participant order
    participant account
    participant db
    * ->> order : Оплатить заказ со счета
    activate order
      critical transaction 
        order ->> db: Создать запись в таблице payment
        order ->> account: Создать операцию списания
        activate account
            account ->> db : Создать запись о списании в withdraw
        deactivate account
      end
    deactivate order
    activate account
      critical transaction
        account ->> db: Заблокировать счёт для списания
        account ->> db: Проверить достаточность средств для списания
        alt if amount < X
          account ->> db: Списать с account
          account ->> db: Перевести операцию списания в PROCESSED
        else иначе
          account ->> db: Перевести операцию списания в REJECTED
        end
        account ->> db: Разблокировать счёт
      end
    deactivate account
```

Операция начисления на счёт обрабатывается аналогично

### Схема взаимодействия с системой начисления баллов

```mermaid
sequenceDiagram
    actor client
    client ->> api: "POST /api/user/orders"
    activate api
        api ->> order: загрузка заказа<br/>в систему
        activate order
            order ->> order: Создать запись в статусе NEW
            order ->> accrual: Поставить в очередь на получение баллов
            activate accrual
            accrual ->> db: Создать PENDING запись в очереди
            accrual ->> order: .
            deactivate accrual
            order ->> api: .
        deactivate order
        api ->> client: .
   deactivate api
```

### Обработка ответа от accrual
```mermaid
flowchart
    response[/response/]
    response --> status200{httpStatus=200}
    status200 --Yes--> statusProcessed{status=PROCESSED}
    status200 --"No"--> status204{httpStatus=204}
    statusProcessed --Yes--> markProcessed[сменить статус заказа на PROCESSED<br/>начислить баллы]
    statusProcessed --"No"--> statusInvalid{status=INVALID}
    statusInvalid --Yes--> markInvalid[сменить статус заказа на INVALID]
    statusInvalid --"No"--> statusRegistered{status=REGISTERED}
    statusRegistered --"No"--> statusProcessing{status=PROCESSING}
    statusRegistered --Yes--> retry
    statusProcessing --"No"-->internalError
    statusProcessing --Yes--> retry
    status204 --Yes--> retry((O))
    status204 --"No"--> status429{httpStatus=429}
    status429 --"No"--> retry
status429 --Yes--> workerPause[Поставить воркеры в паузу на некоторое время]
    retry --"tries = tries - 1"--> tries{tries > 0}
    tries --"No"--> failed[Пометить FAILED<br/>Убрать из очереди]
    tries --Yes--> queue[В конец очереди]
```
