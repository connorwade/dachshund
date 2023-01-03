## This tool is in Alpha!
## Expect breaking changes!
---------
# Dachshund
A tenacious crawler built for quality testing websites easily and quickly.

## Quick Start
Currently, this tool is only distributed as a Go script. In the future, I hope to make it downloadable as a command-line tool.

1. clone repo
2. `cd dachshund` enter the directory
3. `go build` to build the tool
4. `./dachshund config <yourdomain>` to create a config file
5. `./dachshund crawl` to start your site crawl
6. `./dachshund report` to create a csv report of broken links

## Current Features

Dachshund currently supports crawling websites to check links, image srcs, and text content.

Dachshund has three subcommands:
1. dachshund crawl: will crawl the website defined in the YAML configuration file. Use the flag "--report", to instantly write a report for the crawl.
2. dachshund report: writes a report from a JSON file. Currently only supports CSV files
3. dachshund config <yourwebsite>: creates a YAML configuration file for you

### YAML configuration file

Call your configuration file "dachshund.yaml". Please put it in a directory you want your dachshund files to be in. It doesn't currently support a robust file system.

```yaml
starterURL: www.<yourwebsite> # The starting URL
allowedDomains: # A list of websites your crawler is allowed to visit
    - www.<yourwebsite>
selectors:
    get-content: # Selectors for HTML elements you want the text content from
        - h1
    check-links: # Selectors for HTML elements who's links you'd like to visit
        - a[href]
        - img[src]
Colly: # Colly defined variables
    maxDepth: 0 # The max-depth you'd like the crawler to crawl on a website (0 for inifinite, 1 for just the starting URL, 2 for all the links on the starter URL, and so on)
    async: true # Whether to run Colly asynchronously (sends more requests at the same time)
    parallelRequests: 2 # How many asynchronous requests are allowed at a time (CAUTION: do not set too high as you can create significant load to a server)
```
