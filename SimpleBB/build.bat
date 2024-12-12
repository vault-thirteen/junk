:: This script builds the server.
@ECHO OFF

SET build_dir=_build_
SET cert_dir=cert
SET sql_dir=sql
SET sql_init_dir=table_init
SET exe_dir=cmd
SET config_dir=config
SET config_file_ext=json
SET acm_folder=ACM
SET gw_folder=GWM
SET mm_folder=MM
SET nm_folder=NM
SET sm_folder=SM
SET jwt_folder=JWT
SET captcha_folder=RCS
SET captcha_images_folder=rcs_img
SET smtp_folder=SMTP
SET tool_folder=tool
SET argon2_tool_folder=Argon2
SET jwt_tool_folder=MakeJWToken
SET assets_folder=assets
SET frontend_assets_folder=frontend

:: Create the folders.
MKDIR "%build_dir%"
MKDIR "%build_dir%\%tool_folder%"
MKDIR "%build_dir%\%cert_dir%"
MKDIR "%build_dir%\%sql_dir%"
MKDIR "%build_dir%\%assets_folder%"
MKDIR "%build_dir%\%assets_folder%\%frontend_assets_folder%"

:: Show version of Go language.
go version

:: 1. ACM Module.
ECHO 1. ACM Module

:: Build the ACM module (service).
CD "%exe_dir%\%acm_folder%"
go build
IF %Errorlevel% NEQ 0 EXIT /b %Errorlevel%
MOVE "%acm_folder%.exe" ".\..\..\%build_dir%\"
CD ".\..\..\"

:: Copy related files for the ACM module (service).
COPY "%config_dir%\%acm_folder%.%config_file_ext%" "%build_dir%\"
MKDIR "%build_dir%\%cert_dir%\%acm_folder%"
COPY "%cert_dir%\%acm_folder%" "%build_dir%\%cert_dir%\%acm_folder%\"
MKDIR "%build_dir%\%sql_dir%\%acm_folder%"
MKDIR "%build_dir%\%sql_dir%\%acm_folder%\%sql_init_dir%"
COPY "%sql_dir%\%acm_folder%\%sql_init_dir%" "%build_dir%\%sql_dir%\%acm_folder%\%sql_init_dir%\"
MKDIR "%build_dir%\%cert_dir%\%jwt_folder%"
COPY "%cert_dir%\%jwt_folder%" "%build_dir%\%cert_dir%\%jwt_folder%\"

:: 2. Gateway Module.
ECHO 2. Gateway Module

:: Build the Gateway module (service).
CD "%exe_dir%\%gw_folder%"
go build
IF %Errorlevel% NEQ 0 EXIT /b %Errorlevel%
MOVE "%gw_folder%.exe" ".\..\..\%build_dir%\"
CD ".\..\..\"

:: Copy related files for the Gateway module (service).
COPY "%config_dir%\%gw_folder%.%config_file_ext%" "%build_dir%\"
MKDIR "%build_dir%\%cert_dir%\%gw_folder%"
COPY "%cert_dir%\%gw_folder%" "%build_dir%\%cert_dir%\%gw_folder%\"
MKDIR "%build_dir%\%sql_dir%\%gw_folder%"
MKDIR "%build_dir%\%sql_dir%\%gw_folder%\%sql_init_dir%"
COPY "%sql_dir%\%gw_folder%\%sql_init_dir%" "%build_dir%\%sql_dir%\%gw_folder%\%sql_init_dir%\"
COPY "%assets_folder%\%frontend_assets_folder%" "%build_dir%\%assets_folder%\%frontend_assets_folder%\"

:: 3. Message Module.
ECHO 3. Message Module

:: Build the Message module (service).
CD "%exe_dir%\%mm_folder%"
go build
IF %Errorlevel% NEQ 0 EXIT /b %Errorlevel%
MOVE "%mm_folder%.exe" ".\..\..\%build_dir%\"
CD ".\..\..\"

:: Copy related files for the Message module (service).
COPY "%config_dir%\%mm_folder%.%config_file_ext%" "%build_dir%\"
MKDIR "%build_dir%\%cert_dir%\%mm_folder%"
COPY "%cert_dir%\%mm_folder%" "%build_dir%\%cert_dir%\%mm_folder%\"
MKDIR "%build_dir%\%sql_dir%\%mm_folder%\%sql_init_dir%"
COPY "%sql_dir%\%mm_folder%\%sql_init_dir%" "%build_dir%\%sql_dir%\%mm_folder%\%sql_init_dir%\"

:: 4. Notification Module.
ECHO 4. Notification Module

:: Build the Notification module (service).
CD "%exe_dir%\%nm_folder%"
go build
IF %Errorlevel% NEQ 0 EXIT /b %Errorlevel%
MOVE "%nm_folder%.exe" ".\..\..\%build_dir%\"
CD ".\..\..\"

:: Copy related files for the Notification module (service).
COPY "%config_dir%\%nm_folder%.%config_file_ext%" "%build_dir%\"
MKDIR "%build_dir%\%cert_dir%\%nm_folder%"
COPY "%cert_dir%\%nm_folder%" "%build_dir%\%cert_dir%\%nm_folder%\"
MKDIR "%build_dir%\%sql_dir%\%nm_folder%\%sql_init_dir%"
COPY "%sql_dir%\%nm_folder%\%sql_init_dir%" "%build_dir%\%sql_dir%\%nm_folder%\%sql_init_dir%\"

:: 5. Subscription Module.
ECHO 5. Subscription Module

:: Build the Subscription module (service).
CD "%exe_dir%\%sm_folder%"
go build
IF %Errorlevel% NEQ 0 EXIT /b %Errorlevel%
MOVE "%sm_folder%.exe" ".\..\..\%build_dir%\"
CD ".\..\..\"

:: Copy related files for the Subscription module (service).
COPY "%config_dir%\%sm_folder%.%config_file_ext%" "%build_dir%\"
MKDIR "%build_dir%\%cert_dir%\%sm_folder%"
COPY "%cert_dir%\%sm_folder%" "%build_dir%\%cert_dir%\%sm_folder%\"
MKDIR "%build_dir%\%sql_dir%\%sm_folder%\%sql_init_dir%"
COPY "%sql_dir%\%sm_folder%\%sql_init_dir%" "%build_dir%\%sql_dir%\%sm_folder%\%sql_init_dir%\"

:: 6. Captcha Module.
ECHO 6. Captcha Module

:: Build the Captcha module (service).
CD "%exe_dir%\%captcha_folder%"
go build
IF %Errorlevel% NEQ 0 EXIT /b %Errorlevel%
MOVE "%captcha_folder%.exe" ".\..\..\%build_dir%\"
CD ".\..\..\"

:: Copy related files for the Captcha module (service).
COPY "%config_dir%\%captcha_folder%.%config_file_ext%" "%build_dir%\"
MKDIR "%build_dir%\%captcha_images_folder%"

:: 7. SMTP Module.
ECHO 7. SMTP Module

:: Build the SMTP module (service).
CD "%exe_dir%\%smtp_folder%"
go build
IF %Errorlevel% NEQ 0 EXIT /b %Errorlevel%
MOVE "%smtp_folder%.exe" ".\..\..\%build_dir%\"
CD ".\..\..\"

:: Copy related files for the SMTP module (service).
COPY "%config_dir%\%smtp_folder%.%config_file_ext%" "%build_dir%\"

:: 8. Auxiliary tools.
ECHO 8. Auxiliary tools

:: 8.1. Argon tool.
ECHO 8.1. Argon tool
CD "%tool_folder%\%argon2_tool_folder%"
go build
IF %Errorlevel% NEQ 0 EXIT /b %Errorlevel%
MOVE "%argon2_tool_folder%.exe" ".\..\..\%build_dir%\%tool_folder%\"
CD ".\..\..\"

:: 8.2. JWT tool.
ECHO 8.2. JWT tool
CD "%tool_folder%\%jwt_tool_folder%"
go build
IF %Errorlevel% NEQ 0 EXIT /b %Errorlevel%
MOVE "%jwt_tool_folder%.exe" ".\..\..\%build_dir%\%tool_folder%\"
CD ".\..\..\"

ECHO SUCCESSFUL BUILD
