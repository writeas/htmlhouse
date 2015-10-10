
all : local

install : 
	cd css/; $(MAKE) install $(MFLAGS)
	
local : force_look
	cd css/; $(MAKE) $(MFLAGS)
	
development : force_look
	cd css/; $(MAKE) development $(MFLAGS)
	
production : force_look
	cd css/; $(MAKE) production $(MFLAGS)
	
force_look : 
	true
