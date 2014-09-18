
cctools_dir = cctools

all:
	sed -i '/HAS_SYSTEMD_JOURNAL/d' "$(cctools_dir)/configure"
	sed -i '/-lsystemd/{N;d;}' "$(cctools_dir)/configure"
	cd $(cctools_dir) && ./configure && make

clean:
	cd $(cctools_dir) && make clean

