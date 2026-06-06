#!/bin/bash
#
# -- make port 9002 available for db and port 9003 for app in firewall
# chmod +x ~/apps/demo_todo/*.sh ~/apps/demo_todo/goose.exe 
# rm -rf ~/deploy-demo_todo
# ~/apps/demo_todo/run.sh

# run the mariadb container on app server
# (sql scripts in ./migrations are used for database migration)
mkdir -p ~/apps/db-demo_todo
docker run --name db-demo_todo -d --restart unless-stopped -v ~/apps/db-demo_todo:/var/lib/mysql:Z  -p 9002:3306 --env-file ~/apps/demo_todo/.env mariadb:12.2 
#
# run migrations on devdep server
cd ~/apps/demo_todo
./goose.exe up
#
# run the demo_todo container on app server
docker run --name demo_todo -d --restart unless-stopped -p 9003:80 --env-file ~/apps/demo_todo/.env demo_todo:latest
