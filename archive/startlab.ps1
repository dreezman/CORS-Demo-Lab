get-job| stop-job | remove-job 
$dir=get-location
$results=Start-Job -ScriptBlock { param($dir,$configfile) cd $dir; go run main.go  $configfile} -ArgumentList $dir,"static/config.json" 
write-output "`n"
write-output "In 20 seconds, use this URL: http://localhost:9081/ "
start-sleep -s 10
write-output "`n"
Receive-Job -job $results

