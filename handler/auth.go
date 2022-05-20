package fileupload

import "net/http"

//修了它pipline 不做洋葱
func HTTPInterceptor(funcHandler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func (w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			userName := r.Form.Get("userName")
			token := r.Form.Get("token")

			if len(userName) < 3 || !IsTokenVaild(token) {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			funcHandler(w,r)
		})
}
