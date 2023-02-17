openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 30 -nodes
move key.pem 1\key.pem
move cert.pem 1\cert.pem
