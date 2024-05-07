# Проект

# ER диаграмма

```mermaid
erDiagram
user ||--|| account : "user.id = account.user_id"
user ||--o{ order : "user.id = order.user_id"
account ||--o{ withdraw : "account.id = withdraw.user_id"
order ||--o{ withdraw : "order.id = withdraw.order_id"
queue ||--|| order : "queue.order_id = order.id"

user {
  uuid id PK
  string login UK
  string password "Хэш пароля"
}

account {
  uuid id PK
  uuid user_id FK,UK "Хозяин счёта"
  float balance "Текущая сумма баллов на счету"
}

order {
   uuid id PK
   uuid user_id FK "Пользователь, загрузивший заказ"
   string num UK "Номер заказа"
   order_status status "NEW,PROCESSING,PROCESSED,INVALID"
   float accrual "Сумма начисления"
   timestamp uploaded_at "Время создания"
}

withdraw {
   uuid id PK
   uuid account_id FK "Счёт с которого списывают"
   uuid order_id FK "Заказ в счёт которого списываются средства"
%%   withdraw_direction withdraw_direction "DEBIT,CREDIT"
   float amount "Сумма списания"
   withdraw_status status "REGISTERED,PROCESSING,PROCESSED,INVALID"
   timestamp processed_at "Время окончания обработки"
}

queue {
    uuid order_id PK
    queue_status status "PENDING,PROCESSING,PROCESSED,INVALID"
    timestamp status_changed_at "Время последней смены статуса"
    int tries "Осталось попыток на обработку"
}
```

**Примечания к диаграмме**

1. Вместо `float` будет `decimal(20,2)`
2. queue_status, order_status и withdraw_status - type enum


# Взаимодействие в системе

```mermaid
flowchart
    client([client])
    client -- регистрация\nаутентификация\nавторизация --> auth 
    client -- приём номеров заказов\nучёт и ведение списка переданных номеров --> order
    client -- учёт и ведение накопительного счёта --> account
    order -- проверка принятых номеров заказов\nчерез систему расчёта баллов лояльности --> cashback
    cashback -- начисление за каждый подходящий номер заказа\n положенного вознаграждения на\n счёт лояльности пользователя --> account
    account -- списание --> withdraw
```

# Абстрактная схема взаимодействия с системой
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

# Схема взаимодействия с системой начисления баллов

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

## Обработка ответа от accrual
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

# go-musthave-diploma-tpl

Шаблон репозитория для индивидуального дипломного проекта курса «Go-разработчик»

# Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без
   префикса `https://`) для создания модуля

# Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m master template https://github.com/yandex-praktikum/go-musthave-diploma-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/master .github
```

Затем добавьте полученные изменения в свой репозиторий.
