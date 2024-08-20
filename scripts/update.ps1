# Update main
Write-Host "Updating main"
cd ..\main
go get -u
go mod tidy

# Update backend
Write-Host "Updating backend"
cd ..\backend
go get -u
go mod tidy

# Update chatserver
Write-Host "Updating chatserver"
cd ..\chatserver
go get -u
go mod tidy

# Update spacestation
Write-Host "Updating spacestation"
cd ..\spacestation
go get -u
go mod tidy

# Update pipes
Write-Host "Updating pipes"
cd ..\pipes
go get -u
go mod tidy

# Update pipeshandler
Write-Host "Updating pipeshandler"
cd ..\pipeshandler
go get -u
go mod tidy
