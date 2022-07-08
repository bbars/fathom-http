package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	gofathom "github.com/bbars/go-fathom"
	"github.com/bbars/limitedpool"
	"github.com/notnil/chess"
)

type ResponseRes struct {
	Res interface{} `json:"res"`
}

type ResponseErr struct {
	Err string `json:"err"`
}

type CmdPosition struct {
	Position *chess.Position `json:"position"`
}

type CmdPositionExt struct {
	CmdPosition
	UseRule50 bool `json:"useRule50"`
}

type fathomRootResult struct {
	Move    gofathom.TbMove     `json:"move"`
	Details []gofathom.TbResult `json:"details"`
}

type requestProcessor func(r *http.Request) (interface{}, error)

type HttpHandlers struct {
	ctx        context.Context
	fathomPool *limitedpool.LimitedPool
	maxTime    time.Duration
}

func (this *HttpHandlers) processRequest(w http.ResponseWriter, r *http.Request, fn requestProcessor) {
	res, err := fn(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(ResponseErr{
			Err: err.Error(),
		})
	} else {
		json.NewEncoder(w).Encode(ResponseRes{
			Res: res,
		})
	}
}

func (this *HttpHandlers) prepare(r *http.Request, payload interface{}) (gofathom.Fathom, error) {
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		return nil, err
	}
	var ctx context.Context
	if this.maxTime > 0 {
		ctx, _ = context.WithTimeout(this.ctx, this.maxTime)
	} else {
		ctx = context.TODO()
	}
	fathom := this.fathomPool.Get(ctx).(gofathom.Fathom)
	if fathom == nil {
		return nil, fmt.Errorf("http-fathom: unable to get engine in %d", this.maxTime)
	}
	return fathom, nil
}

func (this *HttpHandlers) HandleWDL(w http.ResponseWriter, r *http.Request) {
	this.processRequest(w, r, func(r *http.Request) (interface{}, error) {
		var cmd *CmdPosition
		fathom, err := this.prepare(r, &cmd)
		if err != nil {
			return nil, err
		}
		defer this.fathomPool.Put(fathom)

		res, err := fathom.ProbeWDL(cmd.Position)
		if err != nil {
			return res, err
		}

		return res, nil
	})
}

func (this *HttpHandlers) HandleRoot(w http.ResponseWriter, r *http.Request) {
	this.processRequest(w, r, func(r *http.Request) (interface{}, error) {
		var cmd *CmdPosition
		fathom, err := this.prepare(r, &cmd)
		if err != nil {
			return nil, err
		}
		defer this.fathomPool.Put(fathom)

		resMove, resDetails, err := fathom.ProbeRoot(cmd.Position)
		if err != nil {
			return nil, err
		}
		res := fathomRootResult{
			Move:    resMove,
			Details: resDetails,
		}
		return res, nil
	})
}

func (this *HttpHandlers) HandleRootDTZ(w http.ResponseWriter, r *http.Request) {
	this.processRequest(w, r, func(r *http.Request) (interface{}, error) {
		var cmd *CmdPositionExt
		fathom, err := this.prepare(r, &cmd)
		if err != nil {
			return nil, err
		}
		defer this.fathomPool.Put(fathom)

		res, err := fathom.ProbeRootDTZ(cmd.Position, cmd.UseRule50)
		if err != nil {
			return res, err
		}
		return res, nil
	})
}

func (this *HttpHandlers) HandleRootWDL(w http.ResponseWriter, r *http.Request) {
	this.processRequest(w, r, func(r *http.Request) (interface{}, error) {
		var cmd *CmdPositionExt
		fathom, err := this.prepare(r, &cmd)
		if err != nil {
			return nil, err
		}
		defer this.fathomPool.Put(fathom)

		res, err := fathom.ProbeRootWDL(cmd.Position, cmd.UseRule50)
		if err != nil {
			return res, err
		}
		return res, nil
	})
}
