package main

import (
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"regexp"
)

const (
	Us       = "{Us}"
	Root     = "{Root}"
	Addr     = "{Addr}"
	Protocol = "{Protocol}"
	Domain   = "{Domain}"
)

func main() {
	fmt.Println(Us)
	fmt.Println(Root)
	fmt.Println(Addr)
	if err := http.ListenAndServe(Addr, server()); err != nil {
		panic(err)
	}
}

type _server struct{}

func server() *_server {
	return &_server{}
}
func (this *_server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(fmt.Sprintf("%v", err)))
		}
	}()
	{
		matched := regexp.MustCompile(`^\/([\w-]+)\/([\w-]+)(\.git)?$`).FindStringSubmatch(r.URL.Path)
		l := len(matched)
		if l > 0 {
			mod(w, r, matched, l).do()
			return
		}
	}
	var us string
	_, us, b := r.BasicAuth()
	if !b {
		us = r.Header.Get("us")
		if len(us) == 0 {
			usCookie, err := r.Cookie("us")
			if err == nil {
				us = usCookie.Value
			}
		}
	}
	fmt.Println(us)
	if us != Us {
		w.Header().Set("WWW-Authenticate", "Basic realm=")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	{
		matched := regexp.MustCompile(`^\/([\w-]+)\/([\w-]+)(\.git)?\/info\/refs$`).FindStringSubmatch(r.URL.Path)
		l := len(matched)
		if l > 0 {
			ref(w, r, matched, l).do()
			return
		}
	}
	{
		matched := regexp.MustCompile(`^\/([\w-]+)\/([\w-]+)(\.git)?\/(git-upload-pack|git-receive-pack)$`).FindStringSubmatch(r.URL.Path)
		l := len(matched)
		if l > 0 {
			rpc(w, r, matched, l).do()
			return
		}
	}
}

type _mod struct {
	w       http.ResponseWriter
	r       *http.Request
	matched []string
	l       int
}

func mod(w http.ResponseWriter, r *http.Request, matched []string, l int) *_mod {
	return &_mod{
		w:       w,
		r:       r,
		matched: matched,
		l:       l,
	}
}
func (this *_mod) do() {
	fmt.Println("ref")
	user := this.matched[1]
	fmt.Println(user)
	repo := this.matched[2]
	fmt.Println(repo)
	modName := fmt.Sprintf("%s/%s/%s", Domain, user, repo)
	fmt.Println(modName)
	modUrl := fmt.Sprintf("%s://%s", Protocol, modName)
	fmt.Println(modUrl)
	modHtml := fmt.Sprintf(`<!DOCTYPE html>
<html>
    <head>
        <meta name="go-import" content="%s git %s.git">
    </head>
    <body></body>
</html>`, modName, modUrl)
	fmt.Println(modHtml)
	this.w.Header().Set("Content-Type", "text/html")
	this.w.Header().Set("Expires", "Fri, 01 Jan 1980 00:00:00 GMT")
	this.w.Header().Set("Pragma", "no-cache")
	this.w.Header().Set("Cache-Control", "no-cache, max-age=0, must-revalidate")
	if _, err := this.w.Write([]byte(modHtml)); err != nil {
		fmt.Println(err)
		panic(err)
	}
}

type _ref struct {
	w       http.ResponseWriter
	r       *http.Request
	matched []string
	l       int
}

func ref(w http.ResponseWriter, r *http.Request, matched []string, l int) *_ref {
	return &_ref{
		w:       w,
		r:       r,
		matched: matched,
		l:       l,
	}
}
func (this *_ref) do() {
	fmt.Println("ref")
	user := this.matched[1]
	fmt.Println(user)
	repo := this.matched[2]
	fmt.Println(repo)
	path := fmt.Sprintf("%s/%s/%s.git", Root, user, repo)
	fmt.Println(path)
	service := this.r.URL.Query().Get("service")
	fmt.Println(service)
	this.w.Header().Set("Content-Type", fmt.Sprintf("application/x-%s-advertisement", service))
	this.w.Header().Set("Expires", "Fri, 01 Jan 1980 00:00:00 GMT")
	this.w.Header().Set("Pragma", "no-cache")
	this.w.Header().Set("Cache-Control", "no-cache, max-age=0, must-revalidate")
	pklText := fmt.Sprintf("# service=%s\n", service)
	pklTextLength16 := fmt.Sprintf("%04x", len(pklText)+4)
	if _, err := this.w.Write([]byte(pklTextLength16)); err != nil {
		fmt.Println(err)
		panic(err)
	}
	if _, err := this.w.Write([]byte(pklText)); err != nil {
		fmt.Println(err)
		panic(err)
	}
	if _, err := this.w.Write([]byte("0000")); err != nil {
		fmt.Println(err)
		panic(err)
	}
	cmd := exec.Command("git", service[4:], "--stateless-rpc", "--advertise-refs", path)
	b, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	if _, err := this.w.Write(b); err != nil {
		fmt.Println(err)
		panic(err)
	}
}

type _rpc struct {
	w       http.ResponseWriter
	r       *http.Request
	matched []string
	l       int
}

func rpc(w http.ResponseWriter, r *http.Request, matched []string, l int) *_rpc {
	return &_rpc{
		w:       w,
		r:       r,
		matched: matched,
		l:       l,
	}
}
func (this *_rpc) do() {
	fmt.Println("rpc")
	user := this.matched[1]
	fmt.Println(user)
	repo := this.matched[2]
	fmt.Println(repo)
	path := fmt.Sprintf("%s/%s/%s.git", Root, user, repo)
	fmt.Println(path)
	service := this.matched[this.l-1]
	fmt.Println(service)
	this.w.Header().Set("Content-Type", fmt.Sprintf("application/x-%s-result", service))
	this.w.Header().Set("Connection", "Keep-Alive")
	this.w.Header().Set("Transfer-Encoding", "chunked")
	this.w.Header().Set("X-Content-Type-Options", "nosniff")
	cmd := exec.Command("git", service[4:], "--stateless-rpc", path)
	cmdStdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	cmdStdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	if err := cmd.Start(); err != nil {
		fmt.Println(err)
		panic(err)
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
				this.w.WriteHeader(http.StatusServiceUnavailable)
				this.w.Write([]byte(fmt.Sprintf("%v", err)))
			}
		}()
		if _, err := io.Copy(cmdStdin, this.r.Body); err != nil {
			fmt.Println(err)
			panic(err)
		}
		if err := cmdStdin.Close(); err != nil {
			fmt.Println(err)
			panic(err)
		}
	}()
	if _, err := io.Copy(this.w, cmdStdout); err != nil {
		fmt.Println(err)
		panic(err)
	}
	if err := cmd.Wait(); err != nil {
		fmt.Println(err)
		panic(err)
	}
}
