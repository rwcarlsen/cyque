#!/bin/bash

if [ $# -eq 0 ]; then
    echo "Usage: wq-update.sh <cctools-source-tarball>"
    exit 1
fi

cctools="$(pwd)/$(echo "$@" | sed -r 's/\.[[:alnum:]]+\.[[:alnum:]]+$//')"

wqdir=$(pwd)/work_queue

# untar and remove possibly prev version of work queue files
tar -xzf $@
rm -rf $wqdir
mkdir $wqdir

# copy over needed files

dttools="auth.c auth_all.c auth_unix.c auth_globus.c auth_kerberos.c auth_hostname.c auth_address.c auth_ticket.c b64_encode.c bitmap.c buffer.c catalog_query.c cctools.c change_process_title.c chunk.c console_login.c clean_dir.c copy_stream.c create_dir.c daemon.c datagram.c debug.c debug_file.c debug_journal.c debug_stream.c debug_syslog.c disk_info.c domain_name.c domain_name_cache.c dpopen.c envtools.c fast_popen.c fd.c full_io.c file_cache.c get_canonical_path.c get_line.c getopt.c getopt_aux.c gpu_info.c hdfs_library.c hash_cache.c hash_table.c hmac.c http_query.c itable.c json.c json_aux.c link.c link_auth.c link_nvpair.c list.c load_average.c memory_info.c mergesort.c md5.c mpi_queue.c nvpair.c nvpair_database.c path.c password_cache.c preadwrite.c process.c random_init.c rmonitor.c rmsummary.c set.c sha1.c sleeptools.c sort_dir.c stringtools.c string_array.c text_array.c text_list.c timestamp.c timer.c unlink_recursive.c uptime.c url_encode.c username.c xxmalloc.c"


CFLAGS=-DINSTALL_PATH='"/home/robert/cctools"' -DBUILD_USER='"robert"' -DBUILD_HOST='"rwc"' -DCCTOOLS_VERSION_MAJOR=4 -DCCTOOLS_VERSION_MINOR=2 -DCCTOOLS_VERSION_MICRO='"2"' -DCCTOOLS_VERSION='"4.2.2-RELEASE"' -DCCTOOLS_RELEASE_DATE='"08/25/2014"' -DCCTOOLS_CONFIGURE_ARGUMENTS='""' -DCCTOOLS_SYSTEM_INFORMATION='"Linux rwc 3.15.8-1-ARCH Fri Aug 1 08:51:42 CEST 2014 x86_64 GNU/Linux"' -DCCTOOLS_OPSYS_LINUX -DCCTOOLS_CPU_X86_64

cd $cctools/chirp/src
cp *.h chirp_global.c chirp_multi.c chirp_recursive.c chirp_reli.c chirp_client.c chirp_matrix.c chirp_stream.c chirp_ticket.c $wqdir/
cd ../../dttools/src
cp *.h $dttools $wqdir/
cd ../../work_queue/src
cp *.h batch_job.c work_queue_catalog.c work_queue_resources.c work_queue.c $wqdir/
cd ../../..
rm -rf $cctools

rm -rf wq_all.*
cat $wqdir/*.c >> wq_all.c
rm -f $wqdir/*.c

# create meta includes
#cd $wqdir
#list=$(ls *.c)
#cd ..
#while read -r header; do
#    echo "#include \"$header\"" >> wq_all.c
#done <<< "$list"

