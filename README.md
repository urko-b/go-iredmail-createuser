# iredmail create email account

## iredmail documentation reference
<div style="width:100px ; height:60px">

[![iredmail Reference](https://gitlab.com/urkob/go-iredmail-createuser/-/raw/main/assets/img/iredmail.png)](https://docs.iredmail.org/sql.create.mail.user.html)
</div>
https://docs.iredmail.org/sql.create.mail.user.html

# Use
	
## Install dependencies
```sh
go get
```

## Set environment variables
```sh
REMOTE_USER=user
REMOTE_SERVER=myremoteserver.com
REMOTE_KEY_PATH=/home/user/.ssh/remote_key
REMOTE_PORT=22
API_PORT=4230
SUPPORT_EMAIL=support@myremoteserver.com
TRUSTED_PROXIES="*"
```

## Run
```sh
go run main.go
```


### License

<a rel="license" href="http://creativecommons.org/licenses/by-nc/4.0/"><img alt="Creative Commons License" style="border-width:0" src="https://i.creativecommons.org/l/by-nc/4.0/88x31.png" /></a><br />This work is licensed under a <a rel="license" href="http://creativecommons.org/licenses/by-nc/4.0/">Creative Commons Attribution-NonCommercial 4.0 International License</a>.