killall main ; sleep 3 ; go run main.go TLD 9081 & go run main.go iframe1 3000 & go run main.go iframe2 3001 & 
echo;echo;echo 
echo "In 10 seconds, use this URL: http://localhost:9081/?iframeurl1=http://localhost:3000/iframes.html&iframeurl2=http://localhost:3001/iframes.html"
echo;echo;echo 