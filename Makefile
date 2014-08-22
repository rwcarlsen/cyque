
cyc-cde.tar.gz: sample-sim.xml cyclus
	cde cyclus $<
	rm cyclus.sqlite
	mv cde-package cyc-cde
	tar -czf cyc-cde.tar.gz cyc-cde cyclus
	rm -rf cyc-cde
	rm -f cde.options


