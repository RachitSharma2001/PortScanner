### Go Port Scanner

# How to run
Default: ```go run main.go```
Specifically UDP: ```go run main.go -protocol UDP```
Specifically localhost: ```go run main.go -host 127.0.0.1```

### Other Flags
-host     (Ex: 127.0.0.1 or 10.0.0.86)<br/>
-protocol (either TCP or UDP)<br/>
-beg      (port to start scanning)<br/>
-end      (port to end scanning)<br/>

### Why it works
On my computer, here is the result when I do ```netstat -ano```:
<img src="images/netstat.png" style="display: block; margin-left: 0px;" />

On host <b>127.0.0.1</b>, what are the open ports?<br/>Well, lets first look at the ports on host <b>0.0.0.0</b> (where all requests are routed to):<img src="images/host0.png" style="display: block; margin-left: 0px;">
Here, the open ports are: 135, 445, 5040, 7680, 10533, and 49664-49668, 49670.<br/>
Now, lets look at the ports on host <b>127.0.0.1</b>:<img src="images/host127.png" style="display: block; margin-left: 0px;">
Here, the open ports are: 10533, 55989, and 56989. <br/>

When we run ```go run main.go -protocol TCP -host 127.0.0.1```, here is our output: <img src="images/output.png" style="display: block; margin-left: 0px;">
So we are scanning the open ports!