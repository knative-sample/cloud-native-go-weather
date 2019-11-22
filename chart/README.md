# go-weather
Deploy the go weather

## Installing the Chart
To Install the weather chart:
    ```
    $ helm install chart --name weather --namespace default
    ```
## Uninstalling the Chart

To uninstall/delete the `weather` release but continue to track the release:
    ```
    $ helm delete weather
    ```

To uninstall/delete the `weather` release completely and make its name free for later use:
    ```
    $ helm delete weather --purge
    ```