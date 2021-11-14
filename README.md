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

# create a unique bucket (make sure to update cloudbuild.yaml with that value)
export MOOD_TRACKER_BUCKET="<bucket>"
gsutil md -l EUROPE-WEST3 "gs://${MOOD_TRACKER_BUCKET}"

# create a secret for a token
gcloud secrets create MOOD_TRACKER_BOT_TOKEN
echo -n "${MOOD_TRACKER_BOT_TOKEN}" | gcloud secrets versions add MOOD_TRACKER_BOT_TOKEN --data-file=-

# configure running test and deploy on each push to master (make sure to update repo and owner, authenticate with the UI if necessary)
gcloud beta builds triggers create github \
    --name=deploy \
    --repo-name=mood-tracker \
    --branch-pattern="^master$" \
    --repo-owner=nestoroprysk \
    --build-config=cloudbuild.yaml

# grant cloudbuild access sufficient access
export PROJECT_ID=$(gcloud projects list --format=json | jq -r '.[].projectId')
export PROJECT_NUMBER=$(gcloud projects list --format=json | jq -r '.[].projectNumber')
export CLOUDBUILD_SERVICE="${PROJECT_NUMBER}@cloudbuild.gserviceaccount.com"
gcloud projects add-iam-policy-binding ${PROJECT_ID} --member="serviceAccount:${CLOUDBUILD_SERVICE}" --role="roles/secretmanager.secretAccessor"
gcloud projects add-iam-policy-binding ${PROJECT_ID} --member="serviceAccount:${CLOUDBUILD_SERVICE}" --role="roles/cloudfunctions.developer"
gcloud projects add-iam-policy-binding ${PROJECT_ID} --member="serviceAccount:${CLOUDBUILD_SERVICE}" --role="roles/iam.serviceAccountUser"

# install hooks
git config core.hooksPath hooks

# deploy the function manually (or push a commit to master, which is a preferrable option)
gcloud functions deploy MoodTracker \
    --entry-point=MoodTracker \
    --region=europe-west3 \
    --trigger-http \
    --runtime=go116 \
    --timeout=5s \
    --memory=128MB \
    --max-instances=1 \
    --allow-unauthenticated \
    --update-env-vars=MOOD_TRACKER_BOT_TOKEN=${MOOD_TRACKER_BOT_TOKEN},MOOD_TRACKER_BUCKET=${MOOD_TRACKER_BUCKET}

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
