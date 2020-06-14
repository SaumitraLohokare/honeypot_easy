package main

/*
	Honey Pot for a form entry.

	! Need to make a form with fake field
	! Check if Bot filled form data
	! Block bot IP
*/

import (
	"fmt"
	"net/http"
)

const html_code string = `
<!DOCTYPE html>
<html>
	<head>
		<title>Sweet Sign Up</title>
		<style>
		.sweet-input {
			opacity: 0;
		}
		</style>
	</head>
	<body>
		<form method="post">
			<input name="email" type="email" class="sweet-input">
			<input name="email-confirmation" type="email" class="email-field">
			<input name="password" type="password" class="password-field">
			<input name="submit" type="submit" class="submit-button">
		</form>
	</body>
</html>
`

var blocked_ips = []string{}

func find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func readUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

func sign_up(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		if _, found := find(blocked_ips, readUserIP(req)); found {
			w.Write([]byte("Access Denied"))
		} else {
			ip := readUserIP(req)
			req.ParseForm()
			for key, value := range req.Form {
				if key == "email" && value[0] != "" {
					blocked_ips = append(blocked_ips, ip)
					fmt.Printf("1 User was blocked :	ip = %s\n", ip)
				}
			}
			w.Write([]byte(html_code))
		}
	} else {
		if _, found := find(blocked_ips, readUserIP(req)); found {
			w.Write([]byte("Access Denied"))
		} else {
			w.Write([]byte(html_code))
		}
	}
}

func main() {
	http.HandleFunc("/", sign_up)
	http.ListenAndServe(":80", nil)
}
