# Monzo Trx-Downloader
This project helps me download the transaction records for my budgeting tool automatically.

# Prerequisites
You need to have a monzo developer account and create an API client.
After that, you need to run the OAuth steps manually (open https://auth.monzo.com and login, open email, and get the token from the callback URL) to get access token and refresh token for the first time.
I might automate in the future.

# How to use
Clone the repo.
```
git clone
```

Copy the config.json.example to `~/.config/trx-downloader/`.
```
git clone
mkdir -p ~/.config/trx-downloader/
cp config.json.example ~/.config/trx-downloader/config.json
```

Fill the config.json with necessary fields and you are good to go.

# Download
Run the command below and the CSV file will be generated in the same folder:
```
go run .
```

It accepts `-d` option to specify the oldest recods. Default is 30.
```
go run . -d 10 // downloads the recods since 10 days before
```
