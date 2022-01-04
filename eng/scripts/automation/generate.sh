#!/bin/bash

if [ -z $1 ]; then
    echo "Please input inputfile"
    exit 1
fi
echo $1
cat $1

if [ -z $2 ]; then
    echo "Please input outputfile"
    exit 1
fi

echo $2

env
generator automation-v2 $1 $2

#mock for test
# cat > $2 << EOF
# {
#   "packages": [
#     {
#       "packageName": "armagrifood",
#       "result": "succeeded",
#       "path": [
#         "sdk/resourcemanager/agrifood/armagrifood",
#         "rush.json"
#       ],
#       "packageFolder": "sdk/resourcemanager/agrifood/armagrifood",
#       "changelog": {
#         "content": "Feature: something \n Breaking Changes: something\n",
#         "hasBreakingChange": true
#       },
#       "artifacts": [
#         "sdk/agrifood/azure-arm-agrifood-1.0.0.tgz",
#       ]
#     }
#   ]
# }
# EOF

cat $2