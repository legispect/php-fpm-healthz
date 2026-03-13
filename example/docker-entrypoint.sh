#!/bin/sh
set -eux

if [ "$1" = 'php-fpm' ]; then
	# modify php-fpm pool (www.conf) settings at instance starting up
	# we need our pool to respond health check to fastcgi tool.
	POOL_CONF='/usr/local/etc/php-fpm.d/www.conf'
	if [ -f "$POOL_CONF" ]; then
		sed -i \
			-e "s|^;pm.status_path = /status$|^pm.status_path = /status$|g" \
			"$POOL_CONF"
		echo "[PHP-FPM.conf] successfully swap environment variables"
	else
		echo "[PHP-FPM.conf] www.conf does not exist. Skip to modify"
	fi

	HEALTH_SERVER='/usr/local/bin/health-server'
	CHECK_SCRIPT='/usr/local/bin/php-fpm-healthcheck'
	if [ -x "$HEALTH_SERVER" -a -x "$CHECK_SCRIPT" ]; then
		health-server > /proc/1/fd/1 2> /proc/1/fd/2 &
	else
		echo 'health-server or php-fpm-healthcheck does not exist. liveness probe is unavaliable.'
	fi
fi
exec docker-php-entrypoint "$@"
