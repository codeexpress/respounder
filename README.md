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
Latest versions can be downloaded from the [Release](https://github.com/codeexpress/respounder/releases) tab above.

### Build from source
This is a golang project with no dependencies. Assuming you have golang compiler installed,
the following will build the binary from scratch
```
$ git clone blah
$ cd respounder
$ go build respounder
```

## Usage

Running `respounder` is as simple as invoking it on the command line.
The following will display output on the terminal.
```
$ ./respounder
```
To detect a compromise as soon as it happens, **run respounder as a cron job running every minute**

### Flags

```
$ ./respounder [-json] [-debug]

Flags:
  -json
        Prints a JSON to STDOUT if a responder is detected on
        network. Other text is sent to STDERR
  -debug
        Creates a debug.log file with a trace of the program
  -help
        Displays this help
```

## Demo
![Respounder in action](https://i.imgur.com/ymcDRnJ.gif)

