#!/bin/bash


Telegram_bot_token="6759988959:AAEh920mTsazNnCUqFm66PLM6_WLhcmiXjk"
Telegram_chat_id="954443187"
LOG_FILE="/etc/profile"


# help
display_usage(){
	echo "Usage: $0 [options]"
	echo "options:"
	echo "  -h, --h            Display this help message"
	echo "  -d, --debug        Enable debug mode"
}

# Send notification to Telegram
send_telegram_notification(){
	local message="$1"
	local url="https://api.telegram.org/bot$Telegram_bot_token/sendMessage"
	curl -s -X POST "$url" -d "chat_id=$Telegram_chat_id" -d "text=$message" > /dev/null 2>&1
}


# process PAM
process_pam(){

	while read -r line; do
        if [[ $line == *"pam_unix"*"session opened"* ]]; then
            # Extract relevant information from the log line
            date=$(echo "$line" | awk '{print $1, $2, $3}')
            hostname=$(hostname)
            username=$(echo "$line" | awk -F 'user=' '{print $2}' | awk '{print $1}')
            source_ip=$(echo "$line" | awk -F 'rhost=' '{print $2}' | awk '{print $1}')


            # Prepare notification message
            message="New login event:\n\nDate: $date\nUsername: $username\nHostname: $hostname"
            
            if [[ -n "$source_ip" ]]; then
                message+="\nSource IP: $source_ip"
            fi

            # Send notification
            send_telegram_notification "$message"
        fi
    done
}


while getopts ":hd" option; do
	case $option in
		h | --help)
			display_usage
			exit 0
			;;
		d | --debug)
			DEBUG_MODE=true
			;;
		\?)
			echo "Invalid option: -$OPTARG" >&2
			display_usage
			exit 1
			;;
	esac
done


if [ "$DEBUG_MODE" = true ]; then
	echo "Running in debug mode"
	set -x
fi

if [ "$PAM_TYPE" == "open_session" ]; then
	process_pam
fi


tail -f "$LOG_FILE" | process_pam


