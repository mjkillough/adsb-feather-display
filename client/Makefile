SOURCES = $(wildcard *.py)

.PHONY: venv
venv: venv/bin/activate

venv/bin/activate: requirements.txt
	test -d venv || virtualenv venv
	. venv/bin/activate; pip install -Ur requirements.txt
	touch venv/bin/activate

deps: mpy-requirements.txt
	test -d deps || mkdir deps
	micropython -m upip install -p deps/ -r mpy-requirements.txt

deploy: venv deps $(SOURCES)
	( \
		. venv/bin/activate; \
		for f in $(wildcard deps/*/) ; do \
			ampy -p /dev/ttyUSB0 put $$f; \
		done; \
		 \
		ampy -p /dev/ttyUSB0 put uwebsockets/uwebsockets/client.py uwebsockets/client.py; \
		ampy -p /dev/ttyUSB0 put uwebsockets/uwebsockets/protocol.py uwebsockets/protocol.py; \
		 \
		for f in $(SOURCES); do \
			ampy -p /dev/ttyUSB0 put $$f; \
		done \
	)
