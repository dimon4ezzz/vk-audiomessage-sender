# VK Audiomessage Sender

This tiny script helps you send audio message to person at [VK](https://vk.com).

**Attention!** Your application have to have [direct authorizaion](https://vk.com/dev/auth_direct) ability!

**Attention!** This script supports **only** 2fa!

## Setup

You should set up this constants, hardcoded in `main.go`:
1. `clientID` — every application has client ID
1. `clientSecret` — every application has client secret key; but you have to have [direct authorizaion](https://vk.com/dev/auth_direct) ability!
1. `username` — your username: phone or email
1. `password` — your password in plain text 🤷‍♂️
1. `filename` — your `.ogg` filename (you [should use mono,16KHz,16Kb/s audio](https://vk.com/dev/upload_files_2)) without path
1. `recipient` — recipient user ID

## Usage

This script was tested on Go 1.14.

Start the script with:
```
go run main.go
```

Please, wait for 2fa code prompt and say to shell this code.

At the end of this script you will see link to your message at full VK Web version.