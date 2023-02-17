#!/bin/bash

# Для использования этого скрипта, нужно установить утилиту 'jwt':
# go install github.com/golang-jwt/jwt/cmd/jwt@v3.2.1
# На момент 2021-07-07, актуальная версия утилиты = v3.2.1.

# Параметры скрипта:
#
# 1.  quiet=<yes|no>
#     Включает или выключает интерактивный режим работы.
#     Если задано 'quiet=yes', то скрипт не задаёт вопрос при перезаписи ключа
#     и не спрашивает кодовое слово для ключа.

# Выходим из скрипта в случае ошибки.
set -e

# Параметр 'quiet'.
IS_QUIET=false
if [ "$1" = "quiet=yes" ]
then
  IS_QUIET=true
fi
if [ "$IS_QUIET" = "true" ]
then
  echo "Quiet Mode is ON"
else
  echo "Quiet Mode is OFF"
fi

# Создаём папку, в которую будем складывать ключи (публичный и приватный).
mkdir -p key

# Пути до ключей и токена.
PRIVATE_KEY_PATH="./key/jwt_rsa_RS256.key"
PUBLIC_KEY_PATH="./key/jwt_rsa_RS256.key.pub"
JWT_TOKEN_PATH="./key/jwt_token.txt"

# Текущее время в формате Unix Timestamp.
CURRENT_TIME_UNIX=$(date +%s)
TOKEN_NBF_TIME_UNIX=$CURRENT_TIME_UNIX
TOKEN_EXP_TIME_UNIX=$((CURRENT_TIME_UNIX + 86400))
echo "TOKEN_NBF_TIME_UNIX=$TOKEN_NBF_TIME_UNIX, TOKEN_EXP_TIME_UNIX=$TOKEN_EXP_TIME_UNIX."

# Создаём RSA ключ.
# Приватный ключ создаём с помощью ssh-keygen, публичный -- с помощью openssl.
if [ "$IS_QUIET" = "true" ]
then
  # Этот код работает локально, но в Docker контейнере он не работает!
  #yes y | ssh-keygen -t rsa -b 4096 -m PEM -f $PRIVATE_KEY_PATH
  # Код для Docker контейнера.
  ssh-keygen -t rsa -b 4096 -m PEM -f $PRIVATE_KEY_PATH -q -N ""
else
  ssh-keygen -t rsa -b 4096 -m PEM -f $PRIVATE_KEY_PATH
fi
openssl rsa -in $PRIVATE_KEY_PATH -pubout -outform PEM -out $PUBLIC_KEY_PATH

# Создаём JSON веб токен.
# Время действия токена = с сей секунды до +1 дня (24 часов) от сей секунды.
echo {\"foo\":\"bar\"} | jwt -key $PRIVATE_KEY_PATH -alg RS256 -claim iss=tester -claim sub=test \
  -claim nbf="$TOKEN_NBF_TIME_UNIX" -claim exp="$TOKEN_EXP_TIME_UNIX" -sign - > $JWT_TOKEN_PATH

# Проверка подписи токена сложна и написать её в Shell скрипте не получится.
