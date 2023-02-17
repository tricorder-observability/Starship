#!/bin/bash -e

yarn install
yarn run build
sudo mkdir -p /usr/share/nginx/api-server-html
sudo cp -rf ./dist/* /usr/share/nginx/api-server-html/
sudo cp -rf ./docker/*.conf /usr/share/nginx/api-server-html/
sudo cp -rf ./docker/test.ssl share/nginx/api-server-html/test.ssl
sudo cp -rf ./docker/test.key share/nginx/api-server-html/test.key
sudo nginx -c /usr/share/nginx/api-server-html/nginx_test.conf