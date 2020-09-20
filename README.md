# email2mattermost
Notify mattermost an email receipt from postfix


## Usage
1. Build the app
   ```bash
   go build -o email2mattermost
   sudo mv email2mattemost /usr/local/bin
   ```
1. Create a Bot account on MatterMost
1. Obtain a secret token
1. Update Postfix aliases
   ```
   your-alias: "| email2mattermost -mattemostURL https://mattermost.example.com/ -channelID your-channel-id -bearerToken your-secret-token -myAddress your-alias@example.com"
   ```
1. Restart Postfix
   ```bash
   sudo newaliases
   sudo postfix reload
   ```
