# mango-notify

## Run application
```
nohup ./mango-notify -file=/var/log/auth.log -mailfrom=user@mail.com -pwd=1234 -mailto=myrealmail@mail.com -server=smtp.mail.com:465 -encKey=base64key > server.log 2>&1 &
```

**encKey can be generated when using android app:**
https://github.com/skiarn/Mangodroid-notify

### Example notify on changes in auth log.

#### Setup privileges for user running the application.
##### Install acl 
```
sudo apt-get install acl

Configure acl, open /etc/fstab and add acl example, 
/dev/mmcblk0p2  /               ext4            defaults,noatime,nodiratime,acl 0       0
```
* Setup logroutate privileges.
```
sudo vi /etc/logrotate.d/rsyslog 

/var/log/auth.log
{
        rotate 4
        create 640 root adm
        weekly
        missingok
        notifempty
        compress
        delaycompress
        sharedscripts
        postrotate
                invoke-rc.d rsyslog rotate > /dev/null
                /usr/bin/setfacl -m g:mango:r /var/log/auth.log
        endscript
}
```

**OBS to avoid logging following:**
```
CRON: pam_unix(cron:session): session opened for user root by (uid=0)
CRON: pam_unix(cron:session): session closed for user root
```
Edit: `sudo vi /etc/pam.d/common-session-noninteractive`
Add line: `session     [success=1 default=ignore] pam_succeed_if.so service in cron quiet use_uid` above `session required        pam_unix.so`
Save and restart cron service: `sudo service cron restart`
