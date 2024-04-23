wd=`pwd`
ngx_ver="1.25.5"
ngx_tar_dir="nginx-${ngx_ver}"
setnginx_ver="v0.33"
setnginx_tar_dir="set-misc-nginx-module-${setnginx_ver}"
ngx_dev_kit_ver="v0.3.0"
ngx_dev_kit_tar_dir="ngx_dev_kit-${ngx_dev_kit_ver}"

cd ${wd}
mkdir -pv nginx
ngx_dir="${wd}/nginx"
cd ${ngx_dir}
sudo apt-get --yes --force-yes update && sudo apt-get upgrade
sudo apt-get --yes --force-yes install gcc
sudo apt-get --yes --force-yes make
sudo apt install --yes --force-yes make-guile
sudo apt install --yes --force-yes libnginx-mod-rtmp
sudo apt-get install  --yes --force-yes   zlib1g-dev zlib1g
sudo apt-get install   --yes --force-yes libpcre3 libpcre3-dev
sudo apt-get install   --yes --force-yes libssl-dev
sudo apt-get install  --yes --force-yes  libperl-dev
sudo apt-get install  --yes --force-yes  libgd-dev
sudo apt-get install  --yes --force-yes  net-tools


cd ${ngx_dir}/ ; mkdir -pv ${ngx_tar_dir} ${ngx_dev_kit_tar_dir} ${setnginx_tar_dir}
wget -O  ${ngx_tar_dir}.tar.gz          http://nginx.org/download/${ngx_tar_dir}.tar.gz
tar --strip-components 1 -C ${ngx_tar_dir}         -xzvf ${ngx_tar_dir}.tar.gz
wget -O  ${ngx_dev_kit_tar_dir}.tar.gz  https://github.com/simpl/ngx_devel_kit/archive/${ngx_dev_kit_ver}.tar.gz
tar --strip-components 1 -C ${ngx_dev_kit_tar_dir} -xzvf ${ngx_dev_kit_tar_dir}.tar.gz
wget -O  ${setnginx_tar_dir}.tar.gz     https://github.com/openresty/set-misc-nginx-module/archive/${setnginx_ver}.tar.gz
tar --strip-components 1 -C ${setnginx_tar_dir}    -xzvf ${setnginx_tar_dir}.tar.gz

cd ${ngx_dir}/${ngx_tar_dir}
 ./configure --prefix=/opt/nginx \
     --with-http_ssl_module \
     --add-module=${ngx_dir}/${ngx_dev_kit_tar_dir}/ \
     --add-module=${ngx_dir}/${setnginx_tar_dir}/

sudo  make -j2
sudo  make install

