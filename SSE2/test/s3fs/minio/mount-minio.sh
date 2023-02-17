URL="http://localhost:9002"
BUCKET_NAME="bucket-x"
FOLDER_TO_MOUNT="/media/username/minio"
DEBUG_LEVEL="warning" # "error"
PASSWORD_FILE="./.passwd-s3fs"

s3fs $BUCKET_NAME $FOLDER_TO_MOUNT -f \
	-o dbglevel="$DEBUG_LEVEL" \
	-o url="$URL" \
	-o passwd_file="$PASSWORD_FILE" \
	-o use_path_request_style
