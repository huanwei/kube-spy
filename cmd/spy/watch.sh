#!/usr/bin/env bash
# arg 1: pod name
echo " ">/tmp/watch
begin=$(date +%s)
last=""
count=0

while true;do
    # get annotation
    now=$(kubectl describe pod http-test|grep chaos= )

    # get time
    end=$(date +%s)
    time=`expr $end - $begin`

    # whether newer
    if [ "$now" != "$last" ];then
        echo "$time s:">>/tmp/watch
        echo "$now"  >>/tmp/watch
        last="$now"
    fi

    # check end
    if [ "$now" = "$last" ];then
        count=`expr $count + 1`
        if test $[count] -gt 50;then break;fi
    else
        count=0
    fi

    sleep 0.1
done


