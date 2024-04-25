##
## Build nginx from scratch and install it
##

# get the environment variables from the nginx_build_env.sh
# to config for this environment
#

if [ -f "./nginx/nginx_build_env.sh" ]; then
    echo "Found config file exists."
    . ./nginx/nginx_build_env.sh 
else
    echo "Config file ./nginx/nginx_build_env.sh not found, must invoke from repo root. Exiting."
    exit 1
fi



echo "******** Install pre-req modules **************************"
cd ${ngx_dir}
service nginx stop
sudo apt-mark hold nginx ## prevent any upgrades that wipe out our config
sudo apt-get --yes   update && sudo apt-get --yes upgrade
sudo apt-get --yes   install gcc
sudo apt-get --yes   make
sudo apt install --yes   make-guile
sudo apt install --yes   libnginx-mod-rtmp
sudo apt-get install  --yes     zlib1g-dev zlib1g
sudo apt-get install   --yes   libpcre3 libpcre3-dev
sudo apt-get install   --yes   libssl-dev
sudo apt-get install  --yes    libperl-dev
sudo apt-get install  --yes    libgd-dev
sudo apt-get install  --yes    net-tools

echo "******** Install nginx nonce modules **************************"
cd ${ngx_dir}/ ; mkdir -pv ${ngx_tar_dir} ${ngx_dev_kit_tar_dir} ${setnginx_tar_dir}
wget -O  ${ngx_tar_dir}.tar.gz          http://nginx.org/download/${ngx_tar_dir}.tar.gz
tar --strip-components 1 -C ${ngx_tar_dir}         -xzvf ${ngx_tar_dir}.tar.gz
wget -O  ${ngx_dev_kit_tar_dir}.tar.gz  https://github.com/simpl/ngx_devel_kit/archive/${ngx_dev_kit_ver}.tar.gz
tar --strip-components 1 -C ${ngx_dev_kit_tar_dir} -xzvf ${ngx_dev_kit_tar_dir}.tar.gz
wget -O  ${setnginx_tar_dir}.tar.gz     https://github.com/openresty/set-misc-nginx-module/archive/${setnginx_ver}.tar.gz
tar --strip-components 1 -C ${setnginx_tar_dir}    -xzvf ${setnginx_tar_dir}.tar.gz

echo "******** configure nginx with nonces **************************"
echo  ${ngx_dir} "---------->" ${ngx_tar_dir} " :together: " ${ngx_dir}/${ngx_tar_dir} ${ngx_dir}/${ngx_dev_kit_tar_dir} ${ngx_dir}/${setnginx_tar_dir}

#
# Make sure these exist
#
sudo mkdir -vp /etc/nginx/sites-available /etc/nginx/modules-available
sudo mkdir -vp /etc/nginx/sites-enabled   /etc/nginx/modules-available

cd ${ngx_dir}/${ngx_tar_dir}
# https://www.photographerstechsupport.com/tutorials/hosting-wordpress-on-aws-tutorial-part-2-setting-up-aws-for-wordpress-with-rds-nginx-hhvm-php-ssmtp/#nginx-source
 ./configure --prefix=/opt/nginx \
    --prefix=/etc/nginx \
    --sbin-path=/usr/sbin/nginx \
    --conf-path=/etc/nginx/nginx.conf \
    --error-log-path=/var/log/nginx/error.log \
    --http-log-path=/var/log/nginx/access.log \
    --pid-path=/var/run/nginx.pid \
    --lock-path=/var/run/nginx.lock \
    --http-client-body-temp-path=/var/cache/nginx/client_temp \
    --http-proxy-temp-path=/var/cache/nginx/proxy_temp \
    --http-fastcgi-temp-path=/var/cache/nginx/fastcgi_temp \
    --http-uwsgi-temp-path=/var/cache/nginx/uwsgi_temp \
    --http-scgi-temp-path=/var/cache/nginx/scgi_temp \
    --user=nginx \
    --group=nginx \
    --with-http_ssl_module \
    --with-http_realip_module \
    --with-http_gunzip_module \
    --with-http_gzip_static_module \
    --with-threads \
    --with-file-aio \
    --with-http_v2_module \
    --with-cc-opt='-O2 -g -pipe -Wall -Wp,-D_FORTIFY_SOURCE=2 -fexceptions -fstack-protector --param=ssp-buffer-size=4 -m64 -mtune=native' \
    # --add-module=../headers-more-nginx-module \
    --with-http_ssl_module \
    --add-module=${ngx_dir}/${ngx_dev_kit_tar_dir}/ \
    --add-module=${ngx_dir}/${setnginx_tar_dir}/


echo "******** Build nginx **************************"
sudo  make -j2
sudo  make install
## use this to build docker image
# https://stackoverflow.com/questions/28863126/creating-a-docker-image-with-nginx-compile-options-for-optional-http-modules