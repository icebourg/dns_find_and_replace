# dns_find_and_replace

Ever need to CTRL-F find and replace all on your DNS zones? Well I need to replace thousands of DNS records en masse on Cloudflare and wanted a quick and simple way to do this, so here you go. While I have tested this and even used it in production, I can't guarantee this won't break YOUR production, so please take proper precautions.

The goal of this project is to easily and safely replace many, many records in Cloudflare DNS. Let's say you have 1,000 A records pointing to 1.2.3.4 but need to change them to CNAME them to `app.example.com`, this project easily and safely performs that step en masse.

## Usage
Create a *scoped* API key within Cloudflare with permission to edit whatever DNS zones you desire to find and replace on. This will edit _all_ zones you give your API token access to, so please be sure that's what you want.

Your scoped API key should be available under the environment variable `CLOUDFLARE_API_KEY`.  Here's an example of using to replace values:

### Arguments
- `find` The *content* of the DNS records you want to replace.
- `replacement` The new value for the replacement DNS record.
- `replacementtype` The type of the replacement DNS record.
- `exclude` The name of the record to exclude from being replaced, optional. (can only specify one)
- `batchsize` How many DNS entries to replace at a time. (default: 100)

### Example Usage
```
dns_find_and_replace -find app.example.com -replacement test.example.com -replacementtype CNAME -exclude exclude.example.com
Found record test1.example.com -> app.example.com (CNAME)
Found record testing0.example.com -> app.example.com (CNAME)
Found record testing100.example.com -> app.example.com (CNAME)
Found record testing101.example.com -> app.example.com (CNAME)
Found record testing102.example.com -> app.example.com (CNAME)
Ready to commit 100 replacements? Yes to proceed, any other value to stop...Yes
Replacing test1.example.com
Replacing testing0.example.com
Replacing testing100.example.com
Replacing testing101.example.com
Replacing testing102.example.com
```

I've tested with several thousand records to good effect.

### Installation
Download and install [Go](https://go.dev/doc/install), then it's as easy as:

```
% go install github.com/icebourg/dns_find_and_replace@latest
go: downloading github.com/icebourg/dns_find_and_replace v0.0.0-20230511230158-2969184a79ef
% dns_find_and_replace --help
Usage of dns_find_and_replace:
  -batchsize int
    	the number of DNS records we will operate on at a given time (default 100)
  -exclude string
    	a DNS record to exclude from the replacement
  -find string
    	the DNS record content we are searching for (default "REQUIRED")
  -replacement string
    	the content of the record we want to replace with (default "REQUIRED")
  -replacementtype string
    	the type of the record we want to create (default "REQUIRED")
```

If `dns_find_and_replace` doesn't find the binary after a successful `go install ...`, make sure your Go installation directory is in your `$PATH`. Run `go env` to find where go is installing your binaries.
