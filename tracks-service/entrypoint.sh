#!/bin/sh

mkdir -p /tmp/jave/
cp /usr/bin/ffmpeg /tmp/jave/ffmpeg-amd64-3.5.0
chmod +x /tmp/jave/ffmpeg-amd64-3.5.0
java -jar /opt/app/*.jar &
sleep 5
wait