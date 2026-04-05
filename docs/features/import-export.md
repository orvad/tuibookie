# Import & Export

TuiBookie supports backing up and restoring your bookmarks via JSON files.

## Export

Go to **Settings > Export bookmarks** to save a backup:

- Saves to the current working directory as `bookmarks-backup-YYYY-MM-DD-HHMMSS.json`
- The timestamp ensures you never accidentally overwrite a previous backup

## Import

Go to **Settings > Import bookmarks** to restore from a backup:

- Lists `.json` files in the current directory to choose from
- Or enter a file path manually
- Imported bookmarks are **merged** into existing categories -- nothing is overwritten
