#!/bin/bash
#
# chmod +x ~/deploy-demo_todo/*.sh ~/deploy-demo_todo/goose.exe 
# ~/deploy-demo_todo/deploy.sh

# copy and organize files in destination directory on devdep server
rm -rf ~/apps/demo_todo
mkdir -p ~/apps/demo_todo
mv ~/deploy-demo_todo/{.,}* ~/apps/demo_todo
#
# load the docker image on app server
cd ~/apps/demo_todo
docker load -i ~/apps/demo_todo/demo_todo.tar
docker image ls | grep demo_todo
#