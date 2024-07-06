package work

import (
	"net/http"
	"time"
)
import "github.com/gin-gonic/gin"

func SetupRouter(r *gin.Engine) {
	r.POST("/run_output", func(c *gin.Context) {
		var Res struct {
			Shell        string `json:"shell" binding:"required"`
			TimeOut      int64  `json:"time_out"`
			OutPath      string `json:"out_path"`
			ErrPath      string `json:"err_path"`
			Username     string `json:"username"`
			SudoPassword string `json:"sudo_password"`
			Stdin        string `json:"stdin"`
		}

		if err := c.ShouldBindJSON(&Res); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		cmd := &LocalCmd{}
		stateCode, out, err := cmd.RunOutput(Res.Shell,
			WithStdin(Res.Stdin),
			WithTimeOut(time.Duration(Res.TimeOut)*time.Second),
			WithOutPath(Res.OutPath),
			WithErrPath(Res.ErrPath),
			WithUsername(Res.Username),
			WithSudoPassword(Res.SudoPassword),
		)
		c.JSON(http.StatusOK, gin.H{
			"state_code": stateCode,
			"out":        out,
			"err":        err.Error(),
		})
	})

	r.POST("/start", func(c *gin.Context) {
		var Res struct {
			Shell        string `json:"shell" binding:"required"`
			TimeOut      int64  `json:"time_out"`
			OutPath      string `json:"out_path"`
			ErrPath      string `json:"err_path"`
			Username     string `json:"username"`
			SudoPassword string `json:"sudo_password"`
			Stdin        string `json:"stdin"`
		}

		if err := c.ShouldBindJSON(&Res); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		cmd := &LocalCmd{}
		pid, err := cmd.Start(Res.Shell,
			WithStdin(Res.Stdin),
			WithTimeOut(time.Duration(Res.TimeOut)*time.Second),
			WithOutPath(Res.OutPath),
			WithErrPath(Res.ErrPath),
			WithUsername(Res.Username),
			WithSudoPassword(Res.SudoPassword),
		)
		c.JSON(http.StatusOK, gin.H{
			"pid": pid,
			"err": err.Error(),
		})
	})
}
