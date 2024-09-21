# Commands

```
gcredstash manages credentials using AWS Key Management Service (KMS) and DynamoDB

Usage:
  gcredstash [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  delete      Delete a credential from the store
  get         Get a credential from the store
  getall      Get all credentials from the store
  help        Help about any command
  list        List credentials and their version
  put         Put a credential into the store
  setup       Setup the credential store
  template    Parse a template file with credentials

Flags:
  -h, --help           help for gcredstash
  -t, --table string   DynamoDB table to use for credential storage (default "credential-store")
  -v, --version        version for gcredstash

Use "gcredstash [command] --help" for more information about a command.
```
