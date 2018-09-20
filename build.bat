@echo off

rmdir /S /Q release
cls

echo building...

mkdir release

set GOOS=windows
set GOARCH=386
go build -o launchpad-windows-i386.exe
zip -r9 launchpad-windows-i386.zip fake_tickets launchpad-windows-i386.exe
mv launchpad-windows-i386.zip release/launchpad-windows-i386.zip
del launchpad-windows-i386.exe

set GOOS=windows
set GOARCH=amd64
go build -o launchpad-windows-amd64.exe
zip -r9 launchpad-windows-amd64.zip fake_tickets launchpad-windows-amd64.exe
mv launchpad-windows-amd64.zip release/launchpad-windows-amd64.zip
del launchpad-windows-amd64.exe

set GOOS=linux
set GOARCH=386
go build -o launchpad-linux-i386
zip -r9 launchpad-linux-i386.zip fake_tickets launchpad-linux-i386
mv launchpad-linux-i386.zip release/launchpad-linux-i386.zip
del launchpad-linux-i386

set GOOS=linux
set GOARCH=amd64
go build -o launchpad-linux-amd64
zip -r9 launchpad-linux-amd64.zip fake_tickets launchpad-linux-amd64
mv launchpad-linux-amd64.zip release/launchpad-linux-amd64.zip
del launchpad-linux-amd64

set GOOS=
set GOARCH=
