package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	r "github.com/garyburd/redigo/redis"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"

	"github.com/vieux/gocover.io/server/redis"
)

var (
	docker_socket = flag.String("s", "", "Dockerd socket (e.g., /var/run/docker.sock)")
	docker_addr   = flag.String("d", "", "Dockerd addr (e.g., 127.0.0.1:2375)")
	serveAddr     = flag.String("p", ":8080", "Address and port to serve")
	serveSAddr    = flag.String("ps", ":80443", "Address and port to serve HTTPS")
	redisAddr     = flag.String("r", "127.0.0.1:6379", "redis address")
	redisPass     = flag.String("rp", "", "redis password")
	certPath      = flag.String("tls", "", "cert path")
)

func docker(repo, version string, pool *r.Pool) (int, string) {
	var (
		worker = "vieux/gocover"
		conn   = pool.Get()
	)

	defer conn.Close()

	if version != "" {
		worker = worker + ":" + version
	}

	if version == "" {
		if cached, fresh, err := redis.GetRepo(conn, repo); err != nil {
			return 500, err.Error()
		} else if fresh {
			return 200, string(cached)
		}
	}

	host := ""

	if *docker_socket != "" {
		host = "unix://" + *docker_socket
	} else if *docker_addr != "" {
		host = "tcp://" + *docker_addr
	} else {
		return 500, "cannot connect to docker daemon"
	}

	out, err := exec.Command("docker", "-H", host, "run", "--rm", "-a", "stdout", "-a", "stderr", worker, repo).CombinedOutput()
	if err != nil {
		if strings.Contains(string(out), "Unable to find image") {
			return 500, "go version '" + version + "' not found"
		}
		return 500, string(out)
	}
	re, err := regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	if err != nil {
		return 500, err.Error()
	}
	content := re.ReplaceAllString(string(out), "")
	content = strings.Replace(content, "background: black;", "background: #222222;", 2)

	content = strings.Replace(content, ".cov1 { color: rgb(128, 128, 128) }", ".cov1 { color: #52987D }", 2)
	content = strings.Replace(content, ".cov2 { color: rgb(128, 128, 128) }", ".cov2 { color: #4BA180 }", 2)
	content = strings.Replace(content, ".cov3 { color: rgb(128, 128, 128) }", ".cov3 { color: #44AA83 }", 2)
	content = strings.Replace(content, ".cov4 { color: rgb(128, 128, 128) }", ".cov4 { color: #3DB487 }", 2)
	content = strings.Replace(content, ".cov5 { color: rgb(128, 128, 128) }", ".cov5 { color: #36BD8A }", 2)
	content = strings.Replace(content, ".cov6 { color: rgb(128, 128, 128) }", ".cov6 { color: #2FC68D }", 2)
	content = strings.Replace(content, ".cov7 { color: rgb(128, 128, 128) }", ".cov7 { color: #28D091 }", 2)
	content = strings.Replace(content, ".cov8 { color: rgb(128, 128, 128) }", ".cov8 { color: #21D994 }", 2)
	content = strings.Replace(content, ".cov9 { color: rgb(128, 128, 128) }", ".cov9 { color: #1AE297 }", 2)
	content = strings.Replace(content, "<option value=\"file0\">", "<option value=\"file0\" select=\"selected\">", -1)
	content = strings.Replace(content, "\">"+repo, "\">", -1)

	re = regexp.MustCompile("-- cov:([0-9.]*) --")
	matches := re.FindStringSubmatch(content)
	if len(matches) == 2 {
		cov, err := strconv.ParseFloat(matches[1], 64)
		if err == nil {
			content = strings.Replace(content, "<select id=", fmt.Sprintf("<span class='cov%d'>%s%%</span> | <select id=", int((cov-0.0001)/10), matches[1]), 1)
		}
		if version != "" {
			content = strings.Replace(content, "<select id=", fmt.Sprintf("<span>%s</span> | <select id=", version), 1)
		} else {
			redis.SetCache(conn, repo, content, matches[1])
		}
	} else if version != "" {
		content = strings.Replace(content, "<select id=", fmt.Sprintf("<span>%s</span> | <select id=", version), 1)
	} else {
		redis.SetCache(conn, repo, content, "-1")
	}
	if version == "" {
		redis.SetStats(conn, repo)
	}
	return 200, content
}

func main() {
	flag.Parse()

	var (
		m         = martini.Classic()
		pool, err = redis.NewPool("tcp", *redisAddr, *redisPass)
	)

	if err != nil {
		log.Fatalf("%v", err)
	}

	m.Use(martini.Static("static"))
	m.Use(render.Renderer(render.Options{Layout: "layout"}))
	m.Get("/about", func(r render.Render) {
		r.HTML(200, "about", map[string]interface{}{"about_active": "active"})
	})
	m.Get("/", func(r render.Render) {
		conn := pool.Get()
		defer conn.Close()

		top, err := redis.Top(conn, "top", 4)
		if err != nil {
			log.Println(err.Error())
		}
		last, err := redis.Top(conn, "last", 4)
		if err != nil {
			log.Println(err.Error())
		}
		r.HTML(200, "cover", map[string]interface{}{"top": top, "last": last, "cover_active": "active"})
	})
	m.Post("/_webhook", func(req *http.Request) (int, string) {

		github := struct {
			Repository struct {
				Full_Name string
			}
		}{}

		if err := json.NewDecoder(req.Body).Decode(&github); err != nil {
			return 500, err.Error()
		}

		go docker(github.Repository.Full_Name, "", pool)

		return 202, github.Repository.Full_Name
	})
	m.Get("/_badge/**", func(params martini.Params, r render.Render) {
		var (
			repo = params["_1"]
			conn = pool.Get()
		)

		defer conn.Close()
		if coverage, err := redis.GetCoverage(conn, repo); err != nil {
			r.Redirect(fmt.Sprintf("https://img.shields.io/badge/coverage-error-lightgrey.svg?style=flat"))
		} else if coverage < 25.0 {
			r.Redirect(fmt.Sprintf("https://img.shields.io/badge/coverage-%.1f%%25-red.svg?style=flat", coverage))
		} else if coverage < 50.0 {
			r.Redirect(fmt.Sprintf("https://img.shields.io/badge/coverage-%.1f%%25-orange.svg?style=flat", coverage))
		} else if coverage < 75.0 {
			r.Redirect(fmt.Sprintf("https://img.shields.io/badge/coverage-%.1f%%25-green.svg?style=flat", coverage))
		} else {
			r.Redirect(fmt.Sprintf("https://img.shields.io/badge/coverage-%.1f%%25-brightgreen.svg?style=flat", coverage))
		}

	})
	m.Get("/_cache/**", func(req *http.Request, params martini.Params) (int, string) {
		var (
			repo = params["_1"]
			conn = pool.Get()
		)
		defer conn.Close()

		if req.FormValue("version") == "" {
			if cached, _, err := redis.GetRepo(conn, repo); err != nil {
				return 500, err.Error()
			} else if cached != "" {
				redis.SetStats(conn, repo)
				return 200, string(cached)
			}
		}
		return 404, "No cached version of " + repo
	})
	m.Get("/_/**", func(req *http.Request, params martini.Params) (int, string) {
		return docker(params["_1"], req.FormValue("version"), pool)
	})
	m.Get("/**", func(req *http.Request, params martini.Params, r render.Render) {
		var (
			repo    = params["_1"]
			conn    = pool.Get()
			version = ""
		)
		defer conn.Close()

		if req.ParseForm() == nil {
			version = req.FormValue("version")
		}

		if cached, fresh, err := redis.GetRepo(conn, repo+version); err != nil {
			r.HTML(500, "", map[string]interface{}{"cover_active": "active", "error": err})
		} else if fresh {
			redis.SetStats(conn, repo)
			r.HTML(200, "cached", map[string]interface{}{"repo": repo, "cover_active": "active", "cache": template.HTML(cached), "version": version})
		} else {
			contexts := map[string]interface{}{"repo": repo, "cover_active": "active", "version": version}
			if cached != "" {
				contexts["cache"] = "ok"
			}
			r.HTML(200, "loading", contexts)
		}
	})
	if *certPath != "" {
		go func() {
			log.Println(http.ListenAndServe(*serveAddr, http.HandlerFunc(redir)))
		}()
		log.Fatal(http.ListenAndServeTLS(*serveSAddr, filepath.Join(*certPath, "fullchain.pem"), filepath.Join(*certPath, "privkey.pem"), m))
	} else {
		log.Fatal(http.ListenAndServe(*serveAddr, m))
	}
}

func redir(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "https://gocover.io"+req.RequestURI, http.StatusMovedPermanently)
}
