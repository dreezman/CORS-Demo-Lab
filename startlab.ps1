get-job| stop-job | remove-job 
$dir=get-location
$results=Start-Job -ScriptBlock { param($dir,$configfile) cd $dir; go run main.go  $configfile} -ArgumentList $dir,"static/config.json" 
Receive-Job $results
write-host "`n"
echo "In 20 seconds, use this URL: http://localhost:9081/ "