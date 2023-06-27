get-job| stop-job | remove-job 
$pwd=get-location
start-job -scriptblock {param($pwd,$name,$port) set-location $pwd; go run main.go $name $port} -ArgumentList $pwd,"ParentFrame","9081"
start-job -scriptblock {param($pwd,$name,$port) set-location $pwd; go run main.go $name $port} -ArgumentList $pwd,"iframe1","3000"
start-job -scriptblock {param($pwd,$name,$port) set-location $pwd; go run main.go $name $port} -ArgumentList $pwd,"iframe2","3001"
write-host "`n"
echo "In 20 seconds, use this URL: http://localhost:9081/?iframeurl1=http://localhost:3000/iframes.html&iframeurl2=http://localhost:3001/iframes.html"