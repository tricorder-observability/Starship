# GitHub

This file is not named README.md, because github would pick `.github/README.md`
over `README.md` at the root of the repo. See [docs](
https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-readmes#about-readmes).

GitHub actions workflow configuration files, scripts, and related contents.
* `workflows` defines workflow
* `scripts` includes scripts used during the github workflow execution
* `linters` stores individual linters configuration files for super linter
  **NOTE** All configuration files are hidden files `.<filename>`,
  us `ls -a` or `ls -A` to list.
* `super_linter.env` has the environment variables that enable/disable
  individual linters of super-linters. The file could be placed elsewhere,
  but the current path is a canonical location shared between local
  environment and GitHub action:
  [official doc](https://github.com/github/super-linter/blob/main/docs/run-linter-locally.md#sharing-environment-variables-between-local-and-ci).

