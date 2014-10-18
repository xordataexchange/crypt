# crypt

You can use crypt as a command line tool or as a configuration library:

* [crypt cli](bin/crypt)
* [crypt/config](config)

## Generating gpg keys and keyrings

The crypt cli and config package require gpg keyrings. 

### Create a key and keyring from a batch file

```
vim app.batch
```

```
%echo Generating a configuration OpenPGP key
Key-Type: default
Subkey-Type: default
Name-Real: app
Name-Comment: app configuration key
Name-Email: app@example.com
Expire-Date: 0
%pubring pubring.gpg
%secring secring.gpg
%commit
%echo done
```

Run the following command:

```
gpg2 --batch --gen-key app.batch
```

You should now have two keyrings, pubring.gpg which contains the public keys, and secring.gpg which contains the private keys.

> Note the private keys is not protected by a passphrase.
