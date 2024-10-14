to test stuff:
cd kademlia
go test -coverprofile="name"
to visualize it:
go tool cover -html="name"



to run the program:
go to a terminal that is not in VScode
$env:Path += ";C:\Program Files\Docker\Docker\resources\bin"
.\start_containers.bat
docker attach container# to enter node number "#"
ctrl + p and then ctrl + Q to exit a node without it stopping
once in the node you can write:
"put "File"" to store a string "File"
"get "Hash"" to find the node that has the file and then the value of the file
"exit" to exit a node, and also stopping it

to see all open containers:
docker ps 
To stop the program: 
docker stop $(docker ps -a -q)
