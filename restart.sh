if [ -f homepage.exe ]; then
    rm homepage.exe
fi
echo 'server compiling...'
go build -o homepage.exe

if [ -f homepage.exe ]; then
	echo 'compile server exe success, begin to restart server'
else
	echo 'complie server failed'
	exit 1
fi

process_id=$(ps aux | grep homepage | grep -v 'sudo' | grep -v 'grep' | awk '{print $2 }')
echo current process_id is $process_id
sudo kill -9 $process_id
sudo ./homepage.exe &
sleep 3
process_id=$(ps aux | grep homepage | grep -v 'sudo' | grep -v 'grep' | awk '{print $2 }')
echo new_process_id is $process_id
echo server started
