package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strings"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"

	"github.com/vieux/gocover.io/server/redis"
)

var (
	socket    = flag.String("s", "/var/run/docker.sock", "Dockerd socket")
	serveAddr = flag.String("p", ":8080", "Address and port to serve")
	redisAddr = flag.String("r", "127.0.0.1:6379", "redis address")
)

func main() {
	flag.Parse()

	var (
		pool = redis.NewPool("tcp", *redisAddr)
		m    = martini.Classic()
	)

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
	m.Get("/_score/**", func(params martini.Params) (int, string) {
		var (
			repo = params["_1"]
			conn = pool.Get()
		)

		defer conn.Close()
		if coverage, err := redis.GetCoverage(conn, repo); err != nil {
			return 500, err.Error()
		} else {
			return 200, coverage
		}

	})
	m.Get("/_cache/**", func(params martini.Params) (int, string) {
		var (
			repo = params["_1"]
			conn = pool.Get()
		)
		defer conn.Close()

		if cached, _, err := redis.GetRepo(conn, repo); err != nil {
			return 500, err.Error()
		} else if cached != "" {
			redis.SetStats(conn, repo)
			return 200, string(cached)
		}
		return 404, "No cached version of " + repo
	})
	m.Get("/_/**", func(params martini.Params) (int, string) {
		var (
			repo = params["_1"]
			conn = pool.Get()
		)
		defer conn.Close()

		if cached, fresh, err := redis.GetRepo(conn, repo); err != nil {
			return 500, err.Error()
		} else if fresh {
			return 200, string(cached)
		}

		out, err := exec.Command("docker", "-H", "unix://"+*socket, "run", "--rm", "-a", "stdout", "-a", "stderr", "worker", repo).CombinedOutput()
		if err != nil {
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

		/*
			re, err = regexp.Compile("\\>" + repo + "([\\S\\s]?)\\</option\\>")
			if err != nil {
				return 500, err.Error()
			}
			content = re.ReplaceAllString(result.String(), "\\>...$1\\</option\\>")
		*/

		re = regexp.MustCompile("-- cov:([0-9.]*) --")
		matches := re.FindStringSubmatch(content)
		if len(matches) == 2 {
			redis.SetCache(conn, repo, content, matches[1])
		} else {
			redis.SetCache(conn, repo, content, "-1")
		}
		redis.SetStats(conn, repo)
		return 200, content
	})
	m.Get("/**", func(params martini.Params, r render.Render) {
		var (
			repo = params["_1"]
			conn = pool.Get()
		)
		defer conn.Close()

		if cached, fresh, err := redis.GetRepo(conn, repo); err != nil {
			r.HTML(500, "", map[string]interface{}{"cover_active": "active", "error": err})
		} else if fresh {
			redis.SetStats(conn, repo)
			r.HTML(200, "cached", map[string]interface{}{"repo": repo, "cover_active": "active", "cache": template.HTML(cached)})
		} else {
			contexts := map[string]interface{}{"repo": repo, "cover_active": "active"}
			if cached != "" {
				contexts["cache"] = "ok"
			}
			r.HTML(200, "loading", contexts)
		}
	})
	log.Fatal(http.ListenAndServe(*serveAddr, m))
}
