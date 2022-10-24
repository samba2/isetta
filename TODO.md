Fix/ Investigate
=================
- http client timeout, increase at
Warn: Error when trying access https://www.google.com/, error was: Get "https://www.google.com/": context deadline exceeded (Client.Timeout exceeded while awaiting headers)

- naming: was px_proxy, is: ntlm_proxy
- naming: private_dns_server -> cooperate_dns_server ?

Implement
=========
- version number
- configure additional NO_PROXY entries
- default user, powered by new `gsudo` version


One Day
=======
- plugin system, e.g. for docker config support (or maybe just use podman)

