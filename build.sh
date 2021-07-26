#!/usr/bin/env bash

package=$1
if [[ -z "$package" ]]; then
  echo "usage: $0 <package-name>"
  exit 1
fi
package_split=(${package//\// })
package_name=${package_split[-1]}

platforms=("windows/amd64" "windows/386" "darwin/amd64" "linux/amd64" "linux/386")

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    outputpath='./build/'$GOOS'-'$GOARCH
    output_name=$package_name
    tar_name=$GOOS'-'$GOARCH
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
        tar_name+='.zip'
    else
        tar_name+='.tar.gz'
    fi

    env GOOS=$GOOS GOARCH=$GOARCH go build -o $outputpath'/'$output_name $package
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi

    if [ $GOOS = "windows" ]; then
        (cd $outputpath && zip $tar_name $output_name)
    else
        (cd $outputpath && tar -czvf $tar_name $output_name)
    fi
    
    rm $outputpath'/'$output_name
    mv $outputpath'/'$tar_name './build/'$tar_name
    rm -r $outputpath

done