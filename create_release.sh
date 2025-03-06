#!/bin/bash

package_name="jyggalag"
	
platforms=("darwin/arm64" "darwin/amd64" "linux/amd64" "android/arm64")

for platform in "${platforms[@]}"
do
	platform_split=(${platform//\// })
	GOOS=${platform_split[0]}
	GOARCH=${platform_split[1]}
	output_name='./build/'$package_name'-'$GOOS'-'$GOARCH
	if [ $GOOS == "android" ]; then
		CGO_ENABLED=1
		CC="/etc/android-ndk-r26d/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android28-clang" 
		CXX="/etc/android-ndk-r26d/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android28-clang++"
	fi

	env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name $package
	if [ $? -ne 0 ]; then
   		echo 'An error has occurred! Aborting the script execution...'
		exit 1
	fi
done

