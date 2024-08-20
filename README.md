# bbounty
A tool to make a structure for your bug bounty programs

# Usage examples

```bash
Usage:
  bbounty [command]

Available Commands:
  add         Add a program with domain names
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  list        List discovered domains

Flags:
  -h, --help   help for bbounty

Use "bbounty [command] --help" for more information about a command.

~ ❯ bbounty add
Enter program type (vdp or bbp): bbp
Enter program name: HackerOne
Enter domain names (space-separated): hackerone.com hackertwo.com

... an output from subfinder with domains it discovered
```

The following structure was just created:

```
bbp
|
|--- HackerOne
|   |
|   |--- hackerone.com
|   |   |
|   |   |--- domains.txt
|   |
|   |--- hackertwo.com
|       |
|       |--- domains.txt
```

To display all the domains in directories and their subdirectories, you can run the `bbounty list` command.


```fish
~ ❯ ls -la
total 0
drwxr-xr-x@  3 paractmol  staff   96 Aug 20 21:39 ./
drwxr-xr-x@ 18 paractmol  staff  576 Aug 20 15:20 ../
drwxr-xr-x@  3 paractmol  staff   96 Aug 20 19:28 bbp/

~ ❯ bbounty list
zendesk4.hackerone.com
hackerone.com
mta-sts.forwarding.hackerone.com
...
hackertwo.com
```

# Install

`go install -v github.com/paractmol/bbounty`