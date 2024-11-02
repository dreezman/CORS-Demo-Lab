```
http://localhost:9081/xss-attack?xssvalue=<script>alert("XSS attack!")</script>
document.getElementsByTagName('*')
document.getElementsByTagName("*")[2].innerHTML
document.getElementsByTagName("*")[2].outerHTML

http://localhost:9081/xss-attack?xssvalue="<script>"alert("XSS attack!")</script>
document.getElementsByTagName('*')
document.getElementsByTagName("*")[3].innerHTML
document.getElementsByTagName("*")[3].outerHTML

inner HTML= text contents of the HTML without DOM elements like <p>,<a>, <div>
outer HTML=HTML DOM Element in the DOM, content + HTML elements like <p>,<a>, <div>
```

https://nodeployfriday.com/posts/cors-cyber-attacks/

http://localhost:9081/xss-attack?xssvalue=<script src=http://localhost:3000/xss/xss.js> </script>
document.getElementsByTagName('\*')
document.getElementsByTagName("\*")[2].innerHTML
document.getElementsByTagName("\*")[2].outerHTML

same site created before javascript to protect cookies
cookies controls are protected by samesite NOT same origin
HTML only controls (does not restrict Javascript)

BROWSER controls over cookies management vs JAVA script control over reading reponses which runs in browser created by unknown developer (access to cookies, local storage, JWT, data in responses, JS PUT and POSTS)
cors HTTP requests
csp content like javascript
