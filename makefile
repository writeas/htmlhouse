
all : local

install : 
	./keys.sh gen
	cd static/; $(MAKE) install $(MFLAGS)
	
local : force_look
	cd static/; $(MAKE) $(MFLAGS)
	
development : force_look
	cd static/; $(MAKE) development $(MFLAGS)
	
production : force_look
	cd static/; $(MAKE) production $(MFLAGS)
	./build.sh

hidden : production
	./build-hidden.sh
	
force_look : 
	true
