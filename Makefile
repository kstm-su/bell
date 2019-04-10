install: dependencies
	go build
	mkdir -p /home/kstm/go/bin/bell
	mv bell /home/kstm/go/bin/bell
	sudo cp init_script.sh /etc/init.d/bell
	sudo /etc/init.d/bell start
	sudo echo '@reboot /etc/init.d/bell start' >> /var/spool/cron/crontabs/root

dependencies:
	go get -u github.com/Narsil/alsa-go
	go get -u github.com/youpy/go-wav
	go get -u github.com/go-gonic/gin
