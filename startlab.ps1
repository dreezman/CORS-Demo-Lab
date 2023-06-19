get-job| stop-job | remove-job 
start-job -scriptblock {cd C:\git\gostuff\CORS-Demo-Lab; go run C:\git\gostuff\CORS-Demo-Lab\main.go TLD 8081}
start-job -scriptblock {cd C:\git\gostuff\CORS-Demo-Lab; go run C:\git\gostuff\CORS-Demo-Lab\main.go iframe1 3000}
start-job -scriptblock {cd C:\git\gostuff\CORS-Demo-Lab; go run C:\git\gostuff\CORS-Demo-Lab\main.go iframe2 3001}
