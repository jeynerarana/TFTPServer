# TFTPServer
This is a FTP server side program that is meant to work with build in FTP client in your system
## Ubuntu
So far the application only work with my Linux machine and may not work work in another OS, the 
program was only tested on the Linux System
# How to Use:
1. Start the program by running `sudo go run <location of go file>` in my case I will use `sudo go run .` Since 
I only have **One** go file in my current working directory. But if you have multiple go files please use `sudo go run nameOfFile.go`. You will 
get a banner about the server. You can leave the tab open and go to another tab and run the client.
2. Start Your Own TFTP Client. Since I am on Linux I will just go to my terminal and type: `tftp` then I will get an output as follows:
`tftp> `. This is where I can connect to the server and run commands such as `get` or `put`. The `get` commands sends a read request, and the
`put` command sends a write request. 
3. Connect to the TFTP server with the Linux TFTP client. Navigate to the other tab or terminal windows where you are running your client and
connect to your server using `connect <yourserverIP>`, the terminal should look like so : `tftp> connect 192.168.0.15`
4. Use `get <filename>` or `put <filename>` to send and recieve data
### Note: When you send a file to the server it will be saved in `files/` directory where the program is running. So please make sure that the `files/` is present when you run the server.

