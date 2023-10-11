# config
# export JAVA_HOME=/Users/wangyuan/Library/Java/JavaVirtualMachines/corretto-1.8.0_352/Contents/Home
export ANDROID_HOME=${HOME}/android/android-sdk/cmdline-tools
export ANDROID_NDK_HOME=${HOME}/android/android-ndk-r25
export TOOL=${HOME}/android/android-ndk-r25

android: depend
	gomobile bind -target=android/arm64 -androidapi=23 -o device.aar -v -x ./

ios: depend
	gomobile bind -target=ios -o device.xcframework -v ./

depend:
	cd ./
	gomobile init
	go get golang.org/x/mobile/bind

clean:
	rm -rvf libdevice.*