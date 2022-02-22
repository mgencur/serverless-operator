FROM registry.ci.openshift.org/openshift/openshift-serverless-nightly:serverless-operator-src

# This copies the new content to the image
# Optional: Copy only specific/changed files.
COPY . .