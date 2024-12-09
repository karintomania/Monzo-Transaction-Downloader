# Trx-Downloader
This project downloads my banks' transaction records for my budgeting tool automatically.

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
trx-downloader
```
