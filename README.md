# dns_find_and_replace

Ever need to CTRL-F find and replace all on your DNS zones? Well I need to replace thousands of DNS records en masse on Cloudflare and wanted a quick and simple way to do this, so here you go. While I have tested this and even used it in production, I can't guarantee this won't break YOUR production, so please take proper precautions.

## Usage
Create a *scoped* API key within Cloudflare with permission to edit whatever DNS zones you desire to find and replace on. This will edit _all_ zones you give your API token access to, so please be sure that's what you want.

Your scoped API key should be available under the environment variable `CLOUDFLARE_API_KEY`.  Here's an example of using to replace values:

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
