# mainvspageenforcement.txt

What is enforcing csp

the front end serving the page
or back end making rest calls and then enforcing

NOTE
1) policy is per page, how the page was loaded the policy was set then
2) test it
3) page 1 - load with policy
4) page 2 - same page as 1, but no policy
5) Execute this to test if you can insert into dom
6) Talk about 2 pages with meta head


# javascript to insert and execute script
var script = document.createElement( "script" );
script.text = 'alert("hi")';
document.head.appendChild( script )





# turn off policy

$session = New-Object Microsoft.PowerShell.Commands.WebRequestSession
$session.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"
$session.Cookies.Add((New-Object System.Net.Cookie("password", "SuperSecretPassword", "/", "localhost")))
Invoke-WebRequest -UseBasicParsing -Uri "https://localhost:9381/set-csp-header" `
-Method "POST" `
-WebSession $session `
-Headers @{
"authority"="localhost:9381"
  "method"="POST"
  "path"="/set-csp-header"
  "scheme"="https"
  "accept"="*/*"
  "accept-encoding"="gzip, deflate, br, zstd"
  "accept-language"="en,en-US;q=0.9,de;q=0.8"
  "origin"="https://localhost:9381"
  "referer"="https://localhost:9381/csp/csp.html"
  "sec-ch-ua"="`"Google Chrome`";v=`"123`", `"Not:A-Brand`";v=`"8`", `"Chromium`";v=`"123`""
  "sec-ch-ua-mobile"="?0"
  "sec-ch-ua-platform"="`"Windows`""
  "sec-fetch-dest"="empty"
  "sec-fetch-mode"="cors"
  "sec-fetch-site"="same-origin"
} `
-ContentType "application/json" `
-Body ([System.Text.Encoding]::UTF8.GetBytes("{$([char]10)  `"enabled`": false,$([char]10)  `"cspMode`": `"`",$([char]10)  `"csp-data`": [$([char]10)    {$([char]10)      `"csp-type`": `"default-src`",$([char]10)      `"domains`": [$([char]10)        `"'none'`",$([char]10)        `"https://localhost:3389`"$([char]10)      ]$([char]10)    }$([char]10)  ]$([char]10)}"))


# turn on policy

# enforce
$session = New-Object Microsoft.PowerShell.Commands.WebRequestSession
$session.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"
$session.Cookies.Add((New-Object System.Net.Cookie("password", "SuperSecretPassword", "/", "localhost")))
Invoke-WebRequest -UseBasicParsing -Uri "https://localhost:9381/set-csp-header" `
-Method "POST" `
-WebSession $session `
-Headers @{
"authority"="localhost:9381"
  "method"="POST"
  "path"="/set-csp-header"
  "scheme"="https"
  "accept"="*/*"
  "accept-encoding"="gzip, deflate, br, zstd"
  "accept-language"="en,en-US;q=0.9,de;q=0.8"
  "origin"="https://localhost:9381"
  "referer"="https://localhost:9381/csp/csp.html"
  "sec-ch-ua"="`"Google Chrome`";v=`"123`", `"Not:A-Brand`";v=`"8`", `"Chromium`";v=`"123`""
  "sec-ch-ua-mobile"="?0"
  "sec-ch-ua-platform"="`"Windows`""
  "sec-fetch-dest"="empty"
  "sec-fetch-mode"="cors"
  "sec-fetch-site"="same-origin"
} `
-ContentType "application/json" `
-Body ([System.Text.Encoding]::UTF8.GetBytes("{$([char]10)  `"enabled`": true,$([char]10)  `"cspMode`": `"Content-Security-Policy-Report-Only`",$([char]10)  `"csp-data`": [$([char]10)    {$([char]10)      `"csp-type`": `"default-src`",$([char]10)      `"domains`": [$([char]10)        `"'none'`",$([char]10)        `"https://localhost:3389`"$([char]10)      ]$([char]10)    }$([char]10)  ]$([char]10)}"))

