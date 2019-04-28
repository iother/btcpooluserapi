package api

import (
	"encoding/json"
	"github.com/iother/btcpooluserapi/config"
	"github.com/iother/btcpooluserapi/db"
	"github.com/iother/btcpooluserapi/log"
	"github.com/iother/btcpooluserapi/util"
	"net/http"
)

type API struct {
	source  *db.Client
	cfg     *config.Config
	version string
}

type Resp struct {
	Err_no  int         `json:"err_no"`
	Err_msg string      `json:"err_msg"`
	Data    interface{} `json:"data"`
}

func NewAPI(source *db.Client, cfg *config.Config, version string) *API {
	s := &API{source: source, cfg: cfg, version: version}
	return s
}

func (s *API) VersionHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(s.version)
}

func (s *API) GetUserList(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	last_id := req.FormValue("last_id")

	log.Info("last_id ", last_id)

	miners, err := s.source.QuerySubAccount(last_id, s.cfg.Coin.Limit)

	if err != nil {
		log.Error("QuerySubAccount err  : ", err)
	}

	minersMap := make(map[string]int64)

	lastMaxUserId := util.ToInt64(last_id)

	if len(miners) > 0 {
		for _, v := range miners {
			if v != nil {
				if v.Sub_account_name.String != "" {

					if v.Sub_account_id > lastMaxUserId {
						lastMaxUserId = v.Sub_account_id
					}

					minersMap[v.Sub_account_name.String] = v.Sub_account_id
				}
			}
		}
	}

	log.Info("lastMaxUserId ", lastMaxUserId)

	rs := Resp{}

	rs.Data = minersMap

	json.NewEncoder(w).Encode(rs)
}
