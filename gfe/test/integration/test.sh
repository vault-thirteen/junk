#!/bin/bash

# Данный Shell скрипт производит интеграционное тестирование микро-сервиса.
#
# Порядок действий:
#   - поднятие серверов, нужных для проведения тестирования, и их настройка;
#   - запуск тестов;
#   - остановка серверов.
#
# Требования:
#   - Установленный Docker и Docker-Compose;
#   - Установленный язык Go.

# Выходим из скрипта в случае ошибки.
set -e

# 1. Поднятие и настройка серверов.

# 1.1. PostgreSQL.

