
all : local

install : 
	./keys.sh gen
	cd less/; $(MAKE) install $(MFLAGS)
	
local : force_look
	cd less/; $(MAKE) $(MFLAGS)
	
development : force_look
	cd less/; $(MAKE) development $(MFLAGS)
	
production : force_look
	cd less/; $(MAKE) production $(MFLAGS)
	./build.sh

hidden : production
	./build-hidden.sh
	
force_look : 
	true
