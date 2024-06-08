# Update main
echo "Updating main"
cd main
go get -u
go mod tidy

# Update backend
echo "Updating backend"
cd ../backend
go get -u
go mod tidy

# Update chatserver
echo "Updating chatserver"
cd ../chatserver
go get -u
go mod tidy

# Update spacestation
echo "Updating spacestation"
cd ../spacestation
go get -u
go mod tidy

# Update pipes
echo "Updating pipes"
cd ../pipes
go get -u
go mod tidy

# Update pipeshandler
echo "Updating pipeshandler"
cd ../pipeshandler
go get -u
go mod tidy
