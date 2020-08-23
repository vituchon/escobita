#!/bin/bash

folder=${1:-$(pwd)}
while true;
do
        echo -e "Setting watch on $folder"
        eventsToWatch="-e create -e modify -e moved_to"
        file=`inotifywait -q -r --format '%w%f' ${eventsToWatch} "${folder}"`
        echo " Detected changes on ${file}..."
        if [ ${file: -5} == ".less" ]; then
                lessc --strict-imports $file  `dirname $file | sed -e "s/less/less_compiled/"`/`basename $file | sed -e "s/less/css/"`
        fi
done