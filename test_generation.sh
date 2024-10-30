# Need a linux environment to run this script
if [ -d "tests" ]; then
    rm -r "tests"
fi

mkdir "tests"

# Tests made for projects.json, lab.json and code.json

# Random directory name (openssl rand -base64 12)
dirName=$(openssl rand -base64 12)

# For projects.json
mkdir -p "tests/projects/$dirName/src" && mkdir -p "tests/projects/$dirName/lib"

touch "tests/projects/$dirName/src/main.c" && touch "tests/projects/$dirName/lib/lib.h"