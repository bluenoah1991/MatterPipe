# MatterPipe
Redirect bash output to Mattermost Platform

### How to install to your machine?

##### Step 1

Copy MatterPipe file to /usr/bin  

	sudo cp ./MatterPipe /usr/bin/mp  

##### Step 2

Create configuration file /etc/matterpipe.json  

	echo -e "
	{
	    "address": "your mattermost server address",
	    "clientid": "your user name",
	    "password": "your password",
	    "team": "the team name (not display name)",
	    "channel": "the channel name (not display name)"
	}
	" > /etc/matterpipe.json  

##### Step 3

Okay, Have fun!!!

### How to use it?

	ls -la / | mp  

