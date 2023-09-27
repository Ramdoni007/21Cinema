package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/Ramdoni007/21Cinema/internal/data"
	"github.com/Ramdoni007/21Cinema/internal/jsonlog"
)

const version = "1.0.0"

// Add a db struct field to hold the configuration settings for our database connection
// pool. For now this only holds the DSN, which we will read in from a command-line flag.
// Add maxOpenConns, maxIdleConns and maxIdleTime fields to hold the configuration
// settings for the connection pool.
// Add a models field to hold our new Models struct.
type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenCoons int
		maxIdleCoons int
		maxIdleTime  string
	}
	// Add a new limiter struct containing fields for the requests-per-second and burst
	// values, and a boolean field which we can use to enable/disable rate limiting
	// altogether.
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

// Change the logger field to have the type *jsonlog.Logger, instead of
// *log.Logger.
type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
}

func main() {
	var cfg config

	// Read the value of the port and env command-line flags into the config struct. We
	// default to using the port number 4000 and the environment "development" if no
	// corresponding flags are provided.
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	// Read the DSN value from the db-dsn command-line flag into the config struct. We
	// default to using our development DSN if no flag is provided.
	flag.StringVar(
		&cfg.db.dsn,
		"db-dsn",
		// os.Getenv("TWENTYONE"),
		"postgres://twentyone:secret@localhost/twenty_one?sslmode=disable",
		"Postgresql DSN",
	)

	// Read the connection pool settings from command-line flags into the config struct.
	// Notice the default values that we're using?
	flag.IntVar(&cfg.db.maxOpenCoons, "db-max-open-conns", 25, "Postgresql max open connections")
	flag.IntVar(&cfg.db.maxIdleCoons, "db-max-idle-conns", 25, "Postgresql max idle connections")
	flag.StringVar(
		&cfg.db.maxIdleTime,
		"db-max-idle-Time",
		"15m",
		"Postgresql max connections idle time",
	)

	// Create command line flags to read the setting values into the config struct.
	// Notice that we use true as the default for the 'enabled' setting?
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum request persecond")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enabled rate limiter")

	flag.Parse()

	// Initialize a new jsonlog.Logger which writes any messages *at or above* the INFO
	// severity level to the standard out stream.
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	// Call the openDB() helper function (see below) to create the connection pool,
	// passing in the config struct. If this returns an error, we log it and exit the
	// application immediately.
	db, err := openDB(cfg)
	if err != nil {
		// Use the PrintFatal() method to write a log entry containing the error at the
		// FATAL level and exit. We have no additional properties to include in the log
		// entry, so we pass nil as the second parameter
		logger.PrintFatal(err, nil)
	}

	// Defer a call to db.Close() so that the connection pool is closed before the
	// main() function exits.
	defer db.Close()
	// Likewise use the PrintInfo() method to write a message at the INFO level.
	logger.PrintInfo("database connection pool established", nil)

	// Declare an instance of the application struct, containing the config struct and
	// the logger
	// Use the data.NewModels() function to initialize a Models struct, passing in the
	// connection pool as a paramete
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModel(db),
	}

	//// Declare  HTTP server with some sensible timeout settings, which listens on the
	//// port provided in the config struct and uses the server we created above as the
	//// handler. use the http router instance returned by app.router() as the server handler
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		ErrorLog:     log.New(logger, "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Again, we use the PrintInfo() method to write a "starting server" message at the
	// INFO level. But this time we pass a map containing additional properties (the
	// operating environment and server address) as the final parameter.
	logger.PrintInfo("Server successfully Start", map[string]string{
		"addr": srv.Addr,
		"env":  cfg.env,
	})

	err = srv.ListenAndServe()
	// Use the PrintFatal() method to log the error and exit.
	logger.PrintFatal(err, nil)
}

// The openDB() function returns a sql.DB connection pool.
func openDB(cfg config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config
	// struct.
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	// Set the maximum number of open (in-use + idle) connections in the pool. Note that
	// passing a value less than or equal to 0 will mean there is no limit.
	db.SetMaxOpenConns(cfg.db.maxOpenCoons)

	// Set the maximum number of idle connections in the pool. Again, passing a value
	// less than or equal to 0 will mean there is no limit.
	db.SetMaxIdleConns(cfg.db.maxIdleCoons)

	// Use the time.ParseDuration() function to convert the idle timeout duration string
	// to a time.Duration type.
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	// Set Maximum idle timeout
	db.SetConnMaxIdleTime(duration)

	// Create a context with a 5-second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use PingContext() to establish a new connection to the database, passing in the
	// context we created above as a parameter. If the connection couldn't be
	// established successfully within the 5 second deadline, then this will return an
	// error.
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	// Return the sql.DB connection pool.
	return db, nil
}
