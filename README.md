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
  -h, --help      help for bbounty
  -v, --verbose   Verbose mode

Use "bbounty [command] --help" for more information about a command.

~ ❯ bbounty add
Enter program type (vdp or bbp): vdp
Enter program name: HackerOne
Enter domain names (press Enter to finish):
hackerone.com
hackertwo.com

... an output from subfinder with domains it discovered
```

Alternatively, you can specify the program type and name using these options.
```bash
~ ❯ bbounty add -p bbp -n HackerOne
Enter domain names (press Enter to finish):
hackerone.com
hackertwo.com

```

The following structure was just created:

```fish
❯ tree bbp -n 3
bbp
└── HackerOne
    ├── hackerone.com
    │   └── domains.txt
    └── hackertwo.com
        └── domains.txt
```

To display all the domains in directories and their subdirectories, you can run the `bbounty list` command.


```fish
drwxr-xr-x@  3 paractmol  staff   96 Aug 20 19:28 bbp/

~ ❯ bbounty list
zendesk4.hackerone.com
hackerone.com
mta-sts.forwarding.hackerone.com
...
hackertwo.com
```

## Unix pipeline

```
~ ❯ cat hackerone.csv | grep URL | cut -d, -f1 | bbounty --verbose add -p bbp -n HackerOne
Enter domain names (press Enter to finish):
Executing: echo ma.hacker.one | tee domains.txt
ma.hacker.one

Executing: echo support.hackerone.com | tee domains.txt
support.hackerone.com

Executing: echo hackathon-photos.hackerone-user-content.com | tee domains.txt
hackathon-photos.hackerone-user-content.com
```

# Config

You can set a custom command line that will be triggered when you add a new program with the `add` command.

```bash
~ ❯ cat  ~/.config/bbounty/config.yml
command: "echo %s | tee domains.txt"
```


# Install

`go install -v github.com/paractmol/bbounty@latest`