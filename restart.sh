mode=production
echo original parameters=[$*]
echo original OPTIND=[$OPTIND]
while getopts ":m:" opt
do
    case $opt in
        m)
            echo "this is -a option. OPTARG=[$OPTARG] OPTIND=[$OPTIND]"
			mode=$OPTARG
            ;;
        ?)
            echo "no valid parameter found."
            ;;
    esac
done

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
echo start mode is $mode
sudo nohup ./homepage.exe -mode $mode >>log.txt 2>&1 &
sleep 5
process_id=$(ps aux | grep homepage | grep -v 'sudo' | grep -v 'grep' | awk '{print $2 }')
echo new_process_id is $process_id
echo server started
