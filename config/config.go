package config

import (
	"fmt"
	"github.com/decred/dcrd/dcrutil/v4"
	flags "github.com/jessevdk/go-flags"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	defaultConfigFileName = "dcrnsugar.conf"
	defaultLogLevel       = "info"
	defaultDataDirName    = "data"
	defaultLogDirName     = "logs"
)

var (
	dcrnsugarHomeDir  = dcrutil.AppDataDir("dcrnsugar", false)
	defaultConfigFile = filepath.Join(dcrnsugarHomeDir, defaultConfigFileName)
	defaultLogfile    = filepath.Join(dcrnsugarHomeDir, defaultLogDirName)
	defaultDatafile   = filepath.Join(dcrnsugarHomeDir, defaultDataDirName)

	defaultDBHost  = "127.0.0.1"
	defaultDBPort  = "5432"
	defaultDBUser  = "dcrdata"
	defaultDBPass  = "1234656"
	defaultDBName  = "dcrdata"
	defaultTimeOut = 200

	defaultServerHost = "127.0.0.1:7777"
	//defaultServerHost = "dcrdata.decred.org"
	defaultEnableSSL  = false
)

type Config struct {

	//General application behavior
	ConfigPath             string `short:"c" long:"Config" description:"Path to a custom Configuration file. (~/.dcrnsugar/dcrnsugar.conf)"`
	AppDirectory           string `long:"appdir" description:"Path to application home directory. (~/.dcrnsugar)"`
	DcrnsugarDataDirectory string `long:"dcrnsugardatadir" description:"Path to a dcrnsugar datadir"`
	LogPath                string `long:"logpath" description:"Directory to log output. ([appdir]/logs/)"`
	ShowVersion            bool   `short:"V" long:"version" description:"Display version information and exit"`

	// DB
	DBHost  string `long:"dbhost" description:"DB host"`
	DBPort  string `long:"dbport" description:"DB port"`
	DBUser  string `long:"dbuser" description:"DB user"`
	DBPass  string `long:"dbpass" description:"DB pass"`
	DBName  string `long:"dbname" description:"DB name"`
	TimeOut int    `long:"timeout" description:" timeout per work of db (second) "`

	// Client
	ServerHost string `long:"severHost" description:"explorer host"`
	EnableSSL  bool   `long:"enableSsl" description:"request https host"`
}

var defaultConfig = Config{
	AppDirectory:           dcrnsugarHomeDir,
	ConfigPath:             defaultConfigFile,
	DcrnsugarDataDirectory: defaultDatafile,
	LogPath:                defaultLogfile,

	DBHost:  defaultDBHost,
	DBPort:  defaultDBPort,
	DBUser:  defaultDBUser,
	DBName:  defaultDBName,
	DBPass:  defaultDBPass,
	TimeOut: defaultTimeOut,

	ServerHost: defaultServerHost,
	EnableSSL:  defaultEnableSSL,
}
var DefaultConfig = Config{
	AppDirectory:           dcrnsugarHomeDir,
	ConfigPath:             defaultConfigFile,
	DcrnsugarDataDirectory: defaultDatafile,
	LogPath:                defaultLogfile,

	DBHost:  defaultDBHost,
	DBPort:  defaultDBPort,
	DBUser:  defaultDBUser,
	DBName:  defaultDBName,
	DBPass:  defaultDBPass,
	TimeOut: defaultTimeOut,

	ServerHost: defaultServerHost,
	EnableSSL:  defaultEnableSSL,
}

func LoadConfig() (*Config, error) {
	cfg := defaultConfig

	preCfg := cfg
	preParser := flags.NewParser(&preCfg, flags.HelpFlag|flags.PassDoubleDash)
	_, err := preParser.Parse()

	if err != nil {
		e, ok := err.(*flags.Error)
		if !ok || e.Type != flags.ErrHelp {
			preParser.WriteHelp(os.Stderr)
		}
		if ok && e.Type == flags.ErrHelp {
			preParser.WriteHelp(os.Stdout)
			os.Exit(0)
		}
		return nil, err
	}

	// Show the version and exit if the version flag was specified.
	appName := filepath.Base(os.Args[0])
	appName = strings.TrimSuffix(appName, filepath.Ext(appName))
	if preCfg.ShowVersion {
		fmt.Printf("%s version %s (Go version %s, %s-%s)\n", appName,
			ver.String(), runtime.Version(), runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	parser := flags.NewParser(&cfg, flags.Default)

	if preCfg.AppDirectory != dcrnsugarHomeDir {
		preCfg.AppDirectory = cleanAndExpandPath(preCfg.AppDirectory)
	}

	defaultConfigPath := filepath.Join(preCfg.AppDirectory, defaultConfigFileName)
	if preCfg.ConfigPath == "" {
		preCfg.ConfigPath = defaultConfigPath
	} else {
		preCfg.ConfigPath = cleanAndExpandPath(preCfg.ConfigPath)
	}

	_, err = os.Stat(preCfg.ConfigPath)
	if os.IsNotExist(err) {
		if preCfg.ConfigPath != defaultConfigPath {
			fmt.Fprintln(os.Stderr, "No Configuration file found at "+preCfg.ConfigPath)
			os.Exit(0)
		}
	} else if err != nil {
		fmt.Fprintln(os.Stderr, "failed to stat Configuration file at "+preCfg.ConfigPath)
		return nil, err
	} else {
		err = flags.NewIniParser(parser).ParseFile(preCfg.ConfigPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			parser.WriteHelp(os.Stderr)
			return nil, fmt.Errorf("Unable to parse Configuration file.")
		}
	}

	_, err = parser.Parse()
	if err != nil {
		if e, ok := err.(*flags.Error); !ok || e.Type != flags.ErrHelp {
			parser.WriteHelp(os.Stderr)
		}
		return nil, fmt.Errorf("Error parsing command line arguments: %v", err)
	}

	if cfg.AppDirectory != preCfg.AppDirectory {
		cfg.AppDirectory = cleanAndExpandPath(cfg.AppDirectory)
		err = os.MkdirAll(cfg.AppDirectory, 0700)
		if err != nil {
			return nil, fmt.Errorf("Application directory error: %v", err)
		}
	}

	if cfg.LogPath == "" {
		cfg.LogPath = filepath.Join(cfg.AppDirectory, defaultLogfile)
	} else {
		cfg.LogPath = cleanAndExpandPath(cfg.LogPath)
	}

	if cfg.DcrnsugarDataDirectory == "" {
		cfg.DcrnsugarDataDirectory = defaultDatafile
	}
	cfg.DcrnsugarDataDirectory = cleanAndExpandPath(cfg.DcrnsugarDataDirectory)

	return &cfg, nil
}

// cleanAndExpandPath expands environment variables and leading ~ in the passed
// path, cleans the result, and returns it.
func cleanAndExpandPath(path string) string {
	// NOTE: The os.ExpandEnv doesn't work with Windows cmd.exe-style
	// %VARIABLE%, but the variables can still be expanded via POSIX-style
	// $VARIABLE.
	path = os.ExpandEnv(path)

	if !strings.HasPrefix(path, "~") {
		return filepath.Clean(path)
	}

	// Expand initial ~ to the current user's home directory, or ~otheruser to
	// otheruser's home directory.  On Windows, both forward and backward
	// slashes can be used.
	path = path[1:]

	var pathSeparators string
	if runtime.GOOS == "windows" {
		pathSeparators = string(os.PathSeparator) + "/"
	} else {
		pathSeparators = string(os.PathSeparator)
	}

	userName := ""
	if i := strings.IndexAny(path, pathSeparators); i != -1 {
		userName = path[:i]
		path = path[i:]
	}

	homeDir := ""
	var u *user.User
	var err error
	if userName == "" {
		u, err = user.Current()
	} else {
		u, err = user.Lookup(userName)
	}
	if err == nil {
		homeDir = u.HomeDir
	}
	// Fallback to CWD if user lookup fails or user has no home directory.
	if homeDir == "" {
		homeDir = "."
	}

	return filepath.Join(homeDir, path)
}
