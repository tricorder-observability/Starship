# TODO(jinyun): Move this file to ./Docker subdir, to make things more tidy.

FROM nginx

WORKDIR /usr/src/app/

COPY docker/run.sh .
RUN chmod +x run.sh
COPY ./docker/nginx.conf /etc/nginx/conf.d/default.conf

WORKDIR /usr/share/nginx/api-server-html

# These are for setting up HTTPS access for management UI
# This is part of the effort of
# https://github.com/tricorder-observability/starship/issues/418
COPY ./docker/test.key  /usr/share/nginx/test.key
COPY ./docker/test.ssl  /usr/share/nginx/test.ssl

COPY ./dist  /usr/share/nginx/api-server-html/

EXPOSE 80
# These are for setting up HTTPS access for management UI
# Does not work fully as expected yet.
EXPOSE 443

CMD /usr/src/app/run.sh && nginx -g 'daemon off;'
