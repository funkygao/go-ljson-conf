package ljconf

import (
	"os"
	"os/user"

	"github.com/daviddengcn/go-villa"
	"github.com/daviddengcn/ljson"
)

const INCLUDE_KEY_TAG = "#include#"

// loadFromFile reads configurations from a speicified file. If some error found
// during reading, it will be return, but the conf is still available.
func loadFromFile(fn string) (conf *Conf, err error) {
	if _, err := os.Stat(fn); err != nil {
		return nil, err
	}

	path := findPath(villa.Path(fn))
	conf = &Conf{
		path: path,
		db:   make(map[string]interface{}),
	}

	fin, err := path.Open()
	if err != nil {
		if os.IsNotExist(err) {
			// configuration file not existing is ok, an empty conf
			return conf, nil
		}
		// if file not exists, nothing read (but configuration still usable.)
		return conf, err
	}

	if err := func() error {
		defer fin.Close()

		dec := ljson.NewDecoder(newRcReader(fin))
		return dec.Decode(&conf.db)
	}(); err != nil {
		return conf, err
	}

	loadInclude(conf.db, path.Dir())

	return conf, nil
}

func loadArrayInclude(arr []interface{}, dir villa.Path) {
	for _, el := range arr {
		switch vv := el.(type) {
		case map[string]interface{}:
			loadInclude(vv, dir)
		case []interface{}:
			loadArrayInclude(vv, dir)
		}
	}
}

func loadInclude(db map[string]interface{}, dir villa.Path) {
	for k, v := range db {
		if k == INCLUDE_KEY_TAG {
			switch paths := v.(type) {
			case string:
				//				fmt.Println("Including", paths, "at", dir)
				sub, err := loadFromFile(dir.Join(paths).S())
				if err == nil {
					// merge into current db
					for sk, sv := range sub.db {
						db[sk] = sv
					}
					// remove this entry
					delete(db, k)
				}
				continue
			case []interface{}:
				for _, el := range paths {
					if path, ok := el.(string); ok {
						sub, err := loadFromFile(dir.Join(path).S())
						if err == nil {
							// merge into current db
							for sk, sv := range sub.db {
								db[sk] = sv
							}
						}
					}
				}
				// remove this entry
				delete(db, k)
				continue
			} // switch
		} // if

		switch vv := v.(type) {
		case map[string]interface{}:
			loadInclude(vv, dir)
		case []interface{}:
			loadArrayInclude(vv, dir)
		}
	}
}

func findPath(fn villa.Path) villa.Path {
	if fn.IsAbs() {
		return fn
	}

	if fn.Exists() {
		return fn.AbsPath()
	}

	// Try .exe folder
	tryFn := villa.Path(os.Args[0]).Dir().Join(fn)
	if tryFn.Exists() {
		return tryFn
	}

	// Try user-home folder
	cu, err := user.Current()
	if err == nil {
		tryFn = villa.Path(cu.HomeDir).Join(fn)
		if tryFn.Exists() {
			return tryFn
		}
	}
	return fn.AbsPath()
}
