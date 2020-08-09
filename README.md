# VK Audiomessage Sender

This tiny script helps you send audio message to person at [VK](https://vk.com).

**Attention!** Your application have to have [direct authorizaion](https://vk.com/dev/auth_direct) ability!

**Attention!** This script supports **only** 2fa!

## Setup

You should set up this Auth fields, hardcoded in `main.go`:
1. `clientID` — your application has client ID
1. `clientSecret` — your application has client secret key; you have to have [direct authorizaion](https://vk.com/dev/auth_direct) ability!
1. `username` — your username: phone or email
1. `password` — your password in plain text 🤷‍♂️
1. `filename` — your `.ogg` filename (you [should use **mono**,16KHz,16Kb/s audio](https://vk.com/dev/upload_files_2)) without path
1. `recipient` — recipient user ID

## Usage

This script was tested on Go 1.14. You have to use application with enabled [direct authorizaion](https://vk.com/dev/auth_direct).

Example
```go
package main

import sender "github.com/dimon4ezzz/vk-audiomessage-sender"

func main() {
    auth := sender.Auth{
        ClientID: 111222,
        ClientSecret: "Z6jB1ka2uinYsHhHbZxr",
        Username: "personal@mail.com",
        Password: "dolphins42",
        Filename: "audiomessage.ogg",
        Recipient: 11235813
    }
    sender.Send(auth)
}
```

Please, wait for 2fa code prompt and say to shell this code.

At the end of this script you will see link to your message at full VK Web version.