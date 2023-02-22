# Release

To create a new release:

* Build and push a new version with `tools/scripts/build_and_release_images.sh`
* Go to [helm-charts actions tab](https://github.com/tricorder-observability/helm-charts/actions),
  click the `Release Charts` tab on the left panel, then `Run workflow` on the right side, choose the right tag from
  the `Tags` dropdown menu, and at last click `Run workflow` green button:
  ![image](https://user-images.githubusercontent.com/112656580/220518694-d6a2cee2-0352-4dad-8c80-171a30e525f4.png)
* After the workflow finishes, take a look at the
  [helm-charts releases](https://github.com/tricorder-observability/helm-charts/releases),
  verify with the new release is present, named after the desired chart version.
