# Description

Mood tracker is a telegram bot for tracking mood.

It's built on top ot gcloud.

# Setup

```bash
# ask father for a new bot
python -m webbrowser https://t.me/botfather
export MOOD_TRACKER_BOT_TOKEN=<new-token>
curl https://api.telegram.org/bot$MOOD_TRACKER_BOT_TOKEN/getMe

# install gcloud
brew install google-cloud-sdk

# create a gcloud project
gcloud init

# create a bucket
export MOOD_TRACKER_BUCKET_URL="<gs://uniquelocation>"
gsutil md -l EUROPE-WEST3 ${MOOD_TRACKER_BUCKET_URL}

# deploy the function
gcloud functions deploy MoodTracker \
    --entry-point=MoodTracker \
    --region=europe-west3 \
    --trigger-http \
    --runtime=go116 \
    --timeout=5s \
    --memory=128MB \
    --max-instances=1 \
    --allow-unauthenticated \
    --update-env-vars=MOOD_TRACKER_BOT_TOKEN=${MOOD_TRACKER_BOT_TOKEN}

# set telegram hooks
curl --data "url=$(gcloud functions describe MoodTracker --region=europe-west3 --format=json | jq -r .httpsTrigger.url)" https://api.telegram.org/bot$MOOD_TRACKER_BOT_TOKEN/SetWebhook
```

# Development

```bash
# setup credentials to execute with the permissions a cloud function has
export NAME="$(gcloud iam service-accounts list --format=json | jq -r '.[].email')"
export GOOGLE_APPLICATION_CREDENTIALS="<location>/gcloudcredentials.json"
gcloud iam service-accounts keys create ${GOOGLE_APPLICATION_CREDENTIALS} --iam-account=${NAME}
```
