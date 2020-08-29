# VK Audiomessage Sender

This tiny script helps you send audio message to person at [VK](https://vk.com).

**Attention!** Your application have to have [direct authorizaion](https://vk.com/dev/auth_direct) ability!

## Setup

You should set up this `Auth` fields in call:
1. `ClientID` — your application has client ID
1. `ClientSecret` — your application has client secret key; you have to have [direct authorizaion](https://vk.com/dev/auth_direct) ability!
1. `Username` — your username: phone or email
1. `Password` — your password in plain text 🤷‍♂️
1. `Filename` — your `.ogg` filename (you [should use **mono**,16KHz,16Kb/s audio](https://vk.com/dev/upload_files_2)) without path
1. `Recipient` — recipient user ID

You can set up this `Setup` fields in call:
1. `SaveOauth` — should the application save VK token to file
1. `OauthFile` — custom filename for token; otherwise file has name `vk-token` (see `defaultOauthFile`)

## Usage

This script was tested on Go 1.14. You have to use application with enabled [direct authorizaion](https://vk.com/dev/auth_direct).

Example (w/o token saving)
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
    sender.Send(auth, Setup{})
}
```

If needed (2fa activated): wait for 2fa code prompt and say to shell this code.

At the end of this script you will see link to your message at full VK Web version.

Token is saved with AES, encrypted with your password.
