# A simple go program to update your Guild Wars 2 Arc dps

1. Read through this code and check to see if it's doing something stupid.
2. If you're happy with step 1, then clone the repository.

```bash
git clone https://github.com/srivatsa-scs/go-arcdps-updater
```

3. Make sure you have the necessary requirements such as go 1.16.3
4. Set up your configuration file. If you're using defaults, then rename the config.sample.json -> config.json
5. You can either run or build an exe using

```
$ go run .
```

or

```
$ go build .
```

6. Build option will give you an exe file (in windows) that you can use to run it.
7. The program will generate a log file and you can set the log level to debug mode to get additional information.
8. I've taken steps to ensure that if something goes wrong, things should revert back to before running the program, but if you still see any issues send me a message.
