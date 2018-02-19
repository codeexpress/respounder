<img src="https://i.imgur.com/o1Gya0D.png" width=300px />

# res·pound·er
   <span>/rɪˈspaʊnd dər/</span></span>
   <i>noun</i>
   <ul>
   <li>
   <div style="margin-left:10px; display:inline;">
   A tool that detects presence of a <a href=https://github.com/SpiderLabs/Responder>Responder</a> in the network
   </div>
   </li>
   <li>
   <div style="margin-left:10px; display:inline;">
   Identifies compromised machines before hackers run away with the loot (hashes)
   </div>
   </li>
   </ul>

   Respounder sends LLMNR name resolution requests for made-up hostnames that do not exist.
   In a normal non-adversarial network we do not expect such names to resolve.
   However, a responder, if present in the network, will resolve such queries
   and therefore will be forced to reveal itself.

## Download

### Latest Releases
Respounder is available for 32/64 bit linux, OS X and Windows systems.
Latest versions can be downloaded from the
[Release](https://github.com/codeexpress/respounder/releases) tab above.

### Build from source
This is a golang project with no dependencies. Assuming you have golang compiler installed,
the following will build the binary from scratch
```
$ git clone https://github.com/codeexpress/respounder
$ cd respounder
$ go build -o respounder respounder.go
```

## Usage

Running `respounder` is as simple as invoking it on the command line.
Example invocation:
```bash
$ ./respounder


     .´/
    / (           .----------------.
    [ ]░░░░░░░░░░░|// RESPOUNDER //|
    ) (           '----------------'
    '-'

[wlan0]    Sending probe from 192.168.0.19...   responder not detected
[vmnet1]   Sending probe from 172.16.211.1...   responder not detected
[vmnet8]   Sending probe from 172.16.55.1...    responder detected at 172.16.55.128
```

### Flags

```
$ ./respounder [-json] [-debug] [-hostname testhostname | -rhostname]

Flags:
  -json
        Prints a JSON to STDOUT if a responder is detected on
        the network. Other text is sent to STDERR
  -debug
        Creates a debug.log file with a trace of the program
  -hostname string
        Hostname to search for (default "aweirdcomputername")
  -rhostname
        Searches for a hostname comprised of random string instead
        of the default hostname ("aweirdcomputername")
```

### Typical usage scenario

#### Personal
Detect rogue hosts running responder on public Wi-Fi networks
e.g. like airports, cafés and avoid joining such networks
(especially if you are running windows OS)

#### Corporate
Detect network compromises as soon as they happen by running respounder
in a loop

For eg. the following `crontab` runs respounder every minute and logs a JSON file to syslog
whenever a responder is detected.
```bash
* * * * * /path/to/respounder -json | /usr/bin/logger -t responder-detected
```

Example `syslog` entry:
```bash
code@express:~/$ sudo tail -f /var/log/syslog
Feb  9 03:44:07 responder-detected: [{"interface":"vmnet8","responderIP":"172.16.55.128","sourceIP":"172.16.55.1"}]
```

## Demo
![Respounder in action](https://i.imgur.com/ymcDRnJ.gif)

## Coming Up Next: Android App
There are plans to port this tool to an android app so that adversarial Wi-Fi networks
(eg. WiFi Pineapple or WiFi Pumpkin running responder) can be
detected right from a mobile phone.
