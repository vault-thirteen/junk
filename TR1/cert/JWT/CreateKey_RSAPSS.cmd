openssl genpkey -algorithm rsa-pss -pkeyopt rsa_keygen_bits:8192 -pkeyopt rsa_pss_keygen_md:sha512 -pkeyopt rsa_pss_keygen_mgf1_md:sha512 -pkeyopt rsa_pss_keygen_saltlen:1024 -out jwtPrivateKey.pem

openssl pkey -pubout -in jwtPrivateKey.pem -out jwtPublicKey.pem
