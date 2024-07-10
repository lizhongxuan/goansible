package work

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type ClientCmd struct {
	Url  string
	Ctx  context.Context
	Opts *WorkOptions
}

func (cc *ClientCmd) RunOutput(shell string, opts ...WorkOptionsFunc) (int, string, error) {
	for _, op := range opts {
		op(cc.Opts)
	}
	var resp RunOutputResp
	if err := post(cc.Ctx, cc.Url, WorkReq{
		Shell:        shell,
		TimeOut:      cc.Opts.TimeOut.Seconds(),
		OutPath:      cc.Opts.OutPath,
		ErrPath:      cc.Opts.ErrPath,
		Username:     cc.Opts.Username,
		SudoPassword: cc.Opts.SudoPassword,
		Stdin:        cc.Opts.Stdin,
	}, &resp); err != nil {
		return 0, "", err
	}
	if resp.Err == "" {
		return resp.StateCode, resp.Output, nil
	}
	return resp.StateCode, resp.Output, errors.New(resp.Err)
}

func (cc *ClientCmd) Start(shell string, opts ...WorkOptionsFunc) (int, error) {
	for _, op := range opts {
		op(cc.Opts)
	}
	var resp StartResp
	if err := post(cc.Ctx, cc.Url, WorkReq{
		Shell:        shell,
		TimeOut:      cc.Opts.TimeOut.Seconds(),
		OutPath:      cc.Opts.OutPath,
		ErrPath:      cc.Opts.ErrPath,
		Username:     cc.Opts.Username,
		SudoPassword: cc.Opts.SudoPassword,
		Stdin:        cc.Opts.Stdin,
	}, &resp); err != nil {
		return 0, err
	}
	if resp.Err == "" {
		return resp.Pid, nil
	}
	return resp.Pid, errors.New(resp.Err)
}

func post(ctx context.Context, url string, body interface{}, resqData interface{}, timeout ...time.Duration) error {
	t := 20 * time.Second
	if len(timeout) != 0 {
		t = timeout[0]
	}

	// 设置超时
	ctx, cancel := context.WithTimeout(ctx, t)
	defer cancel()

	req, err := buildRequest(ctx, url, body)
	if err != nil {
		return err
	}

	httpclient := &http.Client{}
	res, err := httpclient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errors.New("http statecode is invalid")
	}

	bodyData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(ctx, err)
		return err
	}
	if err := json.Unmarshal(bodyData, &resqData); err != nil {
		return err
	}
	return nil
}
func buildRequest(ctx context.Context, urlStr string, body interface{}) (*http.Request, error) {
	b, e := json.Marshal(body)
	if e != nil {
		return nil, e
	}
	reader := bytes.NewReader(b)
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	r, e := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), reader)
	if e != nil {
		return nil, e
	}
	r.Header.Set("Content-type", "application/json")
	return r, nil
}
