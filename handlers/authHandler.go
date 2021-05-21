package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/lorezi/golang-bank-app-auth/dto"
	"github.com/lorezi/golang-bank-app-auth/ports"
	"github.com/lorezi/golang-bank-app-auth/utils"
)

type AuthHandler struct {
	Service ports.AuthService
}

func (h AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	req := dto.LoginRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Response(w, http.StatusBadRequest, "invalid data 🥵🥵")
		return
	}

	res, err := h.Service.Login(req)
	if err != nil {
		utils.Response(w, err.Code, err.ShowError())
		return
	}

	utils.Response(w, http.StatusOK, res)

}
