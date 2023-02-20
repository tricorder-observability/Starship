# License

Starship use [skywalking-eyes](https://github.com/apache/skywalking-eyes) to 
check and fix License header.

License header check and fix rule `.licenserc.yaml` defined in the root of starship.

And licnese content follows [AGPL3 license header](https://github.com/licenses/license-templates/blob/master/templates/agpl3-header.txt)

# How to add License
When you create new file in `src/` directory, we need to add and check License header before your push.

```shell
cd starship
make addlicense
```

> Note: skywalking-eyes will add license directly and does not support override License header so far.

# How to check license

```shell
cd starship
make checklicense
```

# License check Pipline
starship use Github actions as License Check Pipeline defined in `.github/license-check.yaml`.
