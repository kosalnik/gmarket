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
account ||--o{ withdraw : "account.id = withdraw.account_id"
user ||--o{ accrual : "user.id = accrual.user_id"
account ||--o{ accrual : "account.id=accrual.account_id"
accrual ||--|o withdraw : "accrual.withdraw_id=withdraw.id"
payment ||--|o withdraw : "payment.withdraw_id=withdraw.id"
user ||--o{ withdraw : "user.id = withdraw.user_id"


user {
  uuid id PK
  string login UK
  string password "Хэш пароля"
}

account {
  uuid id PK
  uuid user_id FK,UK "Хозяин счёта"
  decimal balance "Текущая сумма баллов на счету"
}

accrual {
  uuid id PK
  uuid user_id FK
  uuid account_id FK
  uuid withdraw_id FK
  string order_number UK
  accrual_status status "NEW,PROCESSING,PROCESSED,INVALID"
  timestamp created_at
  timestamp status_changed_at
}

payment {
  uuid id PK
  uuid user_id FK
  uuid account_id FK
  uuid withdraw_id FK
  string order_number 
  payment_status status "NEW,PROCESSING,PROCESSED,INVALID"
  timestamp created_at
  timestamp status_changed_at
}

withdraw {
  uuid id PK
  uuid user_id FK ""
  uuid account_id FK ""
  withdraw_action action "DEBIT,CREDIT"
  decimal amount "Сумма операции"
  withdraw_status status "NEW,PROCESSING,PROCESSED,INVALID"
  timestamp processed_at "Время окончания обработки"
}
```

Все денежные операции сосредоточены в таблицах счёта: `account` и переводов: `withdraw`. Денежные единицы хранятся в 
колонках с типом `decimal(20,2)`.

Таблицы `payment` и `accrual` хранят контекст переводов. Платежи и пополнения.

Информация о пользователе, хранится в таблице `user`. Все остальные таблицы ссылаются на `user` как на корень аггрегата,
дабы в будущем иметь удобный ключ шардирования, если вдруг всдумается раскидывать схему на шарды.


### Взаимодействие в системе

Взаимодействия клиента с компонентами системы можно видеть на следующей диаграме.

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
      order ->> db: Создать запись в таблице payment
      order ->> account: Создать операцию списания
    deactivate order
    activate account
      account ->> db: BEGIN transaction
      account ->> db: Перевести withdraw в статус PROCESSING<br/> с условием что он имел статус REGISTERED
      alt changed 0 rows?
        account ->> db: ROLLBACK
        account ->> *: fail
      else иначе
          account ->> db: COMMIT
      end
      account ->> db: BEGIN
      account ->> db: SELECT FOR UPDATE amount FROM account ...
      alt if amount < X
          account ->> db: Перевести withdraw в статус INVALID
          account ->> db: COMMIT
          account ->> *: invalid
      else иначе
          account ->> db: Уменьшить счёт ID на сумму X
          account ->> db: Перевести withdraw в статус PROCESSED<br/>WHERE status=PROCESSING
          alt changed 0 rows?
            account ->> db: ROLLBACK
            account ->> *: fail
          else иначе
            account ->> db: COMMIT
            account ->> *: Success
          end
      end
    deactivate account
```

```mermaid
sequenceDiagram
    title Пополнение
    * ->> account : "Пополнить на X счёт ID"
    account ->> db: UPDATE SET amount=amount+X WHERE ID
    account ->> *: "Success"
```

### Абстрактная схема взаимодействия с системой
```mermaid
sequenceDiagram
   actor client
   client ->> api: "регистрация пользователя"
   client ->> api: "регистрация заказа"
   activate api
      api ->> cashback: сверка с системой<br/>расчёта баллов
   deactivate api
   activate cashback
   alt есть баллы?
    cashback ->> account : начислить
   end
   deactivate cashback
   client ->> api: списание баллов в счёт<br/>оплаты заказа
   activate api
      api ->> account: списать средства
      activate account
      alt есть баллы?
      account ->> account: списать со счёта
      account ->> withdraw: оставить запись в истории
      end
      deactivate account
   deactivate api
```

### Схема взаимодействия с системой начисления баллов

```mermaid
sequenceDiagram
    actor client
    client ->> api: "POST /api/user/orders"
    activate api
        api ->> order: загрузка заказа<br/>в систему
        activate order
            order ->> order: Создать запись в статусе NEW
            order ->> cashback: Поставить в очередь на получение баллов
            activate cashback
            cashback ->> queue: Создать PENDING запись в очереди
            cashback ->> order: .
            deactivate cashback
            order ->> api: .
        deactivate order
        api ->> client: .
   deactivate api
```

```mermaid
sequenceDiagram
    title Общение с внешней системой accrual на предмет наличия начислений по заказу
    activate cashback
    loop
        cashback ->> queue: Получаем PENDING задачу из очереди<br/>и переводим в статус PROCESSING
        activate queue
        queue ->> cashback: task
        deactivate queue
        cashback ->> order: получаем заказ
        activate order
            order ->> cashback: order
        deactivate order
        cashback ->> accrual: Послать запрос на получение сведений о начислении за заказ
         activate accrual
         accrual ->> cashback: ответ от accrual
         deactivate accrual
            alt status=200
                alt status=PROCESSED
                cashback ->> order: обновить статус заказа
                cashback ->> account: начислить баллы на счёт
                cashback ->> cashback: Убрать заказ из очереди
                else status=INVALID
                cashback ->> order: обновить статус заказа
                cashback ->> cashback: Убрать заказ из очереди
                else status=REGISTERED
                cashback ->> order: обновить статус заказа
                cashback ->> cashback: Отложить в конец очереди
                else status=PROCESSING
                cashback ->> cashback: Отложить в конец очереди
                end
            else status=204
                cashback ->> cashback: Отложить в конец очереди
            else status=429
                cashback ->> cashback: Отложить в конец очереди<br/>и поставить воркеры в паузу на указанное время
            else остальные статусы
               cashback ->> cashback: Ретраить
            end
    end
    deactivate cashback
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
