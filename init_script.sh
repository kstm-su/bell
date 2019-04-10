#!/bin/sh

### BEGIN INIT INFO
# Provides:		bell
# Required-Start:	$local_fs $network $remote_fs
# Required-Stop:	$local_fs $network $remote_fs
# Default-Start:	2 3 4 5
# Default-Stop:		
# Short-Description:	kstm bell server
### END INIT INFO

set -e
umask 022

# cmd="/root/go/src/github.com/kstm-su/bell/bell"
cmd="/home/kstm/go/bin/bell"

. /lib/lsb/init-functions

# Are we running from init?
run_by_init() {
    ([ "$previous" ] && [ "$runlevel" ]) || [ "$runlevel" = S ]
}


export PATH="${PATH:+$PATH:}/usr/sbin:/sbin"
export PORT="3000"
export TOKEN="rUjxtBY797fWPlaKWL43A4PB"

case "$1" in
  start)
	log_daemon_msg "Starting bell" "bell" || true
	if start-stop-daemon --start --background --quiet --oknodo --make-pidfile --pidfile /var/run/bell.pid --exec $cmd ; then
	    log_end_msg 0 || true
	else
	    log_end_msg 1 || true
	fi
	;;
  stop)
	log_daemon_msg "Stopped bell" "bell" || true
	if start-stop-daemon --stop --quiet --oknodo --pidfile /var/run/bell.pid; then
	    log_end_msg 0 || true
	else
	    log_end_msg 1 || true
	fi
	;;

  restart)
	log_daemon_msg "Restarting bell" "bell" || true
	start-stop-daemon --stop --quiet --oknodo --retry 30 --pidfile /var/run/bell.pid
	if start-stop-daemon --start --background --quiet --oknodo --make-pidfile --pidfile /var/run/bell.pid --exec $cmdi ; then
	    log_end_msg 0 || true
	else
	    log_end_msg 1 || true
	fi
	;;

  status)
	status_of_proc -p /var/run/bell.pid $cmd bell && exit 0 || exit $?
	;;

  *)
	log_action_msg "Usage: /etc/init.d/bell {start|stop|restart|status}" || true
	exit 1
esac

exit 0
