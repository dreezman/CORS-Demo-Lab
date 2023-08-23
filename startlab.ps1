get-job| stop-job | remove-job 
$pwd=get-location
start-job -scriptblock {param($pwd,$name,$port) set-location $pwd; go run main.go $name $port} -ArgumentList $pwd,"static/config.json"
start-job -scriptblock {param($pwd,$name,$port) set-location $pwd; go run main.go $name $port} -ArgumentList $pwd,"static/config.json""
start-job -scriptblock {param($pwd,$name,$port) set-location $pwd; go run main.go $name $port} -ArgumentList $pwd,"static/config.json"
write-host "`n"
echo "In 20 seconds, use this URL: http://localhost:9081/"