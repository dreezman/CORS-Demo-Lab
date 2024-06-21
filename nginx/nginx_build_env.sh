##
## Config file for nginx build environment
##

export repo_dir=`pwd` # assuming we are executing from the repo root
export ngx_dir="${repo_dir}/nginx/"
export ngx_ver="1.25.5"
export ngx_tar_dir="nginx-${ngx_ver}"
# setnginx is a module that provides a lot of useful functions for nonces
export setnginx_ver="v0.33"
export setnginx_tar_dir="set-misc-nginx-module-${setnginx_ver}"
# ngx_dev_kit is a module that provides a lot of useful functions for substitution
export ngx_dev_kit_ver="v0.3.0"
export ngx_dev_kit_tar_dir="ngx_dev_kit-${ngx_dev_kit_ver}"

# rootdir for the root index.html
export NGINX__ROOTDIR=/home/dreez/repos/Browser-Security-Lab
export NGINX__PORT=4200

