include $(GOROOT)/src/Make.inc

TARG=launchpad.net/lpad

GOFILES=\
	oauth.go\
	resource.go\
	session.go\
	people.go\
	bugs.go\


include $(GOROOT)/src/Make.pkg

runexample: _obj/launchpad.net/lpad.a
	$(GC) -I _obj -o example.$(O) example.go
	$(LD) -L _obj -o example example.$(O)
	./example
