# Black-Out Back-up Server
A KISS-Backup Server Solution for keeping your uptime high.

## What is this?
This is a small program written in Go which changes your Cloudflare DNS Entries when a ping to your server fails.

## How to use?
1. Download the latest release from the [releases page]()
2. Create a config file in the same directory as the binary called `config.json`
3. Run the program and configure it as a service (see below)
4. Done!

## Config File
To generate an empty config file, run the program without a config file present. The program will generate a config file with empty values.

## Running as a service
### Linux
1. Create a file in `/etc/systemd/system/` called `blackout.service`
2. Paste the following into the file:
```
[Unit]
Description=Black-Out Back-up Server
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/path/to/binary
ExecStart=/path/to/binary/blackout

[Install]
WantedBy=multi-user.target
```
3. Run `systemctl daemon-reload`
4. Run `systemctl enable blackout --now`

### Windows
1. Create a file called `blackout.bat` in the same directory as the binary
2. Paste the following into the file:
```
@echo off
cd /d %~dp0
start /b blackout.exe
```
3. Open the Task Scheduler
4. Create a new task
5. Set the trigger to `At startup`
6. Set the action to `Start a program`
7. Set the program to `blackout.bat`

### MacOS
1. Create a file called `blackout.command` in the same directory as the binary
2. Paste the following into the file:
```
#!/bin/bash
cd /path/to/binary
./blackout
```
3. Open the Terminal
4. Run `chmod +x blackout.command`
5. Open the `System Preferences`
6. Go to `Users & Groups`
7. Go to `Login Items`
8. Add `blackout.command` to the list

## License
This project is licensed under the GPL-3.0 License - see the [LICENSE](LICENSE) file for details
