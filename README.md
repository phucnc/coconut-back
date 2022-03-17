## Requirements
### Docker
* Install Docker Engine.
* Install docker-compose.

### Google cloud credential
Define credential in google cloud storage service and download it as json file

Ex:
```json
{
  "type": "service_account",
  "project_id": "example",
  "private_key_id": "example",
  "private_key": "example",
  "client_email": "example@email.com",
  "client_id": "1234567",
  "auth_uri": "https://example.com/auth",
  "token_uri": "https://example.com/token",
  "auth_provider_x509_cert_url": "https://www.example.com/certs",
  "client_x509_cert_url": "https://www.example.com/certs"
}
```

This credential json file path will be provided for running app

## Database
App using postgres as main database

Migration files are in `./migrations` dir

A postgres instance must be hosted and connection info must be defined as environments for app to connecting

## Bash Scripts
There are some bash scripts in repo for running specified tasks.

Make sure your terminal is in same folder with this **README** when execute them.

### Generate golang bindings from smart contract ABI
Get docker image for ethereum client  

```bash
docker pull ethereum/client-go:alltools-stable
```

Copy
* NFT contract ABI to `./abi/nft.json`
* NFT exchange contract ABI to `./abi/exchange.json`

Then execute 
```bash
./gen_go_bindings
```

The go bindings are generated in `./contracts` dir

### Run on local without docker

Execute 
```bash
./local
```

Please take a look at `./local` to see more details about required environments. 

Required environments:
* Postgres connection info
* Google cloud storage credential and storage bucket name
* Server port for running
* NFT contract address
* NFT exchange contract address