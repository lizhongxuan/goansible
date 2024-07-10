package work

import (
	"net/http"
	"time"
)
import "github.com/gin-gonic/gin"

type WorkReq struct {
	Shell        string  `json:"shell" binding:"required"`
	TimeOut      float64 `json:"time_out"`
	OutPath      string  `json:"out_path"`
	ErrPath      string  `json:"err_path"`
	Username     string  `json:"username"`
	SudoPassword string  `json:"sudo_password"`
	Stdin        string  `json:"stdin"`
}
type RunOutputResp struct {
	StateCode int    `json:"state_code"`
	Output    string `json:"output"`
	Err       string `json:"err"`
}
type StartResp struct {
	Pid int    `json:"pid"`
	Err string `json:"err"`
}

func SetupRouter(r *gin.Engine) {
	r.POST("/run_output", func(c *gin.Context) {
		var req WorkReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		cmd := &LocalCmd{}
		stateCode, out, err := cmd.RunOutput(req.Shell,
			WithStdin(req.Stdin),
			WithTimeOut(time.Duration(req.TimeOut)*time.Second),
			WithOutPath(req.OutPath),
			WithErrPath(req.ErrPath),
			WithUsername(req.Username),
			WithSudoPassword(req.SudoPassword),
		)
		c.JSON(http.StatusOK, RunOutputResp{
			stateCode, out, err.Error(),
		})
	})

	r.POST("/start", func(c *gin.Context) {
		var req WorkReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		cmd := &LocalCmd{}
		pid, err := cmd.Start(req.Shell,
			WithStdin(req.Stdin),
			WithTimeOut(time.Duration(req.TimeOut)*time.Second),
			WithOutPath(req.OutPath),
			WithErrPath(req.ErrPath),
			WithUsername(req.Username),
			WithSudoPassword(req.SudoPassword),
		)
		c.JSON(http.StatusOK, StartResp{
			pid, err.Error(),
		})
	})
}
