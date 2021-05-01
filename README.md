# Системный мониторинг

[![Go Report Card](https://goreportcard.com/badge/github.com/Haba1234/sysmon)](https://goreportcard.com/report/github.com/Haba1234/sysmon)
[![codecov](https://codecov.io/gh/Haba1234/sysmon/branch/master/graph/badge.svg)](https://codecov.io/gh/Haba1234/sysmon)
![workflow](https://github.com/Haba1234/sysmon/actions/workflows/tests.yml/badge.svg)

Демон - программа, собирающая информацию о системе, на которой запущена, и отправляющая её своим клиентам по gRPC.

## Описание
Сбор метрик системы и отправка подписанным клиентам gRPC.

При подписке клиент в запросе указывает параметры:
- N: получение данных каждые N секунд;
- M: получение усредненных данных за последние M секунд.

### Доступные метрики
- Средняя загрузка системы (load average)

### Поддерживаемые ОС
- linux

### Конфигурация демона
- Через командную строку можно задать порт сервера gRPC.
- Через файл конфигурации можно задать какие метрики будет собирать демон.