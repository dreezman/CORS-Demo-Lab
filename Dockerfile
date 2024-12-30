
FROM openresty/openresty:1.25.3.1-2-alpine-fat
EXPOSE 80
EXPOSE 443
RUN apk update && apk add --no-cache \
    openrc \
    curl \
      bash \
      net-tools \
      tcpdump \
      vim

LABEL "org.opencontainers.image.vendor"="Layer8"
LABEL "org.opencontainers.image.title"="Browser Security Training"
LABEL "org.opencontainers.image.authors"="info@layer8.cloud"
LABEL "org.opencontainers.image.source"="TBD"


########## Build the nginx.conf with CSP controls ##########

COPY nginx/conf/mime.types /etc/nginx/mime.types
COPY nginx/entrypoint/ /docker-entrypoint.d/
COPY nginx/entrypoint/docker-entrypoint.sh /docker-entrypoint.sh
COPY nginx/conf/nginx.conf /usr/local/openresty/nginx/conf/
RUN  rm -f entrypoint/docker-entrypoint.sh && \
     chmod +x /docker-entrypoint.sh && \
     chmod +x /docker-entrypoint.d/*.sh
COPY nginx/pubkey/* /usr/local/openresty/nginx/
COPY nginx/local/csp-policy.conf    /usr/share/nginx-config/csp-policy.conf
# clear it out in case it was left over from a previous run
RUN  echo 'echo "" > /usr/share/nginx-config/csp-policy.conf' > /docker-entrypoint.d/10-clear-csp-policy.sh && \
    chmod a+x /docker-entrypoint.d/10-clear-csp-policy.sh
RUN chmod -R a+rwx  /usr/share/nginx-config/csp-policy.conf
# reload ngx every 5 seconds to pick up changes to the csp policy
RUN echo "while  sleep 5; do /usr/local/openresty/nginx/sbin/nginx -s reload &> /tmp/ngxreload   ; done &" > /docker-entrypoint.d/50-start-cron.sh && \
    chmod a+x /docker-entrypoint.d/50-start-cron.sh

# copy front end static files to the web root
COPY fe/ /usr/local/openresty/nginx/html/

CMD ["/docker-entrypoint.sh", "nginx"]
