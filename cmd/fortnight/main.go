package main // import "github.com/motki/fortnight/cmd/fortnight"

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/motki/core/app"
	"github.com/motki/core/app/profile"
	"github.com/motki/core/log"
	"github.com/motki/core/proto"
	"github.com/motki/core/proto/client"

	"github.com/motki/fortnight/localhttp"
	"github.com/motki/fortnight/localstore"
)

var serverAddr = flag.String("server", "motki.org:18443", "Backend server host and port.")
var credentials = flag.String("credentials", "", "Username and password separated by a colon. (ie. \"frank:mypass\")")
var assetsDir = flag.String("assets", "", "Path to static assets.")
var logLevel = flag.String("log-level", "warn", "Log level. Possible values: debug, info, warn, error.")
var insecureSkipVerify = flag.Bool("insecure-skip-verify", false, "INSECURE: Skip verification of server SSL cert.")
var version = flag.Bool("version", false, "Display the application version.")

// fatalf creates a default logger, writes the given message, and exits.
func fatalf(format string, vals ...interface{}) {
	logger := log.New(log.Config{})
	logger.Fatalf(format, vals...)
}

// Config wraps an app.Config and contains optional credentials.
type Config struct {
	*app.Config

	username string
	password string
}

// WithCredentials returns a copy of the Config with the given credentials embedded.
func (c Config) WithCredentials(username, password string) Config {
	return Config{c.Config, username, password}
}

// NewEnv initializes a new CLI environment.
//
// If the given Config contains a username or password, authentication
// will be attempted. If authentication fails, an error is returned.
func NewEnv(conf Config) (*app.ClientEnv, error) {
	env, err := app.NewClientEnv(conf.Config)
	if err != nil {
		return nil, err
	}
	if conf.username != "" || conf.password != "" {
		if err = env.Client.Authenticate(conf.username, conf.password); err != nil {
			return nil, err
		} else {
			env.Logger.Debugf("authenticated as %s", conf.username)
		}
	} else {
		env.Logger.Debugf("running without authentication")
	}
	return env, nil
}

var Version = "dev"

func main() {
	flag.Parse()
	if *version {
		fmt.Printf("%s %s. %s\n", os.Args[0], Version, "https://github.com/motki/fortnight")
		os.Exit(0)
	}

	if *assetsDir == "" {
		wd, err := os.Getwd()
		if err != nil {
			fatalf("fortnight: error getting current working directory: %s", err.Error())
		}
		*assetsDir = wd + "/assets"
	}
	if _, err := os.Stat(*assetsDir); err != nil {
		fatalf("fortnight: could not locate static assets: %s", err.Error())
	}

	// Start the profiler, if enabled via command-line flags.
	pr := profile.New()

	// Trimmed down motki app configuration.
	appConf := &app.Config{
		Backend: proto.Config{
			Kind: proto.BackendRemoteGRPC,
			RemoteGRPC: proto.RemoteConfig{
				ServerAddr:         *serverAddr,
				InsecureSkipVerify: *insecureSkipVerify,
			},
		},
		Logging: log.Config{
			Level: *logLevel,
		},
	}

	// Writing to stderr offers a way to redirect the logger output to a file instead of
	// interrupting the user.
	appConf.Logging.OutputType = log.OutputStderr

	conf := Config{Config: appConf}

	// Add credentials to the configuration, if they were specified.
	if *credentials != "" {
		parts := strings.Split(*credentials, ":")
		if len(parts) != 2 {
			fatalf("fortnight: invalid credentials, expected format \"username:password\"")
		}
		conf = conf.WithCredentials(parts[0], parts[1])
	} else {
		user, pass := os.Getenv("MOTKI_USERNAME"), os.Getenv("MOTKI_PASSWORD")
		if user != "" && pass != "" {
			conf = conf.WithCredentials(user, pass)
		}
	}

	// Initialize the environment.
	env, err := NewEnv(conf)
	if err != nil {
		if err == client.ErrBadCredentials {
			fmt.Println("Invalid username or password.")
		}
		fatalf("fortnight: error initializing application environment: %s", err.Error())
	}

	store, err := localstore.New("data")
	if err != nil {
		fatalf("fortnight: error initialize localstore: %s", err.Error())
	}

	srv := localhttp.NewServer(env.Client, env.Logger, store, *assetsDir)
	go func() {
		err := srv.Serve()
		if err != nil {
			env.Logger.Warnf("fortnight: error returned from server: %s", err.Error())
		}
	}()

	// Block the main loop until SIGINT, SIGHUP, or SIGKILL is received,
	// or until CTRL+C is pressed from the command-line.
	err = env.BlockUntilSignal(make(chan os.Signal))

	// Stop the profiler, if enabled.
	if pr != nil {
		pr.Stop()
	}

	if err != nil {
		env.Logger.Warnf("motki: failed to shutdown cleanly: %s", err.Error())
		os.Exit(1)
	}
}
