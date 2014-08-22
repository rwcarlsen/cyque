
cyc-cde.tar.gz: sample-sim.xml cyclus
	cde cyclus $<
	rm cyclus.sqlite
	mv cde-package cyc-cde
	tar -czf cyc-cde.tar.gz cyc-cde cyclus
	rm -rf cyc-cde
	rm -f cde.options

cyclus_path = $(subst /,\/,$(shell which cyclus))
cyclus: cyclus.in
	sed "s/{{CYCLUS_PATH}}/$(cyclus_path)/" $< > $@
	chmod a+x $@

clean:
	rm -rf cde.options cyc-cde.tar.gz cyclus cde-package cyclus.sqlite

