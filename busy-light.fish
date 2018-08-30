set -x BUSYLIGHT_SERVER localhost:5050
set id_hostname (hostname)
set id_shell_pid %self
set -x BUSYLIGHT_SHELL_ID "$id_hostname-$id_shell_pid"

function preexec_hook --on-event fish_preexec
  /home/mrhn/extcode/busy-light/src/busy-light-send-event/busy-light-send-event -server=$BUSYLIGHT_SERVER -shellid=$BUSYLIGHT_SHELL_ID -event=start -pwd="$PWD" -cmdline="$argv" 
end

function postexec_hook --on-event fish_postexec
  /home/mrhn/extcode/busy-light/src/busy-light-send-event/busy-light-send-event -server=$BUSYLIGHT_SERVER -shellid=$BUSYLIGHT_SHELL_ID -event=stop -pwd="$PWD" -cmdline="$argv" -exitcode=$status
end