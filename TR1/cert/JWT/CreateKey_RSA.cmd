openssl genrsa -out jwtPrivateKey.pem 8192
openssl rsa -pubout -in jwtPrivateKey.pem -out jwtPublicKey.pem
