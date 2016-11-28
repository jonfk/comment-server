# comment-server-rs

## external dependencies
These are the dependencies managed manually outside of the build system
```
# on osx
brew install openssl
```

Following the discussion here: https://github.com/sfackler/rust-openssl/issues/255. Openssl needs the following
environment variables exported:
```bash
export C_INCLUDE_PATH=/usr/local/opt/openssl/include
```
