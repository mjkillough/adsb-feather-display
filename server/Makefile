data:
	test -d data || mkdir data
	wget -O data/StandingData.sqb.gz http://www.virtualradarserver.co.uk/Files/StandingData.sqb.gz
	gunzip data/StandingData.sqb.gz

run:
	go run ./main.go

bin/adsb-server:
	go build -o bin/adsb-server

build: bin/adsb-server
