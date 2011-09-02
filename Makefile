include $(GOROOT)/src/Make.inc

TARG=launchpad.net/lpad

GOFILES=\
	blueprints.go\
	branches.go\
	bugs.go\
	oauth.go\
	people.go\
	projects.go\
	session.go\
	value.go\


include $(GOROOT)/src/Make.pkg

runexample: _obj/launchpad.net/lpad.a
	$(GC) -I _obj -o example.$(O) example.go
	$(LD) -L _obj -o example example.$(O)
	./example
