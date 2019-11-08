echo "proxy to stag1"
port=$1
echo "get port:$port"

while :
do 
echo "`date` start proxy $port"
ssh -oProxyCommand="nc -x 127.0.0.1:1086 %h %p" -L $port:127.0.0.1:$port gdmops@10.6.3.251
echo "`date` proxy $port disconnect  "
sleep 2
done

11
