# AGI-приложение для взаимодействия с Такси-Мастер

Приложение представляет собой часть целого комплекса приложенией.
Только часть функций можно использовать отдельно от остальных.

## Встроенное IVR

```
same => n,AGI(agi://127.0.0.1:4580/ivr)
```

## Получение данных звонящего
### Для описания IVR в диалплане

```
same => n,AGI(agi://127.0.0.1:4580/caller_info)
; выставит переменные:
; CATEGORYID
; CLIENT_GROUP_ID
; IS_DRIVER
; ID_ZAKAZ
; ORDER_ID
; IS_PRIOR
; ORDER_STATE
; ORDER_CONFIRM
; IN_PLACE
; IN_CAR
; CREW_ID
; CREW_GROUP_ID
; DRIVER_PHONE
; DRIVER_TIMECOUNT
; SOUND_COLOR
; SOUND_COLOR
; SOUND_MARK
; SOUND_MARK
; GOSNUMBER
; CREATION_WAY
```

## Получение данных о произвольном номере

```
same => n,AGI(agi://127.0.0.1:4580/phone_info?phone=${CALLERID(number)})
```

## Проговорить одноразовый код 
### Зависит от другого приложения

```
same => n,AGI(agi://127.0.0.1:4580/auth?phone=${CALLERID(number)})
```

## Отправить сообщение всем диспетчерам

```
same => n,AGI(agi://127.0.0.1:4580/show_tm_message?header=HEADER&text=TEXT&timeout=20)
same => n,AGI(agi://127.0.0.1:4580/show_tm_message?header=${URIENCODE(Очередь operators)}&text=${URIENCODE(Очередь звонков: ${QUEUE_WAITING_COUNT(operators)})}&timeout=20)
```

## Отмена заказа

```
same => n,AGI(agi://127.0.0.1:4580/cancel_by_order_id?order_id=${URIENCODE(${ID_ZAKAZ})})
```

## Выставить группу клиента

```
same => n,AGI(agi://127.0.0.1:4580/set_client_group?phone=79876543210&group=34)
```

## Получить номер волителя по позывному

```
same => n,AGI(agi://127.0.0.1:4580/driver_phone_by_crew?crew_code=${EXTEN})
same => n,NoOp(${DRIVER_PHONE})
```

## Получить номер водителя по номеру клиента из текущего заказа

```
same => n,AGI(agi://127.0.0.1:4580/driver_phone_by_caller?phone=${EXTEN})
same => n,NoOp(${DRIVER_PHONE})
```

## Прикрепить запись разговора

```
same => n,Set(NOW=${EPOCH})
same => n,Set(BACK_MONITOR_FILENAME=${REPLACE(MONITOR_FILENAME,/,\\)}.wav)
; type 0 — Исходящий, 1 — Входящий
same => n,AGI(agi://127.0.0.1:4580/record?type=0&length=10&order_id=${URIENCODE(${ID_ZAKAZ})}&date=${URIENCODE(${STRFTIME(${NOW},,%d%m%Y%H%M%S)})}&phone=${URIENCODE(${CALLERID(number)})}&call_id=${UNIQUEID}&path=${URIENCODE(${BACK_MONITOR_FILENAME})})
```
