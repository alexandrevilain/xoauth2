# xoauth2

xoauth2 is a wrapper around golang.org/x/oauth2 which adds access and refresh token storage. 
The main goal of this library is to allow apps reboot without asking user to open the auth url everytime it starts.
It can simplify oauth2 clients integration in a serverless environment.

## Supported storages

- [x] Simple JSON file
- [x] Google Cloud secret manager
- [ ] Hashicorp Vault
- [ ] AWS Secrets Manager

## Example

My next open source project will come with this library, stay tuned!

## Security concerns

Not all storages solutions are equal in terms of security. Choose the solution that suits your requirements. 

## Contributing

All contributions are welcome :) Submit your PR if you want to see another storage solution in this library.
