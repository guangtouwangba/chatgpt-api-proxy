#bin/bash

migrate -database "$(cat config/database.yml | python -c 'import yaml,sys;print( yaml.safe_load(sys.stdin)["database"]'))"