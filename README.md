# taskfactory

---

taskfactory
A scheduler to be used as a cron replacement

The boilerplate for this project has been generated using [goboiler](https://github.com/shammishailaj/goboiler) which is a fork of [go-wagen](https://github.com/groovili/go-wagen).

The generated structure has been further modified as per the [Golang Standards Project Layout](https://github.com/golang-standards/project-layout) alongwith further customizations to modularize the code and add support for HTML templating and [scratch docker image](https://hub.docker.com/_/scratch).


### A. Direct Go Dependencies (3rd-party only)

---

1. [Cobra](https://github.com/spf13/cobra) for implementing the CLI interactions
2. [Go-Homedir](https://github.com/mitchellh/go-homedir) or detecting and expanding the user's home directory without cgo
3. [Go-Resty](https://github.com/go-resty/resty) as its HTTP client library
4. [Govvv](https://github.com/ahmetb/govvv) to add version information during its build process
5. [Logrus](https://github.com/sirupsen/logrus) as its logging library
6. [Dotsql](https://github.com/gchaincl/dotsql) for SQL migrations (not being used currently)
7. [Gopsutil/CPU](https://github.com/shirou/gopsutil/cpu) for CPU information (not being used currently)
8. [Gopsutil/Load](https://github.com/shirou/gopsutil/load) for system load information (not being used currently)
9. [Viper](https://github.com/spf13/viper) for reading configuration files
10. [Times](https://github.com/djherbis/times) for file times (atime, mtime, ctime, btime)
11. [Go-Chi](https://github.com/go-chi/chi) for HTTP routing
12. [Go-Chi/HttpRate](http://github.com/go-chi/httprate) for rate-limiting HTTP requests via [Go-Chi](https://github.com/go-chi/chi)
13. [Minify](https://github.com/tdewolff/minify) provides minifiers for web-formats - `CSS`, `HTML`, `JavaScript`, `JSON`, `SVG` and `XML`
14. [Gofast](https://github.com/yookoala/gofast) provides FastCGI "client" library written purely in go
15. [Fsnotify](https://github.com/fsnotify/fsnotify) - provides cross-platform file system notifications for Go.
16. [Cobra Docs](https://github.com/spf13/cobra/tree/master/doc) - Cobra documentation plugin


### B. Requirements

---

|Requirements Grid||
|---|---|
|Software|Version|
|Docker| &gt;= 19.03.13|
|---|---|
|Docker Compose| &gt;= 2.0.0|
|Ubuntu `build-essential` package||

The default app can be run/built using the following `make` targets.

#### C. `make` targets

---

`run` - Builds and runs the default app using the `docker` image `shammishailaj/gobuilder:0.0.3`

`build` - Builds the project using the `docker` image `shammishailaj/gobuilder:0.0.3` first and then copies all the required project files into a `scratch` image. The resultant is a docker image and not an executable.



#### D. Included Applications

---

- taskfactory CLI and Web server


#### How to run

```
go mod init github.com/shammishailaj/taskfactory
go mod tidy
go mod vendor
git init
git add .gitignore
git commit -m"Initial commit #1 Ignoring files"
git add .
git commit -m"Initial commit #2 Adding everything"
make run
```