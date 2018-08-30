# slackchess
[![Build Status](https://travis-ci.org/loganjspears/slackchess.svg?branch=master)](https://travis-ci.org/loganjspears/slackchess)

How can you play chess on any device?  

Slack!  

This project gives you a turn key solution for tranforming slack into a chess client. With slackchess you can:
- challenge another slack user
- play against @slackbot powered by [Stockfish](https://stockfishchess.org)
- offer draws and resign
- export your game as a PGN

## Screenshot
<img src="https://raw.githubusercontent.com/loganjspears/slackchess/master/screen_shots/screen_shot_1.png" width="600">

## Installation Guide

To start the Slack Integration Guide you will need the IP or URL of your server.  If you don't already have a place to host slackchess, you can follow the Digital Ocean Setup Guide first.  After the Slack Integration is setup, you should follow the Server Setup Guide.  The Server Setup Guide will guide you through pulling down the [docker image](https://hub.docker.com/r/loganjspears/slackchess/) and running it as a container.  

Alternatively, you can use the Azure Container Instances Setup Guide to run your docker container on Microsoft Azure.

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

### Azure Container Instances Setup Guide

This will walk you through getting the `slackchess` container image running on Azure.  You will need an Azure developer account - check out [this site](https://azure.microsoft.com/en-us/free/) to get started.  You will also need the Azure Command Line Interface (CLI) to communicate with Azure.  Follow [these instructions](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli?view=azure-cli-latest) for your platform.

This guide was prepared by modifying the steps in [this tutorial](https://docs.microsoft.com/en-us/azure/container-instances/container-instances-tutorial-prepare-acr) for setting up a docker container instance on Azure.

1. With the Azure CLI installed, open a terminal (command prompt) window and type: `az login`
  * Follow any login steps required.
2. `docker pull loganjspears/slackchess`
3. Create a resource group for your slackchess resources: `az group create --name slackchess --location eastus`
4. Create an Azure Container Registry in the new resource group using `az acr create --resource-group slackchess --name yourACRName --sku Basic --admin-enabled true`
  * Replace `yourACRName` with a unique name of your choosing.  Replace this value in all subsequent commands as well.
5. `az acr login --name yourACRName`
6. `az acr show --name yourACRName --query loginServer --output table` to get the full loginServer name for your container registry.
  * Use the output from this command to replace the value `yourACRLoginServer.azurecr.io` in subsequent commands.
7. `docker tag loganjspears/slackchess yourACRLoginServer.azurecr.io/slackchess:v1`
8. `docker push yourACRLoginServer.azurecr.io/slackchess:v1`
9. Check `az acr repository list --name yourACRName --output table` to see if your push succeeded.  If so, you will see a `slackchess` container.
10. `az acr credential show --name yourACRName --query "passwords[0].value"` to get the password for your container registry.
11. The next couple steps are a bit of a chicken-and-the-egg situation.  You need two things: the Slack integration webhook token, to give to the Azure container instance as an environment variable.  And the fully qualified domain name (FQDN) for your container, to give to the Slack integration as a URL.  
  * Since you will not be able to create an Azure container instance for the server without a unique label, and you will not know if it's an acceptable name until you try, it's best to start with creating the Slack integration.  So:
    1. Run the first 3 steps in the Slack Integration Guide, stopping when you get to the URL.
    2. Copy the value from the Token field of the Slack Integration.
    3. Choose a unique label for your chess server's fully qualified DNS name.   
    4. Run this container creation command in your terminal, substituting `yourSlackIntegrationToken` with the token field value you copied above, `yourDNSLabel` with your chosen label, and `yourPassword` with the container registry password from earlier.
      * `az container create --resource-group slackchess --name slackchess --image yourACRLoginServer.azurecr.io/slackchess:v1 --cpu 1 --memory 1 --registry-login-server yourACRLoginServer.azurecr.io --registry-username yourACRName --registry-password "yourPassword" --dns-name-label yourDNSLabel --environment-variables TOKEN=yourSlackIntegrationToken URL=http://yourDNSLabel.eastus.azurecontainer.io:5000 --ports 80 5000`
      * If `yourDNSLabel` is already in use, this command will abort and tell you so.  Try again with a different value for `yourDNSLabel`, substituting it in *both* locations in the command.
12. Once your command has succeeded, use the same full URL (including port 5000) from the command as the URL field value in the Slack Integration.  Finish the remaining steps in the Slack Integration Guide.
13. Open that same URL in a browser.  If your container was set up successfully, you will get the response `up` in your browser.

## Commands

Play against user:
```
/chess play @user
```

Play against bot:
```
/chess play @slackbot
```

You can view all commands by using the help command.
```
/chess help
```
 
![slackchess](/screen_shots/screen_shot_2.png)

