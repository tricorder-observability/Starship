# Release

To create a new release:

* Build and push a new version with `tools/scripts/build_and_release_images.sh`
* Go to [helm-charts actions tab](https://github.com/tricorder-observability/helm-charts/actions),
  click the `Release Charts` tab on the left panel, then `Run workflow` on the right side, choose the right tag from
  the `Tags` dropdown menu, and at last click `Run workflow` green button:
  ![image](https://user-images.githubusercontent.com/112656580/213858606-0ecdecad-2a31-4117-8499-1ba60af3b076.png)
* After the workflow finishes, take a look at the [helm-charts github page](https://tricorder-observability.github.io/helm-charts/),
  and verify with the new release with helm.
