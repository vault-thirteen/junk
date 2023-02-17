#!/bin/bash

# Exit on Error.
set -e

# Input file.
FILE_TO_PROCESS="$1"
if [ ! -f "$FILE_TO_PROCESS" ]
then
  echo "Input file '$FILE_TO_PROCESS' does not exist."
  exit 1
fi

# Output Folder.
OUTPUT_FOLDER="$2"
if [ ! -d "$OUTPUT_FOLDER" ]
then
  echo "Output folder '$OUTPUT_FOLDER' does not exist."
  exit 1
fi

# Convert the file.
soffice --convert-to pdf "$FILE_TO_PROCESS" --outdir "$OUTPUT_FOLDER" --headless
