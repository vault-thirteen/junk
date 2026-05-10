@ECHO OFF

SET /p USER_ID=Enter UserId:
SET /p SESSION_ID=Enter SessionId:
SET PRIVATE_KEY_FILE_PATH=cert\JWT\jwtPrivateKey.pem
SET PUBLIC_KEY_FILE_PATH=cert\JWT\jwtPublicKey.pem
SET SIGNING_METHOD=RS512

tool\MakeJWToken.exe -uid=%USER_ID% -sid=%SESSION_ID% -private_key=%PRIVATE_KEY_FILE_PATH% -public_key=%PUBLIC_KEY_FILE_PATH% -method=%SIGNING_METHOD%

PAUSE
