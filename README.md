# slackchess
How can you play chess on any device?  

Slack!  

This project gives you a turn key solution for tranforming slack into a chess client. slackchess is built with:
- full [chess engine](https://github.com/loganjspears/chess)
- draw and resignation support
- PGN output
- challenge another player
- previous move highlighting

The roadmap includes:
- play against [Stockfish](https://stockfishchess.org)

## Screenshot
<img src="https://raw.githubusercontent.com/loganjspears/slackchess/master/screen_shots/screen_shot_1.png" width="600">

## Installation Guide

To start the Slack Integration Guide you will need the IP or URL of your server.  If you don't already have a place to host slackchess, you can follow the Digital Ocean Setup Guide first.  After the Slack Integration is setup, you should follow the Server Setup Guide.  The Server Setup Guide will guide you through installing slackchess on a server with docker installed.  

### Digital Ocean Setup Guide

1. Signup for Digital Ocean if you don't have an account (https://m.do.co/c/f4609bed935c referal link to get $10 free)
2. Click Create Droplet
3. Select One-click Apps > Docker 
4. Choose smallest size
5. Make sure to add your SSH keys (if you need help click "New SSH Key > How to use SSH keys")
6. Select 1 droplet
7. Click Create

### Slack Integration Guide

1. Login to Slack and go to https://slack.com/apps
2. Go to Configure > Custom Integrations > Slash Commands > Add Configuration
3. For "Choose a Command" type "/chess" and press "Add Slash Command Integration"
4. Set "URL" to http://45.55.141.331/command where "45.55.141.331" is your IP
5. Make sure "Method" is POST
6. Copy and paste the generated "Token" somewhere, you will need it later
7. For "Customize Name" you can enter anything (ex. "ChessBot")
8. For "Customize Icon" I used this image: https://upload.wikimedia.org/wikipedia/commons/thumb/f/f0/Chess_kdt45.svg/45px-Chess_kdt45.svg.png
9. Click "Save Integration"

![slack integration](/screen_shots/screen_shot_3.png)

### Server Setup Guide

If you added your SSH keys to your local machine (you should have in the Digital Ocean Setup Guide), you can SSH into your machine.  Replace "45.55.141.331" with your IP. 
```bash
ssh root@45.55.141.331
```

Pull down the docker image from Docker Hub.
```bash
docker pull loganjspears/slackchess
```
 
Create and run a docker container with your information. Replace "R546Sk2RoZiAXZltssJ4WfEO" with your Slack token.  Replace "45.55.141.331" with your IP.  
```bash
docker run -d -p 80:5000 -e TOKEN=R546Sk2RoZiAXZltssJ4WfEO -e URL=http://45.55.141.331 loganjspears/slackchess
```

If everything worked the ps command should show the container running.  
```bash
docker ps
```

Thats it!

## Commands

You can view all commands by using the help command.
```
/chess help
```
 
![slackchess](/screen_shots/screen_shot_2.png)

