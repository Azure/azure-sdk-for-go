#!/bin/bash
set -ex

today=`date "+%Y%m%d"`

firstDay=`date -d "${today}" +%Y%m01`

week=`date -d "$firstDay" +%w`

secondSaturday=$((firstDay+(12 - week) % 7 + 8))

if [ $today -gt $secondSaturday ]
then
 echo "The PR generation time of the current month is: [$firstDay-$secondSaturday]"
 exit 1
fi