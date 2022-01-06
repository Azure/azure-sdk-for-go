#!/bin/bash

if [ -z $1 ]; then
    echo "Please input inputfile"
    exit 1
fi

if [ -z $2 ]; then
    echo "Please input outputfile"
    exit 1
fi

generator automation-v2 $1 $2

if [ "$?" != "0" ]; then
  echo "Failed to generate code."
  exit 1
fi

cat $2