# Release

Script to release Starship's docker images, without given any arguments,
this script uses the commit id and branch name as the tag:
`<commit_id>_branchname`

```
build_and_push_images.sh
```

You then need to update the `images.tag` in the `charts/starship/values.yaml` of
[helm-charts](https://github.com/tricorder-observability/helm-charts)

Then bump `version charts/starship/Chart.yaml`, manually release helm
charts in
[Release Charts workflow](https://github.com/tricorder-observability/helm-charts/actions/workflows/release-chart.yaml)

You can also supply tag directly. You can use this to build version tag for testing.
But do not do this for release:
```
build_and_push_images.sh [tag]
# Then deploy this tag with helm
helm upgrade my-starship tricorder-stable/starship --set images.tag=[tag]
```
