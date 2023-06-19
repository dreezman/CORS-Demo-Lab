get-job| stop-job | remove-job 
start-job -scriptblock {param($pwd,$name,$port) set-location $pwd; go run main.go $name $port} -ArgumentList $pwd,"ParentFrame","8081"
start-job -scriptblock {param($pwd,$name,$port) set-location $pwd; go run main.go $name $port} -ArgumentList $pwd,"iframe1","3000"
start-job -scriptblock {param($pwd,$name,$port) set-location $pwd; go run main.go $name $port} -ArgumentList $pwd,"iframe2","3001"
