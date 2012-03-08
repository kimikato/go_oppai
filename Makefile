GOC		 = 6g
GOL		 = 6l
CFLAGS	= -O
PROGRAM = oppai

all:	$(PROGRAM)

oppai: oppai.6
	$(GOL) -o $(PROGRAM) $^

oppai.6: oppai.go
	$(GOC) $^

oppai2: oppai2.6
	$(GOL) -o $(PROGRAM)2 $^

oppai2.6: oppai2.go
	$(GOC) $^

clean:
	$(RM) -f *.6 *~ $(PROGRAM) $(PROGRAM)2

